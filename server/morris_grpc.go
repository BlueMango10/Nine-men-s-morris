package main

import (
	"context"
	"fmt"
	"net"

	"github.com/BlueMango10/Nine-men-s-morris/morris"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MorrisServer struct {
	morris.UnimplementedMorrisServer
	// Data goes here
}

func newMorrisServer() *MorrisServer {
	return &MorrisServer{}
}

func startServer(address string) {
	logI(fmt.Sprintf("STARTING SERVER WITH ADDRESS: %s", address))
	lis, err := net.Listen("tcp", address)
	if err != nil {
		logF(err.Error())
	}
	var opts []grpc.ServerOption
	// Server options here (none yet)

	grpcServer := grpc.NewServer(opts...)
	morris.RegisterMorrisServer(grpcServer, newMorrisServer())
	err = grpcServer.Serve(lis)

	if err != nil {
		logF(err.Error())
	}
	logI("SERVER STARTED")
}

func (s *MorrisServer) GetBoardStream(*morris.Empty, morris.Morris_GetBoardStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method GetBoardStream not implemented")
}

func (s *MorrisServer) MakeMove(context.Context, *morris.Move) (*morris.BoardState, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MakeMove not implemented")
}
