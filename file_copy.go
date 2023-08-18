package main

import (
	"bytes"
	"github.com/dslipak/pdf"
	"log"
	"os"
)

const tempPath = "/tmp/amber_temp.txt"

var f *os.File

func createTempFile() {
	var err error
	f, err = os.Create(tempPath)
	if err != nil {
		log.Panicln("Failed to create dump file:", err)
	}
}

func readPdf(path string) error {
	r, err := pdf.Open(path)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	b, err := r.GetPlainText()
	if err != nil {
		return err
	}
	_, err = buf.ReadFrom(b)
	if err != nil {
		return err
	}
	_, err = f.Write(buf.Bytes())
	return err
}
