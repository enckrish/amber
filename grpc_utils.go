package main

import (
	"amber/pb"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
)

type GRPCConnector struct {
	conn   *grpc.ClientConn
	client pb.AnalyzerClient
}

func NewGRPCConnector(target string) (*GRPCConnector, error) {
	gc := GRPCConnector{}
	opts := []grpc.DialOption{
		// TODO add authentication to prevent unauthorized access to analyzer
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	if conn, err := grpc.Dial(target, opts...); err != nil {
		return nil, err
	} else {
		gc.conn = conn
	}

	gc.client = pb.NewAnalyzerClient(gc.conn)
	return &gc, nil
}

func (gc *GRPCConnector) Close() error {
	return gc.conn.Close()
}

func (gc *GRPCConnector) sendInitRequestType0(service string, historySize uint32) (*pb.UUID, error) {
	req := &pb.InitRequest_Type0{Service: service, HistorySize: historySize}
	res, err := gc.client.InitType0(context.Background(), req)
	if err != nil {
		return nil, err
	}
	return res.Id, nil
}

func receiveFromStream(stream pb.Analyzer_AnalyzeLog_Type0Client, active *bool) {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			*active = false
			return
		}
		if err != nil {
			*active = false
			log.Println("Failed to receive a response:", err)
			log.Println(activeFalseMsg)
			return
		}
		log.Printf("Got message: %v", in)
	}
}
