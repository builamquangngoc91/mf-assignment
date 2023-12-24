package models

import "time"

type User struct {
	UserID    string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
