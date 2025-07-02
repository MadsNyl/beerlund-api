package db


func (p *PostgresStore) InsertParticipant(eventID int, userID string) error {
	query := `
		INSERT INTO participant (event_id, user_id)
		VALUES ($1, $2)
		ON CONFLICT (event_id, user_id) DO NOTHING;
	`

	_, err := p.db.Exec(query, eventID, userID)
	if err != nil {
		return err
	}

	return nil
}