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
			Turn:        morris.BoardSpace_WHITE,
			Board:       make([]morris.BoardSpace, 24),
			WhitePieces: 9,
			BlackPieces: 9,
			WhiteGrave:  0,
			BlackGrave:  0,
			Phase:       morris.Phase_PLACE,
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
	logI("Move attempt: " + fmt.Sprint(move))
	oldState := morris.BoardState{
		Turn:        s.state.Turn,
		Board:       make([]morris.BoardSpace, 24),
		WhitePieces: s.state.WhitePieces,
		BlackPieces: s.state.BlackPieces,
		WhiteGrave:  s.state.WhiteGrave,
		BlackGrave:  s.state.BlackGrave,
		Phase:       s.state.Phase,
	}
	copy(oldState.Board, s.state.Board)
	var (
		ownPieces *int32
		//othersPieces *int32
		//ownGrave     *int32
		othersGrave *int32
	)
	switch s.state.Turn {
	case morris.BoardSpace_WHITE:
		ownPieces = &s.state.WhitePieces
		//othersPieces = &s.state.BlackPieces
		//ownGrave = &s.state.WhiteGrave
		othersGrave = &s.state.BlackGrave
	case morris.BoardSpace_BLACK:
		ownPieces = &s.state.BlackPieces
		//othersPieces = &s.state.WhitePieces
		//ownGrave = &s.state.BlackGrave
		othersGrave = &s.state.WhiteGrave
	}

	switch s.state.Phase {
	case morris.Phase_PLACE:
		// Place in space if the space is free
		if s.state.Board[move.To] == morris.BoardSpace_FREE {
			s.state.Board[move.To] = s.state.Turn
			*ownPieces--
			err := s.handleMorris(move.To, move.Remove, othersGrave, s.state.Turn)
			if err != nil {
				s.state = oldState
				logE(err.Error())
				return nil, err
			}
		} else {
			err := morris.IllegalMoveError{Description: "Piece must be placed on an empty space."}
			logE(err.Error())
			return nil, err
		}
	case morris.Phase_MOVE:
	case morris.Phase_FLY:
	}

	// Update turn
	switch s.state.Turn {
	case morris.BoardSpace_WHITE:
		s.state.Turn = morris.BoardSpace_BLACK
		logI("Turn: BLACK")
		if s.state.BlackPieces > 0 {
			s.state.Phase = morris.Phase_PLACE
			logI("Phase: PLACE")
		} else if s.state.BlackGrave < 3 {
			s.state.Phase = morris.Phase_MOVE
			logI("Phase: MOVE")
		} else {
			s.state.Phase = morris.Phase_FLY
			logI("Phase: FLY")
		}
	case morris.BoardSpace_BLACK:
		s.state.Turn = morris.BoardSpace_WHITE
		logI("Turn: WHITE")
		if s.state.WhitePieces > 0 {
			s.state.Phase = morris.Phase_PLACE
			logI("Phase: PLACE")
		} else if s.state.WhiteGrave < 3 {
			s.state.Phase = morris.Phase_MOVE
			logI("Phase: MOVE")
		} else {
			s.state.Phase = morris.Phase_FLY
			logI("Phase: FLY")
		}
	}
	fmt.Print(s.state.Visualize(false))
	logI(fmt.Sprint(s.state))

	// Broadcast
	s.Broadcast(&s.state)
	return &s.state, nil
}

func (s *MorrisServer) handleMorris(space int32, remove int32, othersGrave *int32, ownColor morris.BoardSpace) error {
	if s.state.HasMorrisAt(space) {
		// Remove piece, but only if it is not in a morris
		if s.state.Board[remove] == ownColor {
			return morris.IllegalMoveError{Description: "You cannot remove your own pieces."}
		}
		if s.state.Board[remove] == morris.BoardSpace_FREE {
			return morris.IllegalMoveError{Description: "You cannot remove a piece from a free space."}
		}
		if s.state.HasMorrisAt(remove) {
			return morris.IllegalMoveError{Description: "You cannot remove a piece from an enemy morris."}
		}
		s.state.Board[remove] = morris.BoardSpace_FREE
		*othersGrave++
	}
	return nil
}

func (s *MorrisServer) Broadcast(state *morris.BoardState) {
	for _, streamChan := range s.streamChans {
		*streamChan <- state
	}
}
