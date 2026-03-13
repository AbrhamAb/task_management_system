package persistence

import (
	"database/sql"
	"task-management-backend/internal/model/db"
)

type Store struct {
	db *sql.DB
}

func NewStore(dbConn *sql.DB) *Store {
	return &Store{db: dbConn}
}

func (s *Store) DB() db.DBTX {
	return s.db
}
