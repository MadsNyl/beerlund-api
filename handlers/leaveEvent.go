package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/user"
)

type leaveEventRequest struct {
	EventID int `json:"event_id"`
}

func (h *Handler) LeaveEvent(w http.ResponseWriter, r *http.Request) {
	claims, ok := clerk.SessionClaimsFromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"access": "unauthorized"}`))
		return
	}

	usr, err := user.Get(r.Context(), claims.Subject)
	if err != nil {
		http.Error(w, fmt.Sprintf("User not found: %v", err), http.StatusUnauthorized)
		return
	}

	userID := usr.ID

	var req leaveEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if req.EventID <= 0 {
		http.Error(w, "event_id must be a positive integer", http.StatusBadRequest)
		return
	}

	err = h.Store.RemoveParticipant(req.EventID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "You are not participating in this event", http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("Failed to unassign from event: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}