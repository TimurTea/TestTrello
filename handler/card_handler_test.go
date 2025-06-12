package handler

import (
	"awesomeProject2/cmd/dto"
	"awesomeProject2/cmd/helper"
	"awesomeProject2/cmd/model"
	"bytes"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleCards_GET(t *testing.T) {
	tests := []struct {
		name            string
		requestBody     dto.CardDTO
		serviceResponse []model.Card
		mockError       error
		expectedStatus  int
	}{
		{
			name:            "success with cards",
			requestBody:     dto.CardDTO{ID: helper.GetPointer(1)},
			serviceResponse: []model.Card{{ID: 1, Title: "Card1"}, {ID: 2, Title: "Card2"}},
			mockError:       nil,
			expectedStatus:  http.StatusOK,
		},
		{
			name:            "empty cards list",
			requestBody:     dto.CardDTO{ID: helper.GetPointer(1)},
			serviceResponse: []model.Card{},
			mockError:       nil,
			expectedStatus:  http.StatusOK,
		},
		{
			name:            "service error",
			requestBody:     dto.CardDTO{ID: helper.GetPointer(1)},
			serviceResponse: nil,
			mockError:       errors.New("fail"),
			expectedStatus:  http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := new(MockCardService)
			logger := zap.NewNop()
			handler := NewCardHandler(mock, logger)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodGet, "/cards", bytes.NewReader(body))
			rec := httptest.NewRecorder()

			mock.On("GetCards", tt.requestBody.ID).Return(tt.serviceResponse, tt.mockError)

			handler.HandleCards(rec, req)

			require.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectedStatus == http.StatusOK {
				var resp []dto.CardDTO
				err := json.NewDecoder(rec.Body).Decode(&resp)
				require.NoError(t, err)
				require.Equal(t, len(tt.serviceResponse), len(resp))
				for i, card := range tt.serviceResponse {
					require.Equal(t, card.ID, *resp[i].ID)
					require.Equal(t, card.Title, resp[i].Title)
				}
			}
		})
	}
}
func TestHandleCards_POST(t *testing.T) {
	tests := []struct {
		name            string
		requestBody     dto.CreateCardDTO
		serviceResponse model.Card
		mockError       error
		expectedStatus  int
		rawBody         string
		expectEncodeErr bool
		setupMock       func(s *MockCardService)
	}{
		{
			name:            "success",
			requestBody:     dto.CreateCardDTO{ListID: 1, Title: "New Card", Description: "desc"},
			serviceResponse: model.Card{ID: 1, ListID: 1, Title: "New Card", Description: "desc"},
			mockError:       nil,
			expectedStatus:  http.StatusOK,
			setupMock: func(s *MockCardService) {
				s.On("CreateCard", model.CardInputCreate{
					ListID:      1,
					Title:       "New Card",
					Description: "desc",
				}).Return(model.Card{ID: 1, ListID: 1, Title: "New Card", Description: "desc"}, nil)
			},
		},
		{
			name:           "missing list id",
			requestBody:    dto.CreateCardDTO{Title: "New Card"},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "missing title",
			requestBody:    dto.CreateCardDTO{ListID: 1},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "service error",
			requestBody:    dto.CreateCardDTO{ListID: 1, Title: "Fail Card"},
			mockError:      errors.New("fail"),
			expectedStatus: http.StatusBadRequest,
			setupMock: func(s *MockCardService) {
				s.On("CreateCard", model.CardInputCreate{
					ListID: 1,
					Title:  "Fail Card",
				}).Return(model.Card{}, errors.New("fail"))
			},
		},
		{
			name:            "error decode",
			rawBody:         `{"title":123}`,
			serviceResponse: model.Card{ID: 123, Title: "EncodeFail"},
			expectedStatus:  http.StatusBadRequest,
		},
		{
			name:            "error from encode",
			requestBody:     dto.CreateCardDTO{Title: "EncodeFail"},
			serviceResponse: model.Card{ID: 123, Title: "EncodeFail"},
			expectedStatus:  http.StatusBadRequest,
			expectEncodeErr: true,
			setupMock: func(s *MockCardService) {
				s.On("CreateCard", mock.Anything).Return(model.Card{ID: 123, Title: "EncodeFail"}, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockCardService)
			logger := zap.NewNop()
			handler := NewCardHandler(mockService, logger)

			var req *http.Request
			if tt.rawBody != "" {
				req = httptest.NewRequest(http.MethodPost, "/cards", bytes.NewBufferString(tt.rawBody))
			} else {
				body, _ := json.Marshal(tt.requestBody)
				req = httptest.NewRequest(http.MethodPost, "/cards", bytes.NewReader(body))
			}

			rec := httptest.NewRecorder()
			var rw http.ResponseWriter = rec

			if tt.expectEncodeErr {
				rw = &badResponseWriter{ResponseWriter: rec}
			}

			if tt.setupMock != nil {
				tt.setupMock(mockService)
			}

			handler.HandleCards(rw, req)

			require.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectedStatus == http.StatusOK && !tt.expectEncodeErr {
				var resp dto.CardDTO
				err := json.NewDecoder(rec.Body).Decode(&resp)
				require.NoError(t, err)

				expected := dto.CardToDTO(tt.serviceResponse)
				require.Equal(t, expected, resp)
			}
		})
	}
}
func TestHandleCards_DELETE(t *testing.T) {
	tests := []struct {
		name            string
		requestBody     dto.DeleteCardDTO
		serviceResponse model.Card
		mockError       error
		expectedStatus  int
	}{
		{
			name:            "success",
			requestBody:     dto.DeleteCardDTO{CardID: 1, ListID: 2},
			serviceResponse: model.Card{ID: 1, Title: "Deleted Card"},
			mockError:       nil,
			expectedStatus:  http.StatusOK,
		},
		{
			name:           "missing list id",
			requestBody:    dto.DeleteCardDTO{CardID: 1},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "missing card id",
			requestBody:    dto.DeleteCardDTO{ListID: 2},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "service error",
			requestBody:    dto.DeleteCardDTO{CardID: 1, ListID: 2},
			mockError:      errors.New("fail"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := new(MockCardService)
			logger := zap.NewNop()
			handler := NewCardHandler(mock, logger)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodDelete, "/cards", bytes.NewReader(body))
			rec := httptest.NewRecorder()

			if tt.requestBody.CardID != 0 && tt.requestBody.ListID != 0 {
				mock.On("DeleteCard", tt.requestBody.ListID, tt.requestBody.CardID).
					Return(tt.serviceResponse, tt.mockError)
			}

			handler.HandleCards(rec, req)

			require.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectedStatus == http.StatusOK {
				var resp dto.CardDTO
				err := json.NewDecoder(rec.Body).Decode(&resp)
				require.NoError(t, err)
				require.Equal(t, tt.serviceResponse.ID, *resp.ID)
				require.Equal(t, tt.serviceResponse.Title, resp.Title)
			}
		})
	}
}
func TestHandleCards_PUT(t *testing.T) {
	tests := []struct {
		name            string
		requestBody     dto.UpdateCardDTO
		serviceResponse model.Card
		mockError       error
		expectedStatus  int
	}{
		{
			name: "success",
			requestBody: dto.UpdateCardDTO{
				ID:          1,
				ListID:      2,
				Title:       "Updated Card",
				Description: "Updated Desc",
			},
			serviceResponse: model.Card{
				ID:          1,
				ListID:      2,
				Title:       "Updated Card",
				Description: "Updated Desc",
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name: "missing list id",
			requestBody: dto.UpdateCardDTO{
				ID:    1,
				Title: "Title",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "missing card id",
			requestBody: dto.UpdateCardDTO{
				ListID: 2,
				Title:  "Title",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "service error",
			requestBody: dto.UpdateCardDTO{
				ID:     1,
				ListID: 2,
				Title:  "Title",
			},
			mockError:      errors.New("fail"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := new(MockCardService)
			logger := zap.NewNop()
			handler := NewCardHandler(mock, logger)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPut, "/cards", bytes.NewReader(body))
			rec := httptest.NewRecorder()

			if tt.requestBody.ID != 0 && tt.requestBody.ListID != 0 {
				mock.On("UpdateCard", model.Card{
					ID:          tt.requestBody.ID,
					ListID:      tt.requestBody.ListID,
					Title:       tt.requestBody.Title,
					Description: tt.requestBody.Description,
				}).Return(tt.serviceResponse, tt.mockError)
			}

			handler.HandleCards(rec, req)

			require.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectedStatus == http.StatusOK {
				var resp dto.CardDTO
				err := json.NewDecoder(rec.Body).Decode(&resp)
				require.NoError(t, err)
				require.Equal(t, tt.serviceResponse.ID, *resp.ID)
				require.Equal(t, tt.serviceResponse.Title, resp.Title)
			}
		})
	}
}
func TestHandleCard_MethodNotAllowed(t *testing.T) {
	mockService := new(MockCardService)
	logger := zap.NewNop()
	handler := NewCardHandler(mockService, logger)

	req := httptest.NewRequest(http.MethodPatch, "/cards", nil)
	rec := httptest.NewRecorder()

	handler.HandleCards(rec, req)

	require.Equal(t, http.StatusMethodNotAllowed, rec.Code)
	require.Contains(t, rec.Body.String(), "Method not allowed")
}
