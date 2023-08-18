package main

import (
	"flag"
	"io"
	"log"
	"os"
)

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
	flag.IntVar(&historySize, "hs", 1, "number of entries to keep in history for analysis context")
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
	test(logPath, service, bufferSize, uint32(historySize))
}
