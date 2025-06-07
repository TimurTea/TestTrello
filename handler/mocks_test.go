package handler

import (
	"awesomeProject2/cmd/model"
	"github.com/stretchr/testify/mock"
)

type MockBoardStorage struct {
	mock.Mock
}
type MockListStorage struct {
	mock.Mock
}
type MockCardStorage struct {
	mock.Mock
}

func (m *MockBoardStorage) CreateBoard(title string) (model.Board, error) {
	args := m.Called(title)
	return args.Get(0).(model.Board), args.Error(1)
}
func (m *MockBoardStorage) GetBoards() ([]model.Board, error) {
	args := m.Called()
	return args.Get(0).([]model.Board), args.Error(1)
}
func (m *MockListStorage) CreateList(input model.ListInputCreate) (model.List, error) {
	args := m.Called(input)
	return args.Get(0).(model.List), args.Error(1)
}
func (m *MockListStorage) GetLists(BoardID *int) ([]model.List, error) {
	args := m.Called(BoardID)
	return args.Get(0).([]model.List), args.Error(1)
}
func (m *MockCardStorage) CreateCard(input model.CardInputCreate) (model.Card, error) {
	args := m.Called(input)
	return args.Get(0).(model.Card), args.Error(1)
}
func (m *MockCardStorage) GetCards(ListID *int) ([]model.Card, error) {
	args := m.Called(ListID)
	return args.Get(0).([]model.Card), args.Error(1)
}
func (m *MockCardStorage) DeleteCard(listID, cardID int) (model.Card, error) {
	args := m.Called(listID, cardID)
	return args.Get(0).(model.Card), args.Error(1)
}
func (m *MockCardStorage) UpdateCard(updated model.Card) (model.Card, error) {
	args := m.Called(updated)
	return args.Get(0).(model.Card), args.Error(1)
}
