package service

import "awesomeProject2/cmd/model"

//go:generate mockgen -source=interfaces.go -destination=../../mocks/mock_service_deps.go -package=mocks
type BoardStorage interface {
	GetBoards() ([]model.Board, error)
	CreateBoard(title string) (model.Board, error)
}

type ListStorage interface {
	GetLists(boardID *int) ([]model.List, error)
	CreateList(input model.ListInputCreate) (model.List, error)
}

type CardStorage interface {
	GetCards(boardID *int) ([]model.Card, error)
	CreateCard(input model.CardInputCreate) (model.Card, error)
	DeleteCard(listID int, cardID int) (model.Card, error)
	UpdateCard(updated model.Card) (model.Card, error)
}
