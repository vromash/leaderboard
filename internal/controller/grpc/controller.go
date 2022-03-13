package grpc_controller

import (
	"context"

	leaderboard "main/internal/controller/protos"
)

func (s *Server) ListScore(_ context.Context, req *leaderboard.ListScoreRequest) (*leaderboard.ListScoreResponse, error) {
	return &leaderboard.ListScoreResponse{
		Results:  []*leaderboard.PlayerScore{},
		AroundMe: []*leaderboard.PlayerScore{},
		Page:     0,
	}, nil
}

func (s *Server) SaveScore(leaderboard.LeaderboardService_SaveScoreServer) error {
	return nil
}
