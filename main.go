package main

import (
	"amber/pb"
	"fmt"
	"time"
)

const AnalyzerTarget = "localhost:50051"

func main() {
	conn, client, err := getStream(AnalyzerTarget)
	if err != nil {
		fmt.Println("Connection error:")
		panic(err)
	}

	services := make([]Service, 1)
	services[0], err = NewService("S1", "/home/krishnendu-wsl/logtests/a.log")
	//services[1], err = NewService("S2", "/home/krishnendu-wsl/logtests/b.log")

	if err != nil {
		panic(err)
	}
	p := NewPipelineConfig(services, pb.ParseMode_Unparsed, 1, 2, 10*time.Second)
	p.Exec(client)

	conn.Close()
}
