package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

var doTimeStamp = true
var doTee = true
var toFile = "/tmp/output.log"

var logFP *os.File

func init() {
	flag.BoolVar(&doTimeStamp, "ts", doTimeStamp, "add timestamps to output")
	flag.BoolVar(&doTee, "tee", doTee, "act like tee (echo input to stdout")
	flag.StringVar(&toFile, "f", toFile, "filename to write")
}

func main() {
	flag.Parse()

	if _, err := os.Stat(toFile); err == nil {
		var filename = toFile
		for i := uint(0); i < ^uint(0); i++ {
			filename = fmt.Sprintf("%s.%d", toFile, i)
			if _, err := os.Stat(filename); os.IsNotExist(err) {
				toFile = filename
				break
			}
		}
	}

	if fp, err := os.Create(toFile); err != nil {
		log.Fatal(err)
	} else {
		logFP = fp
	}

	fmt.Printf("Writing to %s\n\n", toFile)

	var writeTo io.Writer

	if doTee {
		writeTo = io.MultiWriter(logFP, os.Stdout)
	} else {
		writeTo = logFP
	}

	log.SetOutput(writeTo)

	if doTimeStamp {
		log.SetFlags(log.Ldate | log.Lmicroseconds)
	} else {
		log.SetFlags(0)
	}

	stdin := bufio.NewReader(os.Stdin)

	for {
		if line, err := stdin.ReadString('\n'); err == nil {
			log.Print(line)
		} else {
			if line != "" {
				log.Print(line)
			}
			if err != io.EOF {
				log.SetOutput(os.Stderr)
				log.Fatal(err)
			}
			return
		}
	}
}
