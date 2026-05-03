package model

import "time"

type UserInfo struct {
	Name  string
	Email string
}

type User struct {
	ID           string
	UserInfo     UserInfo
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
