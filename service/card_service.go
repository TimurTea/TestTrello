package service

import "awesomeProject2/cmd/model"

type CardService struct {
	storage CardStorage
}

func NewCardService(storage CardStorage) *CardService {
	return &CardService{storage: storage}
}
func (s CardService) GetCards(CardID *int) ([]model.Card, error) {
	return s.storage.GetCards(CardID)
}
func (s CardService) CreateCard(input model.CardInputCreate) (model.Card, error) {
	return s.storage.CreateCard(input)
}
func (s CardService) DeleteCard(listID int, cardID int) (model.Card, error) {
	return s.storage.DeleteCard(listID, cardID)
}
func (s CardService) UpdateCard(updated model.Card) (model.Card, error) {
	return s.storage.UpdateCard(updated)
}
