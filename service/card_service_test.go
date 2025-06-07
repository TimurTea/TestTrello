package service

import (
	"awesomeProject2/cmd/model"
	"errors"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
)

func TestCreateCard(t *testing.T) {
	tests := []struct {
		title         string
		listID        int
		mockReturn    model.Card
		mockError     error
		expectedError bool
	}{
		{
			title:  "success",
			listID: 1,
			mockReturn: model.Card{
				ID:     1,
				ListID: 1,
				Title:  "Test Card",
			},
			mockError:     nil,
			expectedError: false,
		},
		{
			title:         "error, storage returns error",
			listID:        666,
			mockReturn:    model.Card{},
			mockError:     errors.New("failed to create card"),
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			mockStorage := new(MockCardStorage)
			svc := CardService{Storage: mockStorage}

			input := model.CardInputCreate{
				Title:  tt.title,
				ListID: tt.listID,
			}

			mockStorage.On("CreateCard", input).Return(tt.mockReturn, tt.mockError)

			result, err := svc.CreateCard(input)

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
func TestGetCard(t *testing.T) {
	tests := []struct {
		title         string
		listID        int
		mockResult    []model.Card
		mockError     error
		expectedError bool
	}{
		{
			title:  "success",
			listID: 1,
			mockResult: []model.Card{
				{ID: 1, ListID: 1, Title: "Test List"},
				{ID: 2, ListID: 1, Title: "Test List"},
			},
			mockError:     nil,
			expectedError: false,
		},
		{
			title:         "error, storage returns error",
			listID:        666,
			mockResult:    nil,
			mockError:     errors.New("failed to get list"),
			expectedError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			mockStorage := new(MockCardStorage)
			logger := zap.NewNop()
			cardService := NewCardService(mockStorage, logger)
			listID := tt.listID
			mockStorage.On("GetCards", &listID).Return(tt.mockResult, tt.mockError)
			lists, err := cardService.GetCards(&listID)
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
func TestDeleteCard(t *testing.T) {
	tests := []struct {
		title         string
		cardID        int
		listID        int
		mockError     error
		expectedError bool
	}{
		{
			title:         "success",
			cardID:        1,
			listID:        1,
			mockError:     nil,
			expectedError: false,
		},
		{
			title:         "error, storage returns error",
			cardID:        666,
			listID:        666,
			mockError:     errors.New("failed to delete card"),
			expectedError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			mockStorage := new(MockCardStorage)
			service := CardService{Storage: mockStorage}
			mockStorage.On("DeleteCard", tt.listID, tt.cardID).Return(model.Card{}, tt.mockError)
			_, err := service.DeleteCard(tt.cardID, tt.listID)
			if tt.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			mockStorage.AssertExpectations(t)
		})
	}
}
func TestUpdateCard(t *testing.T) {
	tests := []struct {
		title         string
		updated       model.Card
		mockReturn    model.Card
		mockError     error
		expectedError bool
	}{
		{
			title: "success",
			updated: model.Card{
				ID:    1,
				Title: "Update Card",
			},
			mockReturn: model.Card{
				ID:    1,
				Title: "Update Card",
			},
			mockError:     nil,
			expectedError: false,
		},
		{
			title: "error, storage returns error",
			updated: model.Card{
				ID:    666,
				Title: "Doesn't matter",
			},
			mockReturn:    model.Card{},
			mockError:     errors.New("failed to update card"),
			expectedError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			mockStorage := new(MockCardStorage)
			service := CardService{Storage: mockStorage}
			mockStorage.On("UpdateCard", tt.updated).Return(tt.mockReturn, tt.mockError)
			result, err := service.UpdateCard(tt.updated)
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
