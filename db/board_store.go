package storage

import (
	"awesomeProject2/cmd/model"
	"github.com/jmoiron/sqlx"
)

type BoardStorage struct {
	DB *sqlx.DB
}

func NewBoardStorage(db *sqlx.DB) *BoardStorage { return &BoardStorage{db} }
func (s *BoardStorage) GetBoards() ([]model.Board, error) {
	var boards []model.Board
	err := s.DB.Select(&boards, "SELECT * FROM boards")
	return boards, err
}
func (s *BoardStorage) CreateBoard(title string) (model.Board, error) {
	var board model.Board
	query := `INSERT INTO boards (title) VALUES ($1) RETURNING id, title`
	err := s.DB.Get(&board, query, title)
	return board, err
}
