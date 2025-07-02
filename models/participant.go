package models

type Participant struct {
	ID        int    `json:"id"`
	EventID   int    `json:"event_id"`
	UserID    string    `json:"user_id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}