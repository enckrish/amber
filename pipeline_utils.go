package main

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/nxadm/tail"
	"sync"
	"time"
)

type LogInstance struct {
	Service   string
	ServiceId uuid.UUID
	Log       string
}

type BufferItem = LogInstance

type Buffer struct {
	sync.RWMutex
	list []BufferItem
	// Index to push new items at
	head int
	// Number of items currently in buffer
	len int
	// Last time item is inserted in the buffer
	lastInsert time.Time
	// Used to detect buffer timeout
	nsTimeout time.Duration
}

func NewBuffer(size int, nsTimeout time.Duration) *Buffer {
	return &Buffer{list: make([]BufferItem, size), nsTimeout: nsTimeout}
}

func (a *Buffer) Insert(entry BufferItem) {
	a.list[a.head] = entry
	a.head = ClIncr(a.head, len(a.list))
	a.len = Min(a.len+1, len(a.list))
	a.lastInsert = time.Now()
}

func (a *Buffer) InsertMultiple(entries []BufferItem) {
	for _, entry := range entries {
		a.Insert(entry)
	}
}

func (a *Buffer) IsFull() bool {
	return a.len == len(a.list)
}

func (a *Buffer) IsTimeout() bool {
	if a.len == 0 {
		return false
	}
	return time.Now().Sub(a.lastInsert) > a.nsTimeout
}

// Flush Returns all elements stored in the current iteration
// and resets head, len to 0. Flush doesn't remove any elements
func (a *Buffer) Flush() []BufferItem {
	f := a.GetContentSeq()

	a.head = 0
	a.len = 0

	return f
}

func (a *Buffer) GetContentSeq() []BufferItem {
	f := make([]BufferItem, a.len)
	j := ClDecr(a.head, len(a.list))
	for i := a.len - 1; i >= 0; i-- {
		f[i] = a.list[j]
		j = ClDecr(j, len(a.list))
	}

	return f
}

type ParseMode int32

const (
	Parsed ParseMode = iota
	Unparsed
)

type PipelineConfig struct {
	services   []string
	parseMode  ParseMode
	serviceIds []uuid.UUID
	streams    []*tail.Tail
	buffer     *Buffer
	history    *Buffer
}

// NewPipelineConfig
// @parserArg0: Refers to pluginPath instead of format, if usePlugin is enabled
func NewPipelineConfig(services []string, logPaths []string, parseMode ParseMode, bufferSize int, historySize int, bufferTimeout time.Duration) (*PipelineConfig, error) {
	numServices := len(services)

	serviceIds := make([]uuid.UUID, numServices)
	for i, _ := range serviceIds {
		serviceIds[i] = uuid.New()
	}

	streams := make([]*tail.Tail, numServices)
	for i, _ := range streams {
		t, err := tail.TailFile(logPaths[i], tail.Config{Follow: true, ReOpen: false})
		if err != nil {
			return nil, err
		}
		streams[i] = t
	}

	buffer := NewBuffer(bufferSize, bufferTimeout)
	history := NewBuffer(historySize, 0)

	return &PipelineConfig{
		services:   services,
		parseMode:  parseMode,
		serviceIds: serviceIds,
		streams:    streams,
		buffer:     buffer,
		history:    history,
	}, nil
}

func (p *PipelineConfig) Exec() {
	wg := sync.WaitGroup{}

	np := len(p.services)
	wg.Add(np)

	for i := 0; i < np; i++ {
		go p.execStream(i, &wg)
	}

	wg.Wait()
}

func (p *PipelineConfig) execStream(id int, wg *sync.WaitGroup) {
	defer wg.Done()

	// Read operations on local ids only
	t := p.streams[id]
	service := p.services[id]
	serviceId := p.serviceIds[id]
	// Except p.buffer and p.history, no other members are accessed after this point
	buffer := p.buffer
	history := p.history
	for line := range t.Lines {
		it := BufferItem{Service: service, Log: line.Text, ServiceId: serviceId}

		buffer.Lock()
		history.Lock()

		buffer.Insert(it)

		// Allowing concurrent reads might cause double-flushing
		if buffer.IsFull() || buffer.IsTimeout() {
			// TODO IsTimeout should be triggered even without new line
			flBuffer := buffer.Flush()
			flHistory := history.GetContentSeq()
			history.InsertMultiple(flBuffer)
			fmt.Println("Buffer Flush:", flBuffer, "History:", flHistory)

			// TODO send to analyzer
		}
		history.Unlock()
		buffer.Unlock()
	}
}

// TODO check that slice refs are not returned where not wanted
// TODO check that buffer and history load-unload is valid (fn execStream)
