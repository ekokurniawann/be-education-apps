package models

import "time"

type User struct {
	ID         int64      `json:"id" db:"id"`
	Name       string     `json:"name" db:"name"`
	Email      string     `json:"email" db:"email"`
	Password   string     `json:"password" db:"password"`
	Class      *string    `json:"class,omitempty" db:"class"`
	Birthday   *time.Time `json:"birthday,omitempty" db:"birthday"`
	Role       string     `json:"role" db:"role"`
	ProfileURL *string    `json:"profile_url,omitempty" db:"profile_url"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at" db:"updated_at"`
}
