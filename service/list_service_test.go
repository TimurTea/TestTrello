package service

import (
	"awesomeProject2/cmd/model"
	"errors"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
)

func TestCreateList(t *testing.T) {
	tests := []struct {
		title         string
		boardID       int
		mockReturn    model.List
		mockError     error
		expectedError bool
	}{
		{
			title:   "success",
			boardID: 1,
			mockReturn: model.List{
				ID:      1,
				BoardID: 1,
				Title:   "Test List",
			},
			mockError:     nil,
			expectedError: false,
		},
		{
			title:         "error, storage returns error",
			boardID:       666,
			mockReturn:    model.List{},
			mockError:     errors.New("failed to create list"),
			expectedError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			mockStorage := new(MockListService)
			svc := ListService{Storage: mockStorage}

			input := model.ListInputCreate{
				Title:   tt.title,
				BoardID: tt.boardID,
			}

			mockStorage.On("CreateList", input).Return(tt.mockReturn, tt.mockError)

			result, err := svc.CreateList(input)

			if tt.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.mockReturn, result)
			}

			mockStorage.AssertExpectations(t)
		})
	}
}
func TestGetList(t *testing.T) {
	tests := []struct {
		title         string
		boardID       int
		mockResult    []model.List
		mockError     error
		expectedError bool
	}{
		{
			title:   "success",
			boardID: 1,
			mockResult: []model.List{
				{ID: 1, BoardID: 1, Title: "Test List"},
				{ID: 2, BoardID: 1, Title: "Test List"},
			},
			mockError:     nil,
			expectedError: false,
		},
		{
			title:         "error, storage returns error",
			boardID:       666,
			mockResult:    nil,
			mockError:     errors.New("failed to get list"),
			expectedError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			mockStorage := new(MockListService)
			logger := zap.NewNop()
			listService := NewListService(mockStorage, logger)
			boardID := tt.boardID
			mockStorage.On("GetLists", &boardID).Return(tt.mockResult, tt.mockError)
			lists, err := listService.GetLists(&boardID)
			if tt.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.mockResult, lists)
			}
			mockStorage.AssertExpectations(t)
		})
	}
}
