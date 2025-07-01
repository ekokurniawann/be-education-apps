package dto

type StudentSummary struct {
	TotalStudents int            `json:"totalStudents"`
	ClassCounts   map[string]int `json:"classCounts"`
}
