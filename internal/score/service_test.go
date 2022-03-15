package score

import (
	"testing"

	"zgo.at/zcache"

	"github.com/go-test/deep"

	"main/internal/state/ctx"
	"main/internal/storage/db/repository"
)

type mockUserRepo struct {
	repository.UserRepository
	users *zcache.Cache
}

func (m *mockUserRepo) Create(name string) error {
	m.users.SetDefault(name, "")
	return nil
}

type mockScoreRepo struct {
	repository.ScoreRepository
	scores *zcache.Cache
}

func (m *mockScoreRepo) Create(name string, score int64) error {
	m.scores.SetDefault(name, score)
	return nil
}

func (m *mockScoreRepo) Update(name string, score int64) error {
	m.scores.SetDefault(name, score)
	return nil
}

func (m *mockScoreRepo) GetScoreByPlayerName(name string) (*repository.Score, error) {
	all := m.scores.Items()
	score, exists := m.scores.Get(name)

	if !exists {
		return &repository.Score{}, nil
	}

	var rank int64 = 1
	for _, sc := range all {
		if score.(int64) < sc.Object.(int64) {
			rank++
		}
	}

	return &repository.Score{
		Name:  name,
		Score: score.(int64),
		Rank:  rank,
	}, nil
}

func (m *mockScoreRepo) getScores() map[string]int64 {
	items := m.scores.Items()
	result := map[string]int64{}
	for name, item := range items {
		result[name] = item.Object.(int64)
	}
	return result
}

func getMockStructs() (*ctx.AppContext, *mockUserRepo, *mockScoreRepo) {
	userRepo := mockUserRepo{
		users: zcache.New(zcache.NoExpiration, zcache.NoExpiration),
	}

	scoreRepo := mockScoreRepo{
		scores: zcache.New(zcache.NoExpiration, zcache.NoExpiration),
	}

	return &ctx.AppContext{
		Repo: &repository.RepositoryProvider{
			User:  &userRepo,
			Score: &scoreRepo,
		},
	}, &userRepo, &scoreRepo
}

func TestScoreService_SetScore(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name            string
		expectedScores  map[string]int64
		expectedOutcome []*repository.Score
		scores          []struct {
			name  string
			score int64
		}
	}{
		{
			name: "case1",
			expectedScores: map[string]int64{
				"user1": 110,
				"user2": 102,
				"user3": 102,
				"user4": 101,
				"user5": 105,
			},
			expectedOutcome: []*repository.Score{
				{
					Name:  "user1",
					Score: 100,
					Rank:  1,
				},
				{
					Name:  "user2",
					Score: 102,
					Rank:  1,
				},
				{
					Name:  "user1",
					Score: 103,
					Rank:  1,
				},
				{
					Name:  "user3",
					Score: 100,
					Rank:  3,
				},
				{
					Name:  "user4",
					Score: 101,
					Rank:  3,
				},
				{
					Name:  "user2",
					Score: 102,
					Rank:  2,
				},
				{
					Name:  "user3",
					Score: 102,
					Rank:  2,
				},
				{
					Name:  "user5",
					Score: 105,
					Rank:  1,
				},
				{
					Name:  "user1",
					Score: 110,
					Rank:  1,
				},
				{
					Name:  "user5",
					Score: 105,
					Rank:  2,
				},
				{
					Name:  "user5",
					Score: 105,
					Rank:  2,
				},
			},
			scores: []struct {
				name  string
				score int64
			}{
				{"user1", 100},
				{"user2", 102},
				{"user1", 103},
				{"user3", 100},
				{"user4", 101},
				{"user2", 100},
				{"user3", 102},
				{"user5", 105},
				{"user1", 110},
				{"user5", 102},
				{"user5", 101},
			},
		},
	}

	for _, tcase := range testCases {
		tcScoped := tcase
		t.Run(tcScoped.name, func(t *testing.T) {
			t.Parallel()

			appCtx, _, scoreRepo := getMockStructs()
			svc := NewScoreService(appCtx)

			for i, s := range tcScoped.scores {
				out, err := svc.SetScore(s.name, s.score)
				if err != nil {
					t.Errorf("Got an error while setting scores. %v", err)
				}

				if diff := deep.Equal(out, tcScoped.expectedOutcome[i]); diff != nil {
					t.Errorf("SetScore outcome is different from expected: %v", diff)
				}
			}

			if diff := deep.Equal(scoreRepo.getScores(), tcScoped.expectedScores); diff != nil {
				t.Errorf("Final scores are different from expected: %v", diff)
			}
		})
	}
}
