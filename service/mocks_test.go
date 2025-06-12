package service

import (
	"awesomeProject2/cmd/model"
	"github.com/stretchr/testify/mock"
)

type MockBoardService struct {
	mock.Mock
}
type MockListService struct {
	mock.Mock
}
type MockCardService struct {
	mock.Mock
}

func (m *MockBoardService) CreateBoard(title string) (model.Board, error) {
	args := m.Called(title)
	return args.Get(0).(model.Board), args.Error(1)
}
func (m *MockBoardService) GetBoards() ([]model.Board, error) {
	args := m.Called()
	return args.Get(0).([]model.Board), args.Error(1)
}
func (m *MockListService) CreateList(input model.ListInputCreate) (model.List, error) {
	args := m.Called(input)
	return args.Get(0).(model.List), args.Error(1)
}
func (m *MockListService) GetLists(BoardID *int) ([]model.List, error) {
	args := m.Called(BoardID)
	return args.Get(0).([]model.List), args.Error(1)
}
func (m *MockCardService) CreateCard(input model.CardInputCreate) (model.Card, error) {
	args := m.Called(input)
	return args.Get(0).(model.Card), args.Error(1)
}
func (m *MockCardService) GetCards(ListID *int) ([]model.Card, error) {
	args := m.Called(ListID)
	return args.Get(0).([]model.Card), args.Error(1)
}
func (m *MockCardService) DeleteCard(listID, cardID int) (model.Card, error) {
	args := m.Called(listID, cardID)
	return args.Get(0).(model.Card), args.Error(1)
}
func (m *MockCardService) UpdateCard(updated model.Card) (model.Card, error) {
	args := m.Called(updated)
	return args.Get(0).(model.Card), args.Error(1)
}
