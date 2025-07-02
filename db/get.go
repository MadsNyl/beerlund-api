package db

import (
	"database/sql"
	"fmt"

	"api.beerlund.com/m/models"
)

func (p *PostgresStore) GetEvent(id int) (models.EventResponse, error) {
    // 1) Fetch the event row
    const evQ = `
    SELECT
      id, name, description, address, zip_code, city, country,
      start_time, end_time,
      COALESCE(image_url, '') AS image_url,
      max_participants
    FROM event
    WHERE id = $1;
    `

    var ev models.EventResponse
    if err := p.db.QueryRow(evQ, id).Scan(
        &ev.ID,
        &ev.Name,
        &ev.Description,
        &ev.Address,
        &ev.ZipCode,
        &ev.City,
        &ev.Country,
        &ev.StartTime,
        &ev.EndTime,
        &ev.ImageURL,
        &ev.MaxParticipants,
    ); err != nil {
        if err == sql.ErrNoRows {
            return models.EventResponse{}, fmt.Errorf("event %d not found", id)
        }
        return models.EventResponse{}, err
    }

    // 2) Fetch all participants (empty slice if none)
    const partQ = `
    SELECT
    p.id,
    p.event_id,
    p.user_id,
    to_char(p.created_at, 'YYYY-MM-DD"T"HH24:MI:SSZ') AS created_at,
    to_char(p.updated_at, 'YYYY-MM-DD"T"HH24:MI:SSZ') AS updated_at
    FROM participant p
    WHERE p.event_id = $1;
    `

    rows, err := p.db.Query(partQ, id)
    if err != nil {
        return models.EventResponse{}, err
    }
    defer rows.Close()

    ev.Participants = make([]models.Participant, 0)
    for rows.Next() {
        var pt models.Participant
        if err := rows.Scan(
            &pt.ID,
            &pt.EventID,
            &pt.UserID,
            &pt.CreatedAt,
            &pt.UpdatedAt,
        ); err != nil {
            return models.EventResponse{}, err
        }
        ev.Participants = append(ev.Participants, pt)
    }
    if err := rows.Err(); err != nil {
        return models.EventResponse{}, err
    }
    return ev, nil
}

func (p *PostgresStore) GetMaxParticipants(eventID int) (int, error) {
	var maxParticipants int
	err := p.db.QueryRow(
		`SELECT max_participants FROM event WHERE id = $1`,
		eventID,
	).Scan(&maxParticipants)
	if err != nil {
		return 0, err
	}
	return maxParticipants, nil
}