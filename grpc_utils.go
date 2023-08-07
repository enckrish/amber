package main

import (
	"amber/pb"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"time"
)

// Note: defer conn.Close() in main
func getStream(target string) (*grpc.ClientConn, pb.AnalyzerClient, error) {
	opts := []grpc.DialOption{
		// TODO add authentication to prevent unauthorized access to analyzer
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.Dial(target, opts...)
	if err != nil {
		return nil, nil, err
	}

	client := pb.NewAnalyzerClient(conn)
	return conn, client, nil
}

func getAnalyzeLogClient(client *pb.AnalyzerClient) (*pb.Analyzer_AnalyzeLogClient, error) {
	stream, err := (*client).AnalyzeLog(context.Background())
	if err != nil {
		return nil, err
	} else {
		return &stream, nil
	}
}

func sendToStream(stream pb.Analyzer_AnalyzeLogClient, reqChannel chan *pb.AnalyzerRequest) {
	for req := range reqChannel {
		if err := stream.Send(req); err != nil {
			log.Fatalf("Failed to send request: %v", err)
		}
	}
	stream.CloseSend()
}

func receiveFromStream(stream pb.Analyzer_AnalyzeLogClient, reqChannel chan *pb.AnalyzerRequest, requestOpen *bool) {
	for {
		fmt.Println("In:", "in")
		in, err := stream.Recv()
		fmt.Println("In2:", "in")
		if err == io.EOF {
			// prevents closed channel write panic,
			// in case channel is closed after requestOpen check
			*requestOpen = false
			time.AfterFunc(time.Second, func() {
				close(reqChannel)
			})
			return
		}
		if err != nil {
			log.Fatalf("Failed to receive a response : %v", err)
		}
		log.Printf("Got message: %v", in)
	}
}
