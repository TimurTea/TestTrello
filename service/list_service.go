package service

import (
	"awesomeProject2/cmd/model"
	"go.uber.org/zap"
)

type ListService struct {
	Storage ListStorage
	logger  *zap.Logger
}

func NewListService(storage ListStorage, logger *zap.Logger) *ListService {
	return &ListService{
		Storage: storage,
		logger:  logger,
	}
}
func (s ListService) GetLists(ListID *int) ([]model.List, error) {
	return s.Storage.GetLists(ListID)
}
func (s ListService) CreateList(input model.ListInputCreate) (model.List, error) {
	return s.Storage.CreateList(input)
}
