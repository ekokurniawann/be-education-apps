package dto

// UserChapterScore represents a user's score for a specific chapter.
// This is used for combining user information with their chapter scores.
type UserChapterScore struct {
	UserID    int64   `db:"user_id" json:"userId"`
	UserName  string  `db:"user_name" json:"userName"`
	UserClass string  `db:"user_class" json:"userClass"`
	ChapterID int64   `db:"chapter_id" json:"chapterId"`
	Score     float64 `db:"score" json:"score"`
}

// UserScoreEntry represents a single user's summary for chapter scores.
// This struct is specifically designed to be an element within the UserChapterScoresSummary.
type UserScoreEntry struct {
	ID            int64              `json:"id"`
	Name          string             `json:"name"`
	Class         string             `json:"class"`
	ChapterScores map[string]float64 `json:"chapterScores"`
}

// UserChapterScoresSummary represents the summary of all users with their chapter scores.
type UserChapterScoresSummary struct {
	UsersScores []UserScoreEntry `json:"usersScores"` // Now uses the named UserScoreEntry
}
