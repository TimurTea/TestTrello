package handler

import "awesomeProject2/cmd/model"

type BoardService interface {
	GetBoards() ([]model.Board, error)
	CreateBoard(title string) (model.Board, error)
}
type ListService interface {
	GetLists(boardID *int) ([]model.List, error)
	CreateList(input model.ListInputCreate) (model.List, error)
}
type CardService interface {
	GetCards(boardID *int) ([]model.Card, error)
	CreateCard(input model.CardInputCreate) (model.Card, error)
	DeleteCard(listID int, cardID int) (model.Card, error)
	UpdateCard(updated model.Card) (model.Card, error)
}
type SQLExecutor interface {
	Get(dest interface{}, query string, args ...interface{}) error
	Select(dest interface{}, query string, args ...interface{}) error
}
