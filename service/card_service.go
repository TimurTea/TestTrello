package service

import (
	"awesomeProject2/cmd/model"
	"go.uber.org/zap"
)

type CardService struct {
	Storage CardStorage
	logger  *zap.Logger
}

func NewCardService(storage CardStorage, logger *zap.Logger) *CardService {
	return &CardService{
		Storage: storage,
		logger:  logger,
	}
}
func (s CardService) GetCards(CardID *int) ([]model.Card, error) {
	return s.Storage.GetCards(CardID)
}
func (s CardService) CreateCard(input model.CardInputCreate) (model.Card, error) {
	return s.Storage.CreateCard(input)
}
func (s CardService) DeleteCard(listID int, cardID int) (model.Card, error) {
	return s.Storage.DeleteCard(listID, cardID)
}
func (s CardService) UpdateCard(updated model.Card) (model.Card, error) {
	return s.Storage.UpdateCard(updated)
}
