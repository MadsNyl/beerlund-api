package db


func (p *PostgresStore) CountParticipants(eventID int) (int, error) {
	var count int
	err := p.db.QueryRow(
		`SELECT COUNT(*) FROM participant WHERE event_id = $1`,
		eventID,
	).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}