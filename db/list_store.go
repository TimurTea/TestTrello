package storage

import (
	"awesomeProject2/cmd/model"
	"github.com/jmoiron/sqlx"
)

type ListStorage struct {
	DB *sqlx.DB
}

func NewListStorage(db *sqlx.DB) *ListStorage { return &ListStorage{db} }
func (s *ListStorage) GetLists(boardID *int) ([]model.List, error) {
	var lists []model.List
	var err error
	if boardID != nil {
		err = s.DB.Select(&lists, "SELECT * FROM lists WHERE board_id = $1", *boardID)
	} else {
		err = s.DB.Select(&lists, "SELECT * FROM lists")
	}
	return lists, err
}
func (s *ListStorage) CreateList(input model.ListInputCreate) (model.List, error) {
	var list model.List
	query := `INSERT INTO lists (title, board_id) VALUES ($1, $2) RETURNING id, title, board_id`
	err := s.DB.Get(&list, query, input.Title, input.ListID)
	return list, err
}
