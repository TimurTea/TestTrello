package service

import (
	"awesomeProject2/cmd/model"
	"go.uber.org/zap"
)

type BoardService struct {
	Storage BoardStorage
	logger  *zap.Logger
}

func NewBoardService(storage BoardStorage, logger *zap.Logger) *BoardService {
	return &BoardService{
		Storage: storage,
		logger:  logger,
	}
}
func (s BoardService) GetBoards() ([]model.Board, error) {
	return s.Storage.GetBoards()
}
func (s BoardService) CreateBoard(title string) (model.Board, error) {
	return s.Storage.CreateBoard(title)
}
