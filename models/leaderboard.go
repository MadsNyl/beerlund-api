package models

import (
	"github.com/clerk/clerk-sdk-go/v2"
)

type EventLeaderboard struct {
	ID            int  		`json:"id"`
	EventId       int  		`json:"event_id"`
	ParticipantId string  	`json:"participant_id"`
	User 		  *clerk.User `json:"user"`
	Score         int  		`json:"score"`
	Rank		  int  		`json:"rank"`
}