package main

import (
	"github.com/google/uuid"
	"github.com/nxadm/tail"
	"path"
	"strings"
	"sync"
	"time"
)

type Buffer[ItemType any] struct {
	sync.RWMutex
	list []ItemType
	// Index to push new items at
	head int
	// Number of items currently in buffer
	len int
	// Timestamp of last insert operation
	lastInsert time.Time
	// Used to detect buffer timeout
	nsTimeout time.Duration
}

func NewBuffer[ItemType any](size int, nsTimeout time.Duration) *Buffer[ItemType] {
	return &Buffer[ItemType]{list: make([]ItemType, size), nsTimeout: nsTimeout}
}

func (b *Buffer[ItemType]) Insert(entry ItemType) {
	b.list[b.head] = entry
	b.head = ClIncr(b.head, len(b.list))
	b.len = Min(b.len+1, len(b.list))
	b.lastInsert = time.Now()
}

func (b *Buffer[ItemType]) InsertMultiple(entries []ItemType) {
	for _, entry := range entries {
		b.Insert(entry)
	}
}

// Flush Returns all elements stored in the current iteration
// and resets head, len to 0. Flush doesn't remove any elements
func (b *Buffer[ItemType]) Flush() []ItemType {
	f := b.GetContentSeq()
	b.SoftReset()

	return f
}

// SoftReset Resets head and len
func (b *Buffer[ItemType]) SoftReset() {
	b.head = 0
	b.len = 0
}

func (b *Buffer[ItemType]) IsFull() bool {
	return b.len == len(b.list)
}

func (b *Buffer[ItemType]) IsTimeout() bool {
	if b.len == 0 {
		return false
	}
	return time.Now().Sub(b.lastInsert) > b.nsTimeout
}

func (b *Buffer[ItemType]) GetContentSeq() []ItemType {
	f := make([]ItemType, b.len)
	// j starts from b.head-1 and decrements upto b.len iterations
	j := ClDecr(b.head, len(b.list))
	for i := b.len - 1; i >= 0; i-- {
		f[i] = b.list[j]
		j = ClDecr(j, len(b.list))
	}

	return f
}

type Service struct {
	// Name of the log-producing service (Docker, nginx etc.)
	name string
	// UUID generated automatically for the service
	id string
	// File stream for the log file
	stream *tail.Tail
}

func NewService(name string, logPath string) (*Service, error) {
	logPath, err := handleFileType(logPath)
	t, err := tail.TailFile(logPath, tail.Config{Follow: true, ReOpen: true})
	if err != nil {
		return nil, err
	}

	return &Service{name: name, id: uuid.New().String(), stream: t}, nil
}

func handleFileType(logPath string) (string, error) {
	createTempFile()
	var err error

	ext := path.Ext(logPath)
	switch strings.ToLower(ext) {
	case ".pdf":
		err = readPdf(logPath)
	default:
		return logPath, nil
	}
	return tempPath, err
}

type StoreItem struct {
	requestTime time.Time
	logs        []string
	result      *AnalysisResult
}

type AnalysisResult struct {
	StreamId  string   `json:"stream_id"`
	MessageId string   `json:"message_id"`
	Rating    int      `json:"rating"`
	Actions   []string `json:"actions"`
	Review    string   `json:"review"`
	Citation  int      `json:"citation"`
}

func (a *AnalysisResult) StrRating() string {
	ratingMap := []string{"err", "none", "low", "medium", "high", "critical"}
	return ratingMap[a.Rating]
}
