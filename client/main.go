package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/BlueMango10/Nine-men-s-morris/morris"
	"google.golang.org/grpc"
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
	startClient(hostAddress)
}

func startClient(address string) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		logF(err.Error())
	}
	defer conn.Close()

	client := morris.NewMorrisClient(conn)
	boardStream, err := client.GetBoardStream(context.Background(), &morris.Empty{})
	if err != nil {
		logF(err.Error())
	}
	go boardStreamReader(boardStream)
	var (
		from int32
		to   int32
	)
	for {
		fmt.Scanln(&from, &to)
		client.MakeMove(context.Background(), &morris.Move{
			From: from,
			To:   to,
		})
	}
}

func boardStreamReader(stream morris.Morris_GetBoardStreamClient) {
	for {
		board, err := stream.Recv()
		if err == io.EOF {
			logI("EOF")
			break
		}
		if err != nil {
			logF(err.Error())
		}
		logI(board.Turn.Visualize(""))
		fmt.Print(board.Visualize(true))
	}
}
