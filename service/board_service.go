package service

import (
	"awesomeProject2/cmd/model"
	"go.uber.org/zap"
)

type BoardService struct {
	storage BoardStorage
	logger  *zap.Logger
}

func NewBoardService(storage BoardStorage, logger *zap.Logger) *BoardService {
	return &BoardService{
		storage: storage,
		logger:  logger,
	}
}
func (s BoardService) GetBoards() ([]model.Board, error) {
	return s.storage.GetBoards()
}
func (s BoardService) CreateBoard(title string) (model.Board, error) {
	return s.storage.CreateBoard(title)
}
