package pipeline

import (
	"github.com/google/uuid"
	"github.com/nxadm/tail"
	"time"
)

type Service struct {
	// Name of the log-producing service (Docker, nginx etc.)
	name string
	// UUID generated automatically for the service
	id string
	// File stream for the log file
	stream *tail.Tail
}

func NewService(name string, logPath string) (*Service, error) {
	t, err := tail.TailFile(logPath, tail.Config{Follow: true, ReOpen: true})
	if err != nil {
		return nil, err
	}

	return &Service{name: name, id: uuid.New().String(), stream: t}, nil
}

type StoreItem struct {
	requestTime time.Time
	logs        []string
	result      *AnalysisResult
}

func (s StoreItem) RequestTime() time.Time {
	return s.requestTime
}

func (s StoreItem) Logs() []string {
	return s.logs
}

func (s StoreItem) Result() *AnalysisResult {
	return s.result
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
