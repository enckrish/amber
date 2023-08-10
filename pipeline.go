package main

import (
	"amber/pb"
	"context"
	"log"
	"time"
)

const streamPushRetryErr = "Failed to send logs to server. Retrying..."
const maxRetryErr = "Maximum retries exhausted! No new attempts will be made for this stream"
const activeFalseMsg = "Stream is set to inactive"
const sendingLogMsg = "Sending logs for analysis"

// PipelineConfig
//
// Stores parameters for executing analysis of live streams
// Should only be initialized using NewPipelineConfig
type PipelineConfig struct {
	// Stream identifier for use at analyzer
	id      *pb.UUID
	service *Service
	// Stores whether server is taking any more requests
	active        bool
	buffer        *Buffer[string]
	retries       int32
	watchInterval time.Duration
}

func (p *PipelineConfig) Id() *pb.UUID {
	return p.id
}

func (p *PipelineConfig) Service() Service {
	return Service{name: p.service.name, id: pb.UUID{Id: p.service.id.Id}, stream: nil}
}

func (p *PipelineConfig) Active() bool {
	return p.active
}

func (p *PipelineConfig) WatchInterval() time.Duration {
	return p.watchInterval
}

func NewPipelineConfig(
	connector *GRPCConnector,
	service *Service,
	bufferSize int,
	historySize uint32,
	bufferTimeout time.Duration,
) (*PipelineConfig, error) {
	buffer := NewBuffer[string](bufferSize, bufferTimeout)

	id, err := connector.sendInitRequestType0(service.name, historySize)
	if err != nil {
		return nil, err
	}

	return &PipelineConfig{
		id:            id,
		service:       service,
		buffer:        buffer,
		active:        true,
		watchInterval: bufferTimeout / 2,
	}, nil
}

func (p *PipelineConfig) Exec(client pb.AnalyzerClient) error {
	// TODO auth context
	stream, err := client.AnalyzeLog_Type0(context.Background())
	if err != nil {
		return err
	}

	go receiveFromStream(stream, &p.active)
	go p.watchBufferTimeout(stream)

	p.execStream(p.service, stream)

	return stream.CloseSend()
}

func (p *PipelineConfig) execStream(service *Service, stream pb.Analyzer_AnalyzeLog_Type0Client) {
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
			}
		}
		buffer.Unlock()
	}
	log.Println(activeFalseMsg)
	p.active = false
}

// watchBufferTimeout runs at every watchInterval, if buffer is timeout, then tries to send to stream for retries times
// Should be always called as a goroutine.
func (p *PipelineConfig) watchBufferTimeout(stream pb.Analyzer_AnalyzeLog_Type0Client) {
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

// trySendToStream tries to send the buffered logs to the stream
func (p *PipelineConfig) trySendToStream(stream pb.Analyzer_AnalyzeLog_Type0Client) error {
	log.Println(sendingLogMsg)
	req := p.getRequest()
	return stream.Send(req)
}

// getRequest sends the logs in order, wrapped as pb.AnalyzerRequest_Type0
// Does not reset the buffer state
func (p *PipelineConfig) getRequest() *pb.AnalyzerRequest_Type0 {
	flBuffer := p.buffer.GetContentSeq()

	req := &pb.AnalyzerRequest_Type0{
		Id:   p.id,
		Logs: &pb.LogInstance_Type0{Logs: flBuffer},
	}

	return req
}
