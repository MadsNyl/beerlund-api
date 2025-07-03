package db

import (
	"database/sql"

	"api.beerlund.com/m/models"
	_ "github.com/lib/pq"
)

type Store interface {
	Close() error
	Init(connUrl string) (error)
	ListEvents(page, limit int, ended bool) (models.EventListResponse, error)
	GetEvent(id int) (models.EventResponse, error)
	GetMaxParticipants(eventID int) (int, error)
	CountParticipants(eventID int) (int, error)
	InsertParticipant(eventID int, userID string) error
	IsParticipating(eventID int, userID string) (bool, error)
	RemoveParticipant(eventID int, userID string) error
}

type PostgresStore struct {
	db    *sql.DB
}

func NewPostgresStore() *PostgresStore {
	return &PostgresStore{}
}

func (p *PostgresStore) Init(connUrl string) error {
	db, err := sql.Open("postgres", connUrl)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	p.db = db
	return nil
}

func (p *PostgresStore) Close() error {
	return p.db.Close()
}
