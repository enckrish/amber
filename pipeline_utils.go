package main

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/nxadm/tail"
	"time"
)

// TODO should support combined logs from many sources at once, to find connected issues

type Buffer struct {
	list []string
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
	return &Buffer{list: make([]string, size), nsTimeout: nsTimeout}
}

func (a *Buffer) Insert(entry string) {
	a.list[a.head] = entry
	a.head = ClIncr(a.head, len(a.list))
	a.len = Min(a.len+1, len(a.list))
	a.lastInsert = time.Now()
}

func (a *Buffer) InsertMultiple(entries []string) {
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
func (a *Buffer) Flush() []string {
	f := make([]string, a.len)
	j := ClDecr(a.head, len(a.list))
	for i := a.len - 1; i >= 0; i-- {
		f[i] = a.list[j]
		j = ClDecr(j, len(a.list))
	}

	a.head = 0
	a.len = 0

	return f
}

type PipelineConfig struct {
	service     string
	serviceId   uuid.UUID
	ParserFn    TParserFn
	format      string
	logFilePath string
	analyzer    string
	buffer      *Buffer
	history     *Buffer
}

// NewPipelineConfig
// @parserArg0: Refers to pluginPath instead of format, if usePlugin is enabled
func NewPipelineConfig(service string, usePlugin bool, parserArg0, logFilePath string, bufferSize, historySize int, bufferTimeout time.Duration) (PipelineConfig, error) {
	var parserFn TParserFn
	var err error
	format := parserArg0
	if usePlugin {
		parserFn, err = GetPluginParser(parserArg0)
		format = ""
	} else {
		parserFn = ParseLog
	}

	if err != nil {
		return PipelineConfig{}, err
	}

	buffer := NewBuffer(bufferSize, bufferTimeout)
	history := NewBuffer(historySize, 0)
	return PipelineConfig{
		service:     service,
		serviceId:   uuid.New(),
		ParserFn:    parserFn,
		format:      format,
		logFilePath: logFilePath,
		analyzer:    "",
		buffer:      buffer,
		history:     history,
	}, nil
}

func (p PipelineConfig) Exec() error {
	// TODO Multithreading should start from here

	t, err := tail.TailFile(p.logFilePath, tail.Config{Follow: true, ReOpen: false})
	if err != nil {
		return err
	}
	for line := range t.Lines {
		res, err := p.ParserFn(p.format, line.Text)
		if err != nil {
			return err
		}
		fmt.Println(res)

		p.buffer.Insert(res)
		if p.buffer.IsFull() || p.buffer.IsTimeout() {
			// Ideally IsTimeout should be triggered even without new line
			list := p.buffer.Flush()
			fmt.Println("Buffer Flush:", list)
			// TODO send to analyzer
		}
	}

	return nil
}
