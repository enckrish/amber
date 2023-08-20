package main

//
//import (
//	"amber/pb"
//	"os"
//	"time"
//)
//
//const analyzerTarget = "localhost:50051"
//const testTimeout = 2 * time.Second
//
//var brokers = []string{os.Getenv("AMBER_KAFKA_URL")}
//
//func test(logPath string, serviceType string, bufferSize int, historySize uint32) {
//	conn, err := NewGRPCConnector(analyzerTarget)
//	if err != nil {
//		panic(err)
//	}
//	defer func(conn *GRPCConnector) {
//		err := conn.Close()
//		if err != nil {
//			panic(err)
//		}
//	}(conn)
//
//	service, err := NewService(serviceType, logPath)
//	if err != nil {
//		panic(err)
//	}
//
//	p, err := NewPipelineConfig(conn, brokers, service, bufferSize, historySize, testTimeout)
//	if err != nil {
//		panic(err)
//	}
//
//	go runUI(p)
//	p.Exec(pb.NewRouterClient(conn.conn))
//}
