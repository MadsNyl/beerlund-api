package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/user"
)

type joinEventRequest struct {
    EventID int `json:"event_id"`
}

func (h *Handler) JoinEvent(w http.ResponseWriter, r *http.Request) {
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

    // 2) decode JSON body
    var req joinEventRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }
    if req.EventID <= 0 {
        http.Error(w, "event_id must be a positive integer", http.StatusBadRequest)
        return
    }

    // Check if the user is already participating in the event
    isParticipating, err := h.Store.IsParticipating(req.EventID, userID)

    if err != nil {
        http.Error(w, "Failed to check participation status", http.StatusInternalServerError)
        return
    }

    if isParticipating {
        http.Error(w, "You are already participating in this event", http.StatusConflict)
        return
    }

    // 3) load the event’s max_participants
    maxParticipants, err := h.Store.GetMaxParticipants(req.EventID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Event not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

    // 4) count current participants
    participantCount, err := h.Store.CountParticipants(req.EventID)
	if err != nil {
		http.Error(w, "Failed to count participants", http.StatusInternalServerError)
		return
	}

    // 5) check capacity (0 for unlimited)
    if maxParticipants > 0 && participantCount >= maxParticipants {
        http.Error(w, "Event is full", http.StatusForbidden)
        return
    }

	// 6) insert participant
	if err := h.Store.InsertParticipant(req.EventID, userID); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Event not found", http.StatusNotFound)
			return
		}
		if err.Error() == "duplicate key value violates unique constraint" {
			http.Error(w, "You are already participating in this event", http.StatusConflict)
			return
		}
		http.Error(w, "Failed to join event", http.StatusInternalServerError)
		return
	}

    w.WriteHeader(http.StatusCreated)
    fmt.Fprint(w, `{"status":"joined"}`)
}
