package service

import (
	"awesomeProject2/cmd/model"
	"errors"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
)

func TestCreateBoard(t *testing.T) {
	tests := []struct {
		title         string
		inputTitle    string
		mockResult    model.Board
		mockError     error
		expectedBoard model.Board
		expectError   bool
	}{
		{
			title:      "успешное создание доски",
			inputTitle: "Test Board",
			mockResult: model.Board{ID: 1, Title: "Test Board"},
			mockError:  nil,
			expectedBoard: model.Board{
				ID:    1,
				Title: "Test Board",
			},
			expectError: false,
		},
		{
			title:         "ошибка при создании доски",
			inputTitle:    "Bad Board",
			mockResult:    model.Board{},
			mockError:     errors.New("failed to create board"),
			expectedBoard: model.Board{},
			expectError:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			mockStorage := new(MockBoardStorage)
			logger := zap.NewNop()
			service := NewBoardService(mockStorage, logger)
			mockStorage.On("CreateBoard", tt.inputTitle).Return(tt.mockResult, tt.mockError)
			result, err := service.CreateBoard(tt.inputTitle)
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedBoard, result)
			}
			mockStorage.AssertExpectations(t)
		})
	}
}
func TestGetBoard(t *testing.T) {
	tests := []struct {
		title         string
		mockResult    []model.Board
		mockError     error
		expectedError bool
	}{
		{
			title: "success",
			mockResult: []model.Board{
				{ID: 1, Title: "Test Board"},
				{ID: 2, Title: "Test Board"},
			},
			mockError:     nil,
			expectedError: false,
		},
		{
			title:         "error",
			mockResult:    nil,
			mockError:     errors.New("failed to get board"),
			expectedError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			mockStorage := new(MockBoardStorage)
			logger := zap.NewNop()
			boardService := NewBoardService(mockStorage, logger)
			mockStorage.On("GetBoards").Return(tt.mockResult, tt.mockError)
			boards, err := boardService.GetBoards()
			if tt.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.mockResult, boards)
			}
			mockStorage.AssertExpectations(t)
		})
	}
}
