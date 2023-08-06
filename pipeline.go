package main

import (
	"amber/pb"
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
}

func NewPipelineConfig(
	services []Service,
	parseMode pb.ParseMode,
	bufferSize int,
	historySize int,
	bufferTimeout time.Duration,
) (*PipelineConfig, error) {
	buffer := NewBuffer(bufferSize, bufferTimeout)
	history := NewBuffer(historySize, 0)

	return &PipelineConfig{
		Id:        uuid.New(),
		Services:  services,
		ParseMode: parseMode,
		buffer:    buffer,
		history:   history,
	}, nil
}

func (p *PipelineConfig) Exec(stream pb.Analyzer_AnalyzeLogClient) {
	var wg sync.WaitGroup

	np := len(p.Services)
	wg.Add(np)
	for _, serv := range p.Services {
		go p.execStream(serv, stream, &wg)
	}

	wg.Wait()
}

func (p *PipelineConfig) execStream(serv Service, stream pb.Analyzer_AnalyzeLogClient, wg *sync.WaitGroup) {
	defer wg.Done()

	t, service, id := serv.stream, serv.name, serv.id
	buffer, history := p.buffer, p.history

	//fmt.Println("Reading", serv.name)
	for line := range t.Lines {
		buffer.Lock()
		history.Lock()

		//fmt.Println("Line", serv.name, line.Text)
		if len(line.Text) == 0 {
			continue
		}

		it := &BufferItemUnref{Id: &pb.UUID{Id: id.String()}, ServName: service, Log: line.Text}
		buffer.Insert(it)
		//printBuffer("Buffer:", &buffer.list)

		// Allowing concurrent reads might cause double-flushing
		// TODO IsTimeout should be triggered even without new line
		if buffer.IsFull() || buffer.IsTimeout() {
			err := p.sendRequest(stream)
			if err != nil {
				// TODO retry after 5 secs?
			}

		}
		history.Unlock()
		buffer.Unlock()
	}
}

// TODO check that buffer and history load-unload is valid (fn execStream)
func (p *PipelineConfig) sendRequest(stream pb.Analyzer_AnalyzeLogClient) error {
	// Don't reset heads now, that way if stream returns error, we can rollback
	flBuffer := p.buffer.GetContentSeq()
	flHistory := p.history.GetContentSeq()
	{
		//printBuffer("Buffer Flush:", &flBuffer)
		//printBuffer("History:", &flHistory)
	}
	req := pb.AnalyzerRequest{
		Id:        &pb.UUID{Id: p.Id.String()},
		Recent:    flBuffer,
		History:   flHistory,
		ParseMode: p.ParseMode,
	}
	if err := stream.Send(&req); err != nil {
		return err
	}

	p.buffer.SoftReset()
	p.history.InsertMultiple(flBuffer)
	return nil
}
