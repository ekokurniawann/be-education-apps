package repository

import (
	"be-education/dto"
	"be-education/models"
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type UserChapterRepository interface {
	CreateUserChapter(ctx context.Context, userChapter *models.UserChapter) error
	GetUserQuizScoresByUserID(ctx context.Context, userID int64) ([]*dto.UserChapterQuizScoreResponse, error)
	CheckUserChapterCompletion(ctx context.Context, userID int64, chapterID int) (bool, error)
	GetAllUsersWithAllChapterScores(ctx context.Context) ([]*dto.UserChapterScore, error)
}

type userChapterImpl struct {
	db *sqlx.DB
}

func NewUserChapterRepository(db *sqlx.DB) UserChapterRepository {
	return &userChapterImpl{db: db}
}

func (r *userChapterImpl) CheckUserChapterCompletion(ctx context.Context, userID int64, chapterID int) (bool, error) {
	query := `
		SELECT COUNT(id) FROM user_chapters
		WHERE user_id = $1 AND chapter_id = $2`

	var count int
	err := r.db.GetContext(ctx, &count, query, userID, chapterID)
	if err != nil {
		return false, fmt.Errorf("failed to check user chapter completion: %w", err)
	}
	return count > 0, nil
}

func (r *userChapterImpl) CreateUserChapter(ctx context.Context, userChapter *models.UserChapter) error {
	query := `
		INSERT INTO user_chapters (user_id, chapter_id, completed_at, quiz_score, created_at, updated_at)
		VALUES (:user_id, :chapter_id, :completed_at, :quiz_score, :created_at, :updated_at)
		RETURNING id, created_at, updated_at`

	userChapter.CreatedAt = time.Now()
	userChapter.UpdatedAt = time.Now()

	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare named query for user chapter creation: %w", err)
	}
	defer stmt.Close()

	err = stmt.GetContext(ctx, userChapter, userChapter)
	if err != nil {
		return fmt.Errorf("failed to create user chapter: %w", err)
	}
	return nil
}

func (r *userChapterImpl) GetUserQuizScoresByUserID(ctx context.Context, userID int64) ([]*dto.UserChapterQuizScoreResponse, error) {
	query := `
		SELECT
			c.name AS chapter_name,
			uc.quiz_score,
			uc.completed_at
		FROM
			user_chapters uc
		JOIN
			chapters c ON uc.chapter_id = c.id
		WHERE
			uc.user_id = $1`

	var quizScores []*dto.UserChapterQuizScoreResponse
	err := r.db.SelectContext(ctx, &quizScores, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user quiz scores: %w", err)
	}

	return quizScores, nil
}

func (r *userChapterImpl) GetAllUsersWithAllChapterScores(ctx context.Context) ([]*dto.UserChapterScore, error) {
	query := `
        SELECT
            u.id as user_id,
            u.name as user_name,
            COALESCE(u.class, '') as user_class, -- Handle nullable class
            COALESCE(uc.chapter_id, 0) as chapter_id, -- COALESCE chapter_id to 0 if NULL
            COALESCE(uc.quiz_score, 0.0) as score -- Use COALESCE to get 0.0 if quiz_score is NULL
        FROM
            users u
        LEFT JOIN
            user_chapters uc ON u.id = uc.user_id
        WHERE
            u.role = 'mahasiswa' -- Assuming we only care about 'mahasiswa' for this summary
        ORDER BY
            u.id, uc.chapter_id
    `

	var results []*dto.UserChapterScore
	err := r.db.SelectContext(ctx, &results, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all users with chapter scores: %w", err)
	}
	return results, nil
}
