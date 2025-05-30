package service

import (
	"awesomeProject2/cmd/model"
	"go.uber.org/zap"
)

type ListService struct {
	storage ListStorage
	logger  *zap.Logger
}

func NewListService(storage ListStorage, logger *zap.Logger) *ListService {
	return &ListService{
		storage: storage,
		logger:  logger,
	}
}
func (s ListService) GetLists(ListID *int) ([]model.List, error) {
	return s.storage.GetLists(ListID)
}
func (s ListService) CreateList(input model.ListInputCreate) (model.List, error) {
	return s.storage.CreateList(input)
}
