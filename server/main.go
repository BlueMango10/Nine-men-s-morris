package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

const (
	hostAddress = "localhost:32610"
)

var (
	// Logs str
	logI func(str string)
	// Logs str as an error
	logE func(str string)
	// Logs str a fatal error, then panics
	logF func(str string)
)

func main() {
	var buf bytes.Buffer
	var loggerInf = log.New(&buf, "LOG|INFO: ", log.Lshortfile|log.Lmicroseconds)
	var loggerErr = log.New(&buf, "LOG|ERR: ", log.Lshortfile|log.Lmicroseconds)
	var loggerFat = log.New(&buf, "LOG|FATAL: ", log.Lshortfile|log.Lmicroseconds)
	defer func() {
		fmt.Println(&buf)
		os.WriteFile("log.txt", buf.Bytes(), 0644)
	}()
	logI = func(str string) {
		loggerInf.Output(2, str)
		fmt.Println(str)
	}
	logE = func(str string) {
		loggerErr.Output(2, str)
		fmt.Println(str)
	}
	logF = func(str string) {
		loggerFat.Output(2, str)
		panic(str)
	}

	logI(fmt.Sprintf("=== LOG START: %v ===", time.Now().Format(time.RFC1123)))
	var wg sync.WaitGroup
	wg.Add(1)
	go startServer(hostAddress, wg)
	wg.Wait()
	fmt.Println("Press Enter to exit")
	fmt.Scanln()
}
