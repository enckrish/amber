package main

import (
	"amber/pb"
	"fmt"
	"time"
)

const AnalyzerTarget = "localhost:8080"

func main() {
	conn, stream, err := getStream(AnalyzerTarget)
	if err != nil {
		fmt.Println("Connection error:")
		panic(err)
	}

	defer conn.Close()
	defer stream.CloseSend()
	go receive(stream)

	services := make([]Service, 2)
	services[0], err = NewService("S1", "/home/krishnendu-wsl/logtests/a.log")
	services[1], err = NewService("S2", "/home/krishnendu-wsl/logtests/b.log")

	if err != nil {
		panic(err)
	}
	p, err := NewPipelineConfig(services, pb.ParseMode_Unparsed, 5, 2, 10*time.Second)
	p.Exec(stream)
}
