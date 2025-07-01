package dto

type StudentCount struct {
	Class string `db:"class"`
	Total int    `db:"total"`
}
