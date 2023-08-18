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
	client pb.RouterClient
}

func NewGRPCConnector(target string) (*GRPCConnector, error) {
	gc := GRPCConnector{}
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.Dial(target, opts...)
	if err != nil {
		return nil, err
	}
	gc.conn = conn

	gc.client = pb.NewRouterClient(gc.conn)
	return &gc, nil
}

func (gc *GRPCConnector) Close() error {
	return gc.conn.Close()
}

func (gc *GRPCConnector) sendInitRequestType0(service string, historySize uint32) (string, error) {
	req := &pb.InitRequest_Type0{Service: service, HistorySize: historySize}
	res, err := gc.client.Init_Type0(context.Background(), req)
	if err != nil {
		return "", err
	}
	return res.StreamId, nil
}

// This might cause a race condition somewhere so keep in mind
func receiveFromStream(stream pb.Router_RouteLog_Type0Client, active *bool) {
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
