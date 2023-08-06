package main

import (
	"amber/pb"
	"fmt"
	"github.com/google/uuid"
	"github.com/nxadm/tail"
	"sync"
	"time"
)

type BufferItemUnref = pb.LogInstance
type BufferItem = *BufferItemUnref

// TODO Modify BufferItem to be a generic argument of Buffer

type Buffer struct {
	sync.RWMutex
	list []BufferItem
	// Index to push new items at
	head int
	// Number of items currently in buffer
	len int
	// Timestamp of last insert operation
	lastInsert time.Time
	// Used to detect buffer timeout
	nsTimeout time.Duration
}

func NewBuffer(size int, nsTimeout time.Duration) *Buffer {
	return &Buffer{list: make([]BufferItem, size), nsTimeout: nsTimeout}
}

func (b *Buffer) Insert(entry BufferItem) {
	b.list[b.head] = entry
	b.head = ClIncr(b.head, len(b.list))
	b.len = Min(b.len+1, len(b.list))
	b.lastInsert = time.Now()
}

func (b *Buffer) InsertMultiple(entries []BufferItem) {
	for _, entry := range entries {
		b.Insert(entry)
	}
}

// Flush Returns all elements stored in the current iteration
// and resets head, len to 0. Flush doesn't remove any elements
func (b *Buffer) Flush() []BufferItem {
	f := b.GetContentSeq()
	b.SoftReset()

	return f
}

// SoftReset Resets head and len
func (b *Buffer) SoftReset() {
	b.head = 0
	b.len = 0
}

func (b *Buffer) IsFull() bool {
	return b.len == len(b.list)
}

func (b *Buffer) IsTimeout() bool {
	if b.len == 0 {
		return false
	}
	return time.Now().Sub(b.lastInsert) > b.nsTimeout
}

func (b *Buffer) GetContentSeq() []BufferItem {
	f := make([]BufferItem, b.len)
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
	id uuid.UUID
	// File stream for the log file
	stream *tail.Tail
}

func NewService(name string, logPath string) (Service, error) {
	t, err := tail.TailFile(logPath, tail.Config{Follow: true, ReOpen: true})
	if err != nil {
		return Service{}, err
	}

	return Service{name: name, id: uuid.New(), stream: t}, nil
}

func printBuffer(title string, item *[]BufferItem) {
	fmt.Println(title)
	for i, it := range *item {
		fmt.Println(i, it)
	}
}
