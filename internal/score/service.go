package score

import (
	"github.com/rs/zerolog/log"

	"main/internal/state/ctx"
	"main/internal/storage/db/repository"
)

type ScoreService struct {
	appCtx *ctx.AppContext
}

func NewScoreService(appCtx *ctx.AppContext) *ScoreService {
	return &ScoreService{appCtx}
}

func (s *ScoreService) SetScore(name string, score int64) (*repository.Score, error) {
	playerScore, err := s.appCtx.Repo.Score.GetScoreByPlayerName(name)
	if err != nil {
		return nil, err
	}

	if playerScore == nil {
		return s.createScore(name, score)
	}

	if playerScore.Score < score {
		return s.updateScore(name, score)
	}

	return playerScore, nil
}

func (s *ScoreService) createScore(name string, score int64) (*repository.Score, error) {
	if err := s.appCtx.Repo.User.Create(name); err != nil {
		return nil, err
	}

	log.Info().Str("name", name).Msg("user registered")

	if err := s.appCtx.Repo.Score.Create(name, score); err != nil {
		return nil, err
	}

	return s.retrieveScore(name)
}

func (s *ScoreService) updateScore(name string, score int64) (*repository.Score, error) {
	if err := s.appCtx.Repo.Score.Update(name, score); err != nil {
		return nil, err
	}

	return s.retrieveScore(name)
}

func (s *ScoreService) retrieveScore(name string) (*repository.Score, error) {
	return s.appCtx.Repo.Score.GetScoreByPlayerName(name)
}
