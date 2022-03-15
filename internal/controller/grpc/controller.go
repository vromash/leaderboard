package grpc_controller

import (
	"context"
	"io"

	"github.com/rs/zerolog/log"

	leaderboard "main/internal/controller/protos"
)

func (s *Server) ListScore(_ context.Context, req *leaderboard.ListScoreRequest) (*leaderboard.ListScoreResponse, error) {
	return &leaderboard.ListScoreResponse{
		Results:  []*leaderboard.PlayerScore{},
		AroundMe: []*leaderboard.PlayerScore{},
		Page:     0,
	}, nil
}

func (s *Server) SaveScore(stream leaderboard.LeaderboardService_SaveScoreServer) error {
	for {
		score, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			log.Error().Err(err).Msg("failed to receive new score")
			return err
		}

		updatedScore, err := s.scoreSvc.SetScore(score.Name, score.Score)
		if err != nil {
			log.Error().Err(err).Msg("failed to save new score")
			return err
		}

		resp := &leaderboard.SaveScoreResponse{
			Name: updatedScore.Name,
			Rank: updatedScore.Rank,
		}
		if err := stream.Send(resp); err != nil {
			log.Error().Err(err).Msg("failed to send response after new score save")
			return err
		}
	}
}
