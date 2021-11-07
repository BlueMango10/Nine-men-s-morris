package main

import (
	"context"
	"fmt"
	"net"
	"sync"

	"github.com/BlueMango10/Nine-men-s-morris/morris"
	"google.golang.org/grpc"
)

type MorrisServer struct {
	morris.UnimplementedMorrisServer
	state       morris.BoardState
	streamChans []*chan *morris.BoardState
}

func newMorrisServer() *MorrisServer {
	return &MorrisServer{
		state: morris.BoardState{
			Turn:  morris.BoardSpace_WHITE,
			Board: make([]morris.BoardSpace, 24),
		},
	}
}

func startServer(address string, wg sync.WaitGroup) {
	logI(fmt.Sprintf("STARTING SERVER WITH ADDRESS: %s", address))
	lis, err := net.Listen("tcp", address)
	if err != nil {
		logF(err.Error())
	}
	var opts []grpc.ServerOption
	// Server options here (none yet)

	grpcServer := grpc.NewServer(opts...)
	morris.RegisterMorrisServer(grpcServer, newMorrisServer())
	wg.Done()
	logI("SERVER STARTED")

	err = grpcServer.Serve(lis)
	if err != nil {
		logF(err.Error())
	}
}

func (s *MorrisServer) GetBoardStream(_ *morris.Empty, stream morris.Morris_GetBoardStreamServer) error {
	logI("Board stream requested.")
	streamChan := make(chan *morris.BoardState)
	s.streamChans = append(s.streamChans, &streamChan)
	stream.Send(&s.state)
	for {
		err := stream.Send(<-streamChan)
		if err != nil {
			logE(err.Error())
		}
	}
}

func (s *MorrisServer) MakeMove(ctx context.Context, move *morris.Move) (*morris.BoardState, error) {
	s.state.Board[move.From] = morris.BoardSpace_FREE
	s.state.Board[move.To] = s.state.Turn
	switch s.state.Turn {
	case morris.BoardSpace_WHITE:
		s.state.Turn = morris.BoardSpace_BLACK
	case morris.BoardSpace_BLACK:
		s.state.Turn = morris.BoardSpace_WHITE
	}
	s.Broadcast(&s.state)
	return &s.state, nil
}

func (s *MorrisServer) Broadcast(state *morris.BoardState) {
	for _, streamChan := range s.streamChans {
		*streamChan <- state
	}
}
