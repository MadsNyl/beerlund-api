package db


func (p *PostgresStore) RemoveParticipant(eventID int, userID string) error {
	_, err := p.db.Exec(
		`DELETE FROM participant WHERE event_id = $1 AND user_id = $2`,
		eventID, userID,
	)
	if err != nil {
		return err
	}
	return nil
}