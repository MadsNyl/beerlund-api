package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/user"
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

	// Add Clerk users to the event participants and leaderboard
	uidSet := map[string]struct{}{}
    for _, pt := range event.Participants {
        uidSet[pt.UserID] = struct{}{}
    }
    for _, lb := range event.Leaderboard {
        uidSet[lb.ParticipantId] = struct{}{}
    }

	userIDs := make([]string, 0, len(uidSet))
	for uid := range uidSet {
		userIDs = append(userIDs, uid)
	}

	params := &user.ListParams{
		UserIDs: userIDs,
	}

	usrList, err := user.List(r.Context(), params)
	if err != nil {
		http.Error(w, "Failed to fetch user details", http.StatusInternalServerError)
		log.Printf("Error fetching user details: %v", err)
		return
	}

	users := make(map[string]*clerk.User, len(usrList.Users))
	for _, u := range usrList.Users {
		users[u.ID] = u
	}

	for i := range event.Participants {
		event.Participants[i].User = users[event.Participants[i].UserID]
	}
	for i := range event.Leaderboard {
		event.Leaderboard[i].User = users[event.Leaderboard[i].ParticipantId]
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(event); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}