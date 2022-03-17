package score

import (
	"sort"
	"testing"
	"time"

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
	scores              *zcache.Cache
	scoresForListing    *zcache.Cache
	requestMonthDate    time.Time
	scoreMonthDate      time.Time
	requsestAllTimeDate time.Time
	scoreAllTimeDate    time.Time
}

type score struct {
	score       int64
	rankAllTime int64
	rankMonth   int64
	updatedAt   time.Time
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

func (m *mockScoreRepo) GetInRange(from, to int64, maxDate time.Time) ([]*repository.Score, error) {
	scores := m.scoresForListing.Items()

	var result []*repository.Score
	for user, sc := range scores {
		s := sc.Object.(*score)

		rank := s.rankMonth
		if maxDate == m.requsestAllTimeDate {
			rank = s.rankAllTime
		}

		if rank < from ||
			rank > to {
			continue
		}

		if s.updatedAt.After(maxDate) {
			result = append(result, &repository.Score{
				Name:  user,
				Score: s.score,
				Rank:  rank,
			})
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Rank < result[j].Rank
	})

	return result, nil
}

func (m *mockScoreRepo) GetScoreByPlayerNameInTimeRange(name string, maxDate time.Time) (*repository.Score, error) {
	sc, ok := m.scoresForListing.Get(name)
	if !ok {
		return nil, nil
	}

	if sc.(*score).updatedAt.Before(maxDate) {
		return nil, nil
	}

	rank := sc.(*score).rankMonth
	if maxDate == m.requsestAllTimeDate {
		rank = sc.(*score).rankAllTime
	}

	return &repository.Score{
		Name:  name,
		Score: sc.(*score).score,
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

	requestMonthDate := time.Date(2022, 3, 1, 10, 10, 10, 10, time.UTC)
	scoreMonthDate := requestMonthDate.Add(24 * time.Hour)
	requestAllTimeDate := requestMonthDate.Add(-24 * 30 * time.Hour)
	scoreAllTimeDate := requestAllTimeDate.Add(24 * time.Hour)

	scoreRepo := mockScoreRepo{
		scores:              zcache.New(zcache.NoExpiration, zcache.NoExpiration),
		scoresForListing:    zcache.New(zcache.NoExpiration, zcache.NoExpiration),
		requestMonthDate:    requestMonthDate,
		scoreMonthDate:      scoreMonthDate,
		requsestAllTimeDate: requestAllTimeDate,
		scoreAllTimeDate:    scoreAllTimeDate,
	}

	listingData := map[string]*score{
		"user_1":  {110, 8, 4, scoreMonthDate},
		"user_2":  {50, 17, 30, scoreAllTimeDate},
		"user_3":  {201, 2, 2, scoreMonthDate},
		"user_4":  {165, 4, 30, scoreAllTimeDate},
		"user_5":  {82, 14, 9, scoreMonthDate},
		"user_6":  {100, 9, 5, scoreMonthDate},
		"user_7":  {89, 11, 7, scoreMonthDate},
		"user_8":  {169, 3, 30, scoreAllTimeDate},
		"user_9":  {15, 20, 11, scoreMonthDate},
		"user_10": {150, 5, 30, scoreAllTimeDate},
		"user_11": {120, 7, 3, scoreMonthDate},
		"user_12": {85, 13, 30, scoreAllTimeDate},
		"user_13": {221, 1, 1, scoreMonthDate},
		"user_14": {68, 16, 30, scoreAllTimeDate},
		"user_15": {88, 12, 8, scoreMonthDate},
		"user_16": {46, 18, 30, scoreAllTimeDate},
		"user_17": {95, 10, 6, scoreMonthDate},
		"user_18": {77, 15, 30, scoreAllTimeDate},
		"user_19": {20, 19, 10, scoreMonthDate},
		"user_20": {123, 6, 30, scoreAllTimeDate},
	}

	for user, sc := range listingData {
		scoreRepo.scoresForListing.SetDefault(user, sc)
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

func TestScoreService_ListScores(t *testing.T) {
	t.Parallel()

	appCtx, _, scoreRepo := getMockStructs()
	svc := NewScoreService(appCtx)

	testCases := []struct {
		name     string
		expected *ScoreList
		user     string
		page     int64
		maxDate  time.Time
	}{
		{
			name: "no user, page 1, no period",
			expected: &ScoreList{
				PagedScore: []*repository.Score{
					{"user_13", 221, 1},
					{"user_3", 201, 2},
					{"user_11", 120, 3},
					{"user_1", 110, 4},
					{"user_6", 100, 5},
					{"user_17", 95, 6},
					{"user_7", 89, 7},
					{"user_15", 88, 8},
					{"user_5", 82, 9},
					{"user_19", 20, 10},
				},
				ScoreAroundPlayer: nil,
			},
			user:    "",
			page:    1,
			maxDate: scoreRepo.requestMonthDate,
		},
		{
			name: "no user, page 2, no period",
			expected: &ScoreList{
				PagedScore: []*repository.Score{
					{"user_9", 15, 11},
				},
				ScoreAroundPlayer: nil,
			},
			user:    "",
			page:    2,
			maxDate: scoreRepo.requestMonthDate,
		},
		{
			name: "user_9, page 1, no period",
			expected: &ScoreList{
				PagedScore: []*repository.Score{
					{"user_13", 221, 1},
					{"user_3", 201, 2},
					{"user_11", 120, 3},
					{"user_1", 110, 4},
					{"user_6", 100, 5},
					{"user_17", 95, 6},
					{"user_7", 89, 7},
					{"user_15", 88, 8},
					{"user_5", 82, 9},
					{"user_19", 20, 10},
				},
				ScoreAroundPlayer: []*repository.Score{
					{"user_5", 82, 9},
					{"user_19", 20, 10},
					{"user_9", 15, 11},
				},
			},
			user:    "user_9",
			page:    1,
			maxDate: scoreRepo.requestMonthDate,
		},
		{
			name: "user_9, page 2, no period",
			expected: &ScoreList{
				PagedScore: []*repository.Score{
					{"user_9", 15, 11},
				},
				ScoreAroundPlayer: nil,
			},
			user:    "user_9",
			page:    2,
			maxDate: scoreRepo.requestMonthDate,
		},
		{
			name: "no user, page 1, all time",
			expected: &ScoreList{
				PagedScore: []*repository.Score{
					{"user_13", 221, 1},
					{"user_3", 201, 2},
					{"user_8", 169, 3},
					{"user_4", 165, 4},
					{"user_10", 150, 5},
					{"user_20", 123, 6},
					{"user_11", 120, 7},
					{"user_1", 110, 8},
					{"user_6", 100, 9},
					{"user_17", 95, 10}},
				ScoreAroundPlayer: nil,
			},
			user:    "",
			page:    1,
			maxDate: scoreRepo.requsestAllTimeDate,
		},
		{
			name: "user_18, page 1, all time",
			expected: &ScoreList{
				PagedScore: []*repository.Score{
					{"user_13", 221, 1},
					{"user_3", 201, 2},
					{"user_8", 169, 3},
					{"user_4", 165, 4},
					{"user_10", 150, 5},
					{"user_20", 123, 6},
					{"user_11", 120, 7},
					{"user_1", 110, 8},
					{"user_6", 100, 9},
					{"user_17", 95, 10},
				},
				ScoreAroundPlayer: []*repository.Score{
					{"user_12", 85, 13},
					{"user_5", 82, 14},
					{"user_18", 77, 15},
					{"user_14", 68, 16},
					{"user_2", 50, 17},
				},
			},
			user:    "user_18",
			page:    1,
			maxDate: scoreRepo.requsestAllTimeDate,
		},
	}

	for _, tcase := range testCases {
		tcScoped := tcase
		t.Run(tcScoped.name, func(t *testing.T) {
			t.Parallel()

			var user *string
			if tcScoped.user != "" {
				user = &tcScoped.user
			}
			listedScores, err := svc.ListScores(user, tcScoped.page, tcScoped.maxDate)
			if err != nil {
				t.Errorf("Got an error while listing scores. %v", err)
			}

			if diff := deep.Equal(listedScores, tcScoped.expected); diff != nil {
				t.Errorf("Listed scores are different from expected: %v", diff)
			}
		})
	}
}
