package main

import (
	"context"

	"github.com/BlueMango10/Nine-men-s-morris/morris"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MorrisServer struct {
	morris.UnimplementedMorrisServer
	// Data goes here
}

func (s *MorrisServer) GetBoardStream(*morris.Empty, morris.Morris_GetBoardStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method GetBoardStream not implemented")
}
func (s *MorrisServer) MakeMove(context.Context, *morris.Move) (*morris.BoardState, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MakeMove not implemented")
}
