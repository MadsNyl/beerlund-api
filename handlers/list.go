package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"api.beerlund.com/m/logger"
	"api.beerlund.com/m/models"
)

func (h *Handler) ListEvents(w http.ResponseWriter, r *http.Request) {
	logger.Info("Request to ListEvents", map[string]any{
		"method": r.Method,
		"url":    r.URL.String(),
		"remote_addr": r.RemoteAddr,
		"usr_agent": r.UserAgent(),
	})

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	ended := r.URL.Query().Get("ended")
	if ended == "" {
		ended = "false" // Default to false if not specified
	}
	endedBool, err := strconv.ParseBool(ended)
	if err != nil {
		endedBool = false // Default to false if parsing fails
	}
	
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	logger.Info("Fetching events from db", map[string]any{
		"page": page,
		"limit": limit,
		"ended": endedBool,
		"remote_addr": r.RemoteAddr,
		"usr_agent": r.UserAgent(),
	})

	events, err := h.Store.ListEvents(page, limit, endedBool)
	if err != nil {
		logger.Error("Failed to list events: "+err.Error(), map[string]any{
			"page": page,
			"limit": limit,
			"ended": endedBool,
			"remote_addr": r.RemoteAddr,
			"usr_agent": r.UserAgent(),
		})
		http.Error(w, "Failed to list events", http.StatusInternalServerError)
		return
	}

	response := models.EventListResponse{
		Events:     events.Events,
		NextPage:   events.NextPage,
		PrevPage:   events.PrevPage,
		TotalCount: events.TotalCount,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	logger.Info("Successfully listed events", map[string]any{
		"page": page,
		"limit": limit,
		"ended": endedBool,
		"remote_addr": r.RemoteAddr,
		"usr_agent": r.UserAgent(),
	})
}
