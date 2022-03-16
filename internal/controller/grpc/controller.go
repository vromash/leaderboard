package grpc_controller

import (
	"context"
	"fmt"
	"io"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	leaderboard "main/internal/controller/protos"
	"main/internal/storage/db/repository"
)

func (s *Server) ListScore(_ context.Context, req *leaderboard.ListScoreRequest) (*leaderboard.ListScoreResponse, error) {
	useAllTimeRecords := s.shouldUseAllTimeRecords(req.Period)
	page, nextPage, err := s.getPageOptions(req.Page, useAllTimeRecords)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	scores, err := s.scoreSvc.ListScores(req.Name, page)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &leaderboard.ListScoreResponse{
		Results:  s.convertRepoScoreIntoLeaderboardScore(scores.PagedScore),
		AroundMe: s.convertRepoScoreIntoLeaderboardScore(scores.ScoreAroundPlayer),
		Page:     nextPage,
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
			return status.Error(codes.Internal, err.Error())
		}

		updatedScore, err := s.scoreSvc.SetScore(score.Name, score.Score)
		if err != nil {
			log.Error().Err(err).Msg("failed to save new score")
			return status.Error(codes.Internal, err.Error())
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

func (s *Server) shouldUseAllTimeRecords(period *leaderboard.TimePeriod) bool {
	return period != nil && *period == leaderboard.TimePeriod_TIME_PERIOD_ALL
}

// getPageOptions returns (currentPage, maxPage, error)
func (s *Server) getPageOptions(page *int64, allTime bool) (int64, int64, error) {
	maxPage, err := s.scoreSvc.GetMaxPage(allTime)
	if err != nil {
		log.Error().Err(err).Msg("failed to get page options")
		return 0, 0, status.Error(codes.Internal, "internal error")
	}

	currentPage, err := s.validatePage(page, maxPage)
	if err != nil {
		return 0, 0, status.Error(codes.InvalidArgument, err.Error())
	}

	nextPage := s.getNextPage(currentPage, maxPage)

	return currentPage, nextPage, nil
}

func (s *Server) validatePage(page *int64, maxPage int64) (int64, error) {
	if page == nil {
		return 1, nil
	}

	if *page <= 0 {
		return 0, fmt.Errorf("page can't be less than 1")
	}

	if *page > maxPage {
		return 0, fmt.Errorf("page doesn't exist")
	}

	return *page, nil
}

func (s *Server) getNextPage(page, maxPage int64) int64 {
	if page < maxPage {
		return page + 1
	}
	return 0
}

func (s *Server) convertRepoScoreIntoLeaderboardScore(
	data []*repository.Score,
) []*leaderboard.PlayerScore {
	result := make([]*leaderboard.PlayerScore, len(data))
	for i, sc := range data {
		result[i] = &leaderboard.PlayerScore{
			Name:  sc.Name,
			Score: sc.Score,
			Rank:  sc.Rank,
		}
	}
	return result
}
