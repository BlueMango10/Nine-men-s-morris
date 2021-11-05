package main

import (
	"bytes"
	"fmt"
	"log"
	"time"
)

var (
	logf func(str string)
)

func main() {
	var buf bytes.Buffer
	var logger = log.New(&buf, "LOG: ", log.Lshortfile)
	defer fmt.Println(&buf)
	logf = func(str string) {
		logger.Output(2, str)
	}
	logf(fmt.Sprintf("=== LOG START: %v ===", time.Now().Format(time.RFC1123)))
}
