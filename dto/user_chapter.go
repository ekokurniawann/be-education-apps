package dto

import "time"

type CreateUserChapterRequest struct {
	ChapterID   int64      `json:"chapter_id" binding:"required"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	QuizScore   *float64   `json:"quiz_score,omitempty"`
}

type UserChapterQuizScoreResponse struct {
	ChapterName string     `json:"chapter_name" db:"chapter_name"`
	QuizScore   *float64   `json:"quiz_score,omitempty" db:"quiz_score"`
	CompletedAt *time.Time `json:"completed_at,omitempty" db:"completed_at"`
}

type CheckChapterCompletionRequest struct {
	ChapterID int `json:"chapter_id" binding:"required"`
}
