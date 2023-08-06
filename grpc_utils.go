package main

import (
	"amber/pb"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
)

// Note: defer conn.Close() in main
func getStream(target string) (*grpc.ClientConn, pb.Analyzer_AnalyzeLogClient, error) {
	conn, err := grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}

	client := pb.NewAnalyzerClient(conn)
	stream, err := client.AnalyzeLog(context.Background())
	if err != nil {
		return nil, nil, nil
	}
	return conn, stream, nil
}

func receive(stream pb.Analyzer_AnalyzeLogClient) {
	waitc := make(chan pb.AnalyzerResponse)
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				close(waitc)
			}
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("Analyzer Response [", in.Id, "]:", in)
		}
	}()
}
