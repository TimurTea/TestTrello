package storage

import (
	"awesomeProject2/cmd/model"
	"github.com/jmoiron/sqlx"
)

type CardStorage struct {
	DB *sqlx.DB
}

func NewCardStorage(db *sqlx.DB) *CardStorage { return &CardStorage{db} }
func (s *CardStorage) GetCards(listID *int) ([]model.Card, error) {
	var cards []model.Card
	var err error
	if listID != nil {
		err = s.DB.Select(&cards, "SELECT * FROM cards WHERE list_id = $1", *listID)
	} else {
		err = s.DB.Select(&cards, "SELECT * FROM cards")
	}
	return cards, err
}
func (s *CardStorage) CreateCard(input model.CardInputCreate) (model.Card, error) {
	query := `INSERT INTO cards(title,board_id,list_id, description) VALUES ($1, $2, $3) RETURNING  title, board_id, description, list_id`
	var card model.Card
	err := s.DB.Get(&card, query, input.Title, input.Description, input.ListID)
	return card, err
}
func (s *CardStorage) DeleteCard(listID int, cardID int) (model.Card, error) {
	query := `DELETE FROM cards WHERE id = $1 AND list_id = $2 RETURNING id, list_id`
	var card model.Card
	err := s.DB.Get(&card, query, cardID, listID)
	return card, err
}

func (s *CardStorage) UpdateCard(updated model.Card) (model.Card, error) {
	query := `UPDATE cards SET title = $1, description = $2, list_id = $3 WHERE id = $4 RETURNING id, title, description, list_id`
	var card model.Card
	err := s.DB.Get(&card, query, updated.Title, updated.Description, updated.ListID, updated.ID)
	return card, err
}
