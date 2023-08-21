package main

import (
	"amber/pb"
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/google/uuid"
	orderedmap "github.com/wk8/go-ordered-map"
	"log"
	"time"
)

const streamPushRetryErr = "Failed to send logs to server. Retrying..."
const maxRetryErr = "Maximum retries exhausted! No new attempts will be made for this stream"
const activeFalseMsg = "Stream is set to inactive"
const sendingLogMsg = "Sending logs for analysis"
const analysisStoreTopic = "topic.log.analysis.result.1"

// Pipeline
//
// Stores parameters for executing analysis of live streams
// Should only be initialized using NewPipelineConfig
type Pipeline struct {
	// Stream identifier for use at analyzer
	streamId string
	service  *Service
	// Stores whether server is taking any more requests
	active        bool
	buffer        *Buffer[string]
	retries       int32
	watchInterval time.Duration
	store         *orderedmap.OrderedMap
	consumer      sarama.Consumer
	OnUpdateFn    func(pipeline *Pipeline)
	// TODO add trend
}

func (p *Pipeline) Id() string {
	return p.streamId
}

func (p *Pipeline) Service() Service {
	return Service{name: p.service.name, id: p.service.id, stream: nil}
}

func (p *Pipeline) Active() bool {
	return p.active
}

func (p *Pipeline) WatchInterval() time.Duration {
	return p.watchInterval
}

func NewPipelineConfig(
	connector *GRPCConnector,
	kafkaBrokers []string,
	service *Service,
	bufferSize int,
	historySize uint32,
	bufferTimeout time.Duration,
) (*Pipeline, error) {
	buffer := NewBuffer[string](bufferSize, bufferTimeout)

	id, err := connector.sendInitRequestType0(service.name, historySize)
	if err != nil {
		panic(err)
	}

	consumer, err := sarama.NewConsumer(kafkaBrokers, nil)
	if err != nil {
		panic(err)
	}

	return &Pipeline{
		streamId:      id,
		service:       service,
		buffer:        buffer,
		active:        true,
		retries:       10,
		watchInterval: bufferTimeout / 2,
		store:         orderedmap.New(),
		consumer:      consumer,
		OnUpdateFn:    onUpdateSample,
	}, nil
}

func (p *Pipeline) Exec(client pb.RouterClient) error {
	stream, err := client.RouteLog_Type0(context.Background())
	if err != nil {
		return err
	}

	go receiveFromStream(stream, &p.active)
	go p.watchBufferTimeout(stream)
	go p.Listen(p.streamId)

	p.execStream(p.service, stream)

	return stream.CloseSend()
}

func (p *Pipeline) execStream(service *Service, stream pb.Router_RouteLog_Type0Client) {
	t := service.stream
	buffer := p.buffer

	for line := range t.Lines {
		//fmt.Println(line.Text)
		l := line.Text
		if len(l) == 0 {
			continue
		}
		it := l
		buffer.Insert(it)

		// Lock is used to guarantee concurrency with watchBufferTimeout
		buffer.Lock()
		if buffer.IsFull() {
			if !p.active {
				return
			}
			if err := p.trySendToStream(stream); err != nil {
				log.Println(streamPushRetryErr)
			} else {
				p.buffer.SoftReset()
				// To prevent rapidly overflowing the to-router traffic
				time.Sleep(time.Duration(0.01 * float64(time.Second)))
			}
		}
		buffer.Unlock()
	}
	log.Println(activeFalseMsg)
	p.active = false
}

// watchBufferTimeout runs at every watchInterval, if buffer is timeout, then tries to send to stream for retries times
// Should be always called as a goroutine.
func (p *Pipeline) watchBufferTimeout(stream pb.Router_RouteLog_Type0Client) {
	maxRetries := p.retries
	for range time.Tick(p.watchInterval) {
		if !p.buffer.IsTimeout() {
			continue
		}
		if !p.active {
			return
		}
		if p.retries < 0 {
			break
		}
		// Lock is used to guarantee concurrency with execStream
		p.buffer.Lock()

		err := p.trySendToStream(stream)
		if err != nil {
			log.Println(streamPushRetryErr)
			p.retries--
		} else {
			p.retries = maxRetries
			p.buffer.SoftReset()
		}

		p.buffer.Unlock()
	}
	p.active = false
	log.Println(maxRetryErr)
	log.Println(activeFalseMsg)

}

// trySendToStream tries to send the buffered logsBox to the stream
func (p *Pipeline) trySendToStream(stream pb.Router_RouteLog_Type0Client) error {
	log.Println(sendingLogMsg)
	req := p.getRequest()
	p.store.Set(req.MessageId, StoreItem{
		requestTime: time.Now(),
		logs:        req.Logs,
		result:      nil,
	})
	p.OnUpdateFn(p)
	return stream.Send(req)
}

// getRequest sends the logsBox in order, wrapped as pb.AnalyzerRequest_Type0
// Does not reset the buffer state
func (p *Pipeline) getRequest() *pb.AnalyzerRequest_Type0 {
	flBuffer := p.buffer.GetContentSeq()

	req := &pb.AnalyzerRequest_Type0{
		StreamId:  p.streamId,
		MessageId: uuid.New().String(),
		Logs:      flBuffer,
	}

	return req
}

func (p *Pipeline) Listen(streamId string) {
	topic := analysisStoreTopic                                   //e.g. user-created-topic
	partitionList, _ := p.consumer.Partitions(analysisStoreTopic) //get all partitions
	initialOffset := sarama.OffsetOldest                          //offset to start reading message from
	for _, partition := range partitionList {
		pc, _ := p.consumer.ConsumePartition(topic, partition, initialOffset)
		go getMessagesOfId(pc, streamId, p.store, func() {
			p.OnUpdateFn(p)
		})
	}
}

func getMessagesOfId(pc sarama.PartitionConsumer, streamId string, store *orderedmap.OrderedMap, onUpdateFn func()) {
	for msg := range pc.Messages() {
		result := AnalysisResult{}
		if err := json.Unmarshal(msg.Value, &result); err != nil {
			log.Println(err)
			continue
		}
		if result.StreamId != streamId {
			continue
		}
		log.Printf("%+v\n", result)

		prev, ok := store.Get(result.MessageId)
		if !ok {
			fmt.Println("Result available for unaccounted messageId. Skipping")
			continue
		}
		p := prev.(StoreItem)
		item := StoreItem{requestTime: p.requestTime, logs: p.logs, result: &result}
		store.Delete(result.MessageId)
		store.Set(result.MessageId, item)
		onUpdateFn()
	}
}

func (p *Pipeline) GetResultsOfMessage(messageId string) (StoreItem, bool) {
	res, ok := p.store.Get(messageId)
	if !ok {
		return StoreItem{}, ok
	}
	return res.(StoreItem), ok
}

const v = 4

func (p *Pipeline) GetUIDataOverview() [][v]string {
	res := make([][v]string, 0)
	for pair := p.store.Newest(); pair != nil; pair = pair.Prev() {
		messageId := pair.Key.(string)
		data := pair.Value.(StoreItem)
		timestamp := data.requestTime.Round(0).String()
		resultsFetched := data.result != nil
		var status, severity string
		if resultsFetched {
			status = "DONE"
			severity = data.result.StrRating()
		} else {
			status = "WAITING"
			severity = ""
		}
		values := [v]string{timestamp, status, messageId, severity}
		res = append(res, values)
	}
	return res
}

func onUpdateSample(p *Pipeline) {
	//app.Draw()
	if overviewTable != nil {
		updateOverviewTableData(p)
		app.Draw()
	}
}
