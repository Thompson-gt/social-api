package types

import (
	"time"
)

type RequestUser struct {
	UserID       string    `json:"userId"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	ProfilePic   string    `json:"profilePicture"`
	CoverPic     string    `json:"coverPicture"`
	Follwers     []string  `json:"follwers"`
	Follwings    []string  `json:"follwings"`
	IsAdmin      bool      `json:"isAdmin"`
	Desc         string    `json:"desc"`
	City         string    `json:"city"`
	From         string    `json:"from"`
	Relationship int       `json:"relationship"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"` // need to update this whenever changing data
}
