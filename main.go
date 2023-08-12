package main

import (
	"amber/pb"
	"time"
)

const analyzerTarget = "localhost:50051" // "0.tcp.in.ngrok.io:17400"

const testLogService = "ABC"
const testLogPath = "/home/krishnendu-wsl/log-tests/a.log"

const testBufferSize = 3
const testHistorySize = 10
const testTimeout = 2 * time.Second

func main() {
	conn, err := NewGRPCConnector(analyzerTarget)
	if err != nil {
		panic(err)
	}
	defer func(conn *GRPCConnector) {
		err := conn.Close()
		if err != nil {
			panic(err)
		}
	}(conn)

	service, err := NewService(testLogService, testLogPath)
	if err != nil {
		panic(err)
	}

	p, err := NewPipelineConfig(conn, service, testBufferSize, testHistorySize, testTimeout)
	if err != nil {
		panic(err)
	}

	err = p.Exec(pb.NewRouterClient(conn.conn))
	if err != nil {
		panic(err)
	}
}
