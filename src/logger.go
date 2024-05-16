package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

func logInit(filename string, console bool) *log.Logger {
	var f *os.File
	var err error

	if filename != "" {
		f, err = os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm) // #nosec G304
		if err != nil {
			fmt.Printf("Can't open the log file: %v\n", err)
			console = true
		}
	}

	var wr io.Writer
	if console {
		wr = io.MultiWriter(os.Stdout, f)
	} else {
		wr = f
	}
	return log.New(wr, "", log.LstdFlags)
}
