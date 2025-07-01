package models

import "time"

type UserChapter struct {
	ID          int64      `json:"id" db:"id"`
	UserID      int64      `json:"user_id" db:"user_id"`
	ChapterID   int64      `json:"chapter_id" db:"chapter_id"`
	CompletedAt *time.Time `json:"completed_at,omitempty" db:"completed_at"`
	QuizScore   *float64   `json:"quiz_score,omitempty" db:"quiz_score"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}
