package models

import "time"

type EventList struct {
	ID          int    		`json:"id"`
	Name        string 		`json:"name"`
	Description string 		`json:"description"`
	Address     string 		`json:"address"`
	ZipCode     string 		`json:"zip_code"`
	City        string 		`json:"city"`
	Country     string 		`json:"country"`
	StartTime   time.Time 	`json:"start_time"`
	EndTime     time.Time 	`json:"end_time"`
	ImageURL    string 		`json:"image_url"`
	Participants int   		`json:"participants"`
}

type EventListResponse struct {
	Events []EventList 	`json:"events"`
	NextPage int    	`json:"next_page"`
	PrevPage int    	`json:"prev_page"`
	TotalCount int 		`json:"total_count"`
}

type EventResponse struct {
	ID          	int       		`json:"id"`
	Name        	string    		`json:"name"`
	Description 	string    		`json:"description"`
	Address     	string    		`json:"address"`
	ZipCode     	string    		`json:"zip_code"`
	City        	string    		`json:"city"`
	Country     	string    		`json:"country"`
	StartTime   	time.Time 		`json:"start_time"`
	EndTime     	time.Time 		`json:"end_time"`
	ImageURL    	string    		`json:"image_url"`
	Participants 	[]Participant 	`json:"participants"`
}

type EventRequest struct {
	Name        string    `json:"name" validate:"required"`
	Description string    `json:"description"`
	Address     string    `json:"address"`
	ZipCode     string    `json:"zip_code"`
	City        string    `json:"city"`
	Country     string    `json:"country"`
	StartTime   time.Time `json:"start_time" validate:"required"`
	EndTime     time.Time `json:"end_time" validate:"required"`
	ImageURL    string    `json:"image_url"`
}