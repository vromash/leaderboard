package score

import (
	"fmt"
	"math"
	"time"

	"github.com/rs/zerolog/log"

	"main/internal/state/config"
	"main/internal/state/ctx"
	"main/internal/storage/db/repository"
)

type ScoreService struct {
	appCtx *ctx.AppContext
}

type ScoreList struct {
	PagedScore        []*repository.Score
	ScoreAroundPlayer []*repository.Score
}

func NewScoreService(appCtx *ctx.AppContext) *ScoreService {
	return &ScoreService{appCtx}
}

func (s *ScoreService) ListScores(
	name *string,
	page int64,
	maxDate time.Time,
) (*ScoreList, error) {
	rankFrom, rankTo := s.getScoreRange(page)

	playerScore, err := s.listPlayerScore(name, maxDate)
	if err != nil {
		return nil, err
	}

	scoresAtPage, err := s.appCtx.Repo.Score.GetInRange(rankFrom, rankTo, maxDate)
	if err != nil {
		return nil, err
	}

	scoreAround, err := s.getAroundPlayer(playerScore, rankFrom, rankTo, maxDate)
	if err != nil {
		return nil, err
	}

	return &ScoreList{scoresAtPage, scoreAround}, nil
}

func (s *ScoreService) listPlayerScore(name *string, maxDate time.Time) (*repository.Score, error) {
	if name == nil || *name == "" {
		return nil, nil
	}

	playerScore, err := s.appCtx.Repo.Score.GetScoreByPlayerNameInTimeRange(*name, maxDate)
	if err != nil {
		return nil, err
	}

	if playerScore == nil {
		return nil, fmt.Errorf("user doesn't exist")
	}

	return playerScore, nil
}

func (s *ScoreService) getAroundPlayer(
	playerScore *repository.Score,
	rankFrom, rankTo int64,
	maxDate time.Time,
) ([]*repository.Score, error) {
	if playerScore == nil {
		return nil, nil
	}

	if playerScore.Rank < rankFrom {
		return nil, nil
	}

	if playerScore.Rank >= rankFrom &&
		playerScore.Rank <= rankTo {
		return nil, nil
	}

	aroundFrom := playerScore.Rank - 2
	aroundTo := playerScore.Rank + 2

	scoresAtPage, err := s.appCtx.Repo.Score.GetInRange(aroundFrom, aroundTo, maxDate)
	if err != nil {
		return nil, err
	}

	return scoresAtPage, nil
}

func (s *ScoreService) GetMaxPage(maxDate time.Time) (int64, error) {
	recordNumber, err := s.appCtx.Repo.Score.GetRecordNumber(maxDate)
	if err != nil {
		return 0, fmt.Errorf("failed to get max page: %w", err)
	}

	resultsPerPage := config.GetResultsPerPage()
	maxPage := float64(recordNumber) / float64(resultsPerPage)
	return int64(math.Ceil(maxPage)), nil
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

func (s *ScoreService) getScoreRange(page int64) (from, to int64) {
	resultsPerPage := config.GetResultsPerPage()

	if page == 1 {
		return 1, resultsPerPage
	}

	from = resultsPerPage*(page-1) + 1
	to = from + resultsPerPage - 1
	return
}
