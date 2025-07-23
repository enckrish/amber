package main

import (
	"amber/pipeline"
	"amber/pipeline/pb"
	"flag"
	"io"
	"log"
	"os"
)
import (
	"time"
)

const analyzerTarget = "localhost:50051"
const timeout = 2 * time.Second

var brokers = []string{os.Getenv("AMBER_KAFKA_URL")}
var (
	logPath     string
	service     string
	bufferSize  int
	historySize int
	showHelp    bool
)

func init() {
	flag.StringVar(&logPath, "p", "", "path to the log file to be analyzed")
	flag.StringVar(&service, "t", "", "name of the log producing service")
	flag.IntVar(&bufferSize, "bs", 1, "number of log lines to read before analyzing")
	flag.IntVar(&historySize, "hs", 0, "number of entries to keep in history for analysis context")
	flag.BoolVar(&showHelp, "h", false, "show help")
}
func main() {

	flag.Parse()

	if showHelp {
		flag.Usage()
		os.Exit(2)
	}

	if logPath == "" || service == "" {
		log.Panicln("Invalid values for `p` or `t`")
	}
	log.SetOutput(io.Discard)
	//test(logPath, service, bufferSize, uint32(historySize))

	conn, err := pipeline.NewGRPCConnector(analyzerTarget)
	if err != nil {
		panic(err)
	}
	defer func(conn *pipeline.GRPCConnector) {
		err := conn.Close()
		if err != nil {
			panic(err)
		}
	}(conn)

	service, err := pipeline.NewService(service, logPath)
	if err != nil {
		panic(err)
	}

	p, err := pipeline.NewPipelineConfig(conn, brokers, service, bufferSize, uint32(historySize), timeout, onUpdateSample)
	if err != nil {
		panic(err)
	}

	go runUI(p)
	p.Exec(pb.NewRouterClient(conn.Conn()))
}
