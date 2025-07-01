package dto

import "time"

type CreateUserRequest struct {
	Name     string     `json:"name" binding:"required"`
	Email    string     `json:"email" binding:"required,email"`
	Password string     `json:"password" binding:"required,min=6"`
	Class    *string    `json:"class,omitempty"`
	Birthday *time.Time `json:"birthday,omitempty"`
}

type LoginUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Class      string    `json:"class"`
	Birthday   string    `json:"birthday"`
	Role       string    `json:"role"`
	ProfileURL string    `json:"profile_url"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
