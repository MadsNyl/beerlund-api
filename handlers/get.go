package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func (h *Handler) GetEvent(w http.ResponseWriter, r *http.Request) {
    parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/events/"), "/")
    if len(parts) == 0 || parts[0] == "" {
        http.Error(w, "Event ID is required", http.StatusBadRequest)
        return
    }
    idStr := parts[0]

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid Event ID", http.StatusBadRequest)
		return
	}

	event, err := h.Store.GetEvent(id)
	if err != nil {
		log.Printf("Error fetching event with ID %d: %v", id, err)
		http.Error(w, "Event not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(event); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}