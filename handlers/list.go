package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"api.beerlund.com/m/models"
)

func (h *Handler) ListEvents(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	events, err := h.Store.ListEvents(page, limit)
	if err != nil {
		log.Printf("Error listing events: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
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
}
