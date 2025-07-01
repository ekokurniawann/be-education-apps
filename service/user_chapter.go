package service

import (
	"be-education/dto"
	"be-education/models"
	"be-education/repository"
	"context"
	"fmt"
	"sort"
)

type UserChapterService interface {
	CreateUserChapter(ctx context.Context, userChapter *models.UserChapter) error
	GetUserQuizScoresByUserID(ctx context.Context, userID int64) ([]*dto.UserChapterQuizScoreResponse, error)
	CheckUserChapterCompleted(ctx context.Context, userID int64, chapterID int) (bool, error)
	GetAllUsersChapterScoresSummary(ctx context.Context) (*dto.UserChapterScoresSummary, error)
}

type userChapterServiceImpl struct {
	userChapterRepo repository.UserChapterRepository
}

func NewUserChapterService(userChapterRepo repository.UserChapterRepository) UserChapterService {
	return &userChapterServiceImpl{userChapterRepo: userChapterRepo}
}

func (s *userChapterServiceImpl) CreateUserChapter(ctx context.Context, userChapter *models.UserChapter) error {

	err := s.userChapterRepo.CreateUserChapter(ctx, userChapter)
	if err != nil {
		return fmt.Errorf("service failed to create user chapter: %w", err)
	}
	return nil
}

func (s *userChapterServiceImpl) GetUserQuizScoresByUserID(ctx context.Context, userID int64) ([]*dto.UserChapterQuizScoreResponse, error) {
	quizScores, err := s.userChapterRepo.GetUserQuizScoresByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("service failed to get user quiz scores: %w", err)
	}
	return quizScores, nil
}

func (s *userChapterServiceImpl) CheckUserChapterCompleted(ctx context.Context, userID int64, chapterID int) (bool, error) {
	completed, err := s.userChapterRepo.CheckUserChapterCompletion(ctx, userID, chapterID)
	if err != nil {
		return false, fmt.Errorf("service failed to check user chapter completion: %w", err)
	}
	return completed, nil
}

func (s *userChapterServiceImpl) GetAllUsersChapterScoresSummary(ctx context.Context) (*dto.UserChapterScoresSummary, error) {
	rawScores, err := s.userChapterRepo.GetAllUsersWithAllChapterScores(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get raw chapter scores from repository: %w", err)
	}

	userScoresMap := make(map[int64]dto.UserScoreEntry)

	for _, rs := range rawScores {
		if _, ok := userScoresMap[rs.UserID]; !ok {
			userScoresMap[rs.UserID] = dto.UserScoreEntry{
				ID:            rs.UserID,
				Name:          rs.UserName,
				Class:         rs.UserClass,
				ChapterScores: make(map[string]float64),
			}
		}
		userEntry := userScoresMap[rs.UserID]

		if rs.ChapterID > 0 {
			chapterKey := fmt.Sprintf("C%d", rs.ChapterID) // Format as "C1", "C2", etc.
			userEntry.ChapterScores[chapterKey] = rs.Score
		}
		userScoresMap[rs.UserID] = userEntry
	}

	var usersScores []dto.UserScoreEntry
	for _, userEntry := range userScoresMap {
		usersScores = append(usersScores, userEntry)
	}

	sort.Slice(usersScores, func(i, j int) bool {
		return usersScores[i].ID < usersScores[j].ID
	})

	summary := &dto.UserChapterScoresSummary{
		UsersScores: usersScores,
	}

	return summary, nil
}
