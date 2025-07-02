package db

import (
    "api.beerlund.com/m/models"
)

func (p *PostgresStore) ListEvents(page, limit int, ended bool) (models.EventListResponse, error) {
    // The 3rd parameter ($3) is our 'ended' flag
    query := `
    SELECT 
      e.id,
      e.name,
      e.description,
      e.address,
      e.zip_code,
      e.city,
      e.country,
      e.start_time,
      e.end_time,
      COALESCE(e.image_url, '') AS image_url,
      e.max_participants,
      COUNT(p.user_id) AS participant_count
    FROM event e
    LEFT JOIN participant p ON e.id = p.event_id
    WHERE
      (
        $3::boolean = TRUE  AND e.end_time <  NOW()  -- past events
      ) OR (
        $3::boolean = FALSE AND e.end_time >= NOW()  -- future/ongoing
      )
    GROUP BY
      e.id, e.name, e.description, e.address, e.zip_code,
      e.city, e.country, e.start_time, e.end_time,
      e.image_url, e.max_participants
    ORDER BY e.start_time DESC
    LIMIT $1 OFFSET $2;
    `
    offset := (page - 1) * limit

	events := make([]models.EventList, 0, limit)

    rows, err := p.db.Query(query, limit, offset, ended)
    if err != nil {
        return models.EventListResponse{}, err
    }
    defer rows.Close()

    for rows.Next() {
        var ev models.EventList
        if err := rows.Scan(
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
            &ev.Participants,
        ); err != nil {
            return models.EventListResponse{}, err
        }
        events = append(events, ev)
    }
    if err := rows.Err(); err != nil {
        return models.EventListResponse{}, err
    }

    // total count for pagination (not filtering by ended)
    var totalCount int
    if err := p.db.
        QueryRow(`SELECT COUNT(*) FROM event`).
        Scan(&totalCount); err != nil {
        return models.EventListResponse{}, err
    }

    // compute next/prev pages
    nextPage := page + 1
    if page*limit >= totalCount {
        nextPage = 0
    }
    prevPage := page - 1
    if prevPage < 1 {
        prevPage = 0
    }

    return models.EventListResponse{
        Events:     events,
        NextPage:   nextPage,
        PrevPage:   prevPage,
        TotalCount: totalCount,
    }, nil
}
