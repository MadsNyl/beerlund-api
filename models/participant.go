package models

import "github.com/clerk/clerk-sdk-go/v2"

type Participant struct {
	ID        int         `json:"id"`
	EventID   int         `json:"event_id"`
	UserID    string      `json:"user_id"`
	User      *clerk.User `json:"user"`
	CreatedAt string      `json:"created_at"`
	UpdatedAt string      `json:"updated_at"`
}