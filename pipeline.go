package main

import (
	"amber/pb"
	"context"
	"fmt"
	"github.com/google/uuid"
	"sync"
	"time"
)

// PipelineConfig
//
// Stores parameters for executing analysis of live streams
// Should only be initialized using NewPipelineConfig
//
// PipelineConfig.buffer and PipelineConfig.history may exchange references:
// their contents should not be updated outside PipelineConfig methods
type PipelineConfig struct {
	// Stream identifier for use at analyzer
	Id       uuid.UUID
	Services []Service
	// Sets whether log will be parsed before analysis or passed as it is
	ParseMode pb.ParseMode
	buffer    *Buffer
	history   *Buffer
	// Stores whether server is taking any more requests
	active bool
}

func NewPipelineConfig(
	services []Service,
	parseMode pb.ParseMode,
	bufferSize int,
	historySize int,
	bufferTimeout time.Duration,
) *PipelineConfig {
	buffer := NewBuffer(bufferSize, bufferTimeout)
	history := NewBuffer(historySize, 0)

	return &PipelineConfig{
		Id:        uuid.New(),
		Services:  services,
		ParseMode: parseMode,
		buffer:    buffer,
		history:   history,
		active:    true,
	}
}

func (p *PipelineConfig) Exec(client pb.AnalyzerClient) {
	stream, err := client.AnalyzeLog(context.Background())
	if err != nil {
		panic(err)
	}
	defer stream.CloseSend()

	reqChannel := make(chan *pb.AnalyzerRequest)
	go sendToStream(stream, reqChannel)
	go receiveFromStream(stream, reqChannel, &p.active)

	wg := sync.WaitGroup{}
	np := len(p.Services)
	wg.Add(np)
	for _, serv := range p.Services {
		go p.execStream(serv, reqChannel, &wg)
	}

	wg.Wait()
}

func (p *PipelineConfig) execStream(serv Service, reqChannel chan *pb.AnalyzerRequest, wg *sync.WaitGroup) {
	defer wg.Done()

	t, service, id := serv.stream, serv.name, serv.id
	buffer, history := p.buffer, p.history

	fmt.Println("Reading", serv.name)
	for line := range t.Lines {
		buffer.Lock()
		history.Lock()

		fmt.Println("Line", serv.name, line.Text)
		if len(line.Text) == 0 {
			continue
		}

		it := &BufferItemUnref{Id: &pb.UUID{Id: id.String()}, ServName: service, Log: line.Text}
		buffer.Insert(it)
		//printBuffer("Buffer:", &buffer.list)

		// Allowing concurrent reads might cause double-flushing
		// TODO IsTimeout should be triggered even without new line
		if buffer.IsFull() || buffer.IsTimeout() {
			if !p.active {
				return
			}
			reqChannel <- p.getRequests()
		}
		history.Unlock()
		buffer.Unlock()
	}
}

// TODO check that buffer and history load-unload is valid
func (p *PipelineConfig) getRequests() *pb.AnalyzerRequest {
	flBuffer := p.buffer.Flush()
	flHistory := p.history.GetContentSeq()

	req := &pb.AnalyzerRequest{
		Id:        &pb.UUID{Id: p.Id.String()},
		Recent:    flBuffer,
		History:   flHistory,
		ParseMode: p.ParseMode,
	}
	p.history.InsertMultiple(flBuffer)

	return req
}

func (p *PipelineConfig) isActive() bool {
	return p.active
}
