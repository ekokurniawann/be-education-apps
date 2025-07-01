package dto

type AdminSummary struct {
	TotalAdmins int            `json:"totalAdmins"`
	Admins      []UserResponse `json:"admins"`
}
