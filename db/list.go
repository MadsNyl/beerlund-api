package db

import "api.beerlund.com/m/models"

func (p *PostgresStore) ListEvents(page, limit int) (models.EventListResponse, error) {
	query := `
		SELECT e.id, e.name, e.description, e.address, e.zip_code, e.city, e.country, e.start_time, e.end_time, e.image_url,
			   COUNT(p.user_id) AS participant_count
		FROM event e
		LEFT JOIN participant p ON e.id = p.event_id
		GROUP BY e.id, e.name, e.description, e.address, e.zip_code, e.city, e.country, e.start_time, e.end_time, e.image_url
		ORDER BY e.start_time DESC
		LIMIT $1 OFFSET $2;
	`
	offset := (page - 1) * limit

	rows, err := p.db.Query(query, limit, offset)
	if err != nil {
		return models.EventListResponse{}, err
	}
	defer rows.Close()

	var events []models.EventList
	for rows.Next() {
		var event models.EventList
		if err := rows.Scan(&event.ID, &event.Name, &event.Description, &event.Address,
			&event.ZipCode, &event.City, &event.Country, &event.StartTime,
			&event.EndTime, &event.ImageURL); err != nil {
			return models.EventListResponse{}, err
		}
		events = append(events, event)
	}

	totalCountQuery := `SELECT COUNT(*) FROM event;`
	var totalCount int
	if err := p.db.QueryRow(totalCountQuery).Scan(&totalCount); err != nil {
		return models.EventListResponse{}, err
	}

	nextPage := page + 1
	if (page * limit) >= totalCount {
		nextPage = 0 // No more pages
	}

	prevPage := page - 1
	if prevPage < 1 {
		prevPage = 0 // No previous page
	}

	return models.EventListResponse{
		Events:      events,
		NextPage:    nextPage,
		PrevPage:    prevPage,
		TotalCount:  totalCount,
	}, nil
}