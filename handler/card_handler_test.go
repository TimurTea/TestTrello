package handler

import (
	"awesomeProject2/cmd/dto"
	"awesomeProject2/cmd/model"
	"bytes"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleCards_GET(t *testing.T) {
	id := 1
	tests := []struct {
		name           string
		requestBody    dto.CardDTO
		mockReturn     []model.Card
		mockError      error
		expectedStatus int
	}{
		{
			name:           "success with cards",
			requestBody:    dto.CardDTO{ID: &id},
			mockReturn:     []model.Card{{ID: 1, Title: "Card1"}, {ID: 2, Title: "Card2"}},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "empty cards list",
			requestBody:    dto.CardDTO{ID: &id},
			mockReturn:     []model.Card{},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "service error",
			requestBody:    dto.CardDTO{ID: &id},
			mockReturn:     nil,
			mockError:      errors.New("fail"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := new(MockCardStorage)
			logger := zap.NewNop()
			handler := NewCardHandler(mock, logger)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodGet, "/cards", bytes.NewReader(body))
			rec := httptest.NewRecorder()

			mock.On("GetCards", tt.requestBody.ID).Return(tt.mockReturn, tt.mockError)

			handler.HandleCards(rec, req)

			require.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectedStatus == http.StatusOK {
				var resp []dto.CardDTO
				err := json.NewDecoder(rec.Body).Decode(&resp)
				require.NoError(t, err)
				require.Equal(t, len(tt.mockReturn), len(resp))
				for i, card := range tt.mockReturn {
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
		mockReturn      model.Card
		mockError       error
		expectedStatus  int
		rawBody         string
		expectEncodeErr bool
	}{
		{
			name:           "success",
			requestBody:    dto.CreateCardDTO{ListID: 1, Title: "New Card", Description: "desc"},
			mockReturn:     model.Card{ID: 1, ListID: 1, Title: "New Card", Description: "desc"},
			mockError:      nil,
			expectedStatus: http.StatusOK,
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
			mockReturn:     model.Card{},
			mockError:      errors.New("fail"),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "error decode",
			rawBody:        `{"title":123}`,
			mockReturn:     model.Card{ID: 123, Title: "EncodeFail"},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:            "error from encode",
			requestBody:     dto.CreateCardDTO{Title: "EncodeFail"},
			mockReturn:      model.Card{ID: 123, Title: "EncodeFail"},
			expectedStatus:  http.StatusBadRequest,
			expectEncodeErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := new(MockCardStorage)
			logger := zap.NewNop()
			handler := NewCardHandler(mock, logger)

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

			if tt.requestBody.ListID != 0 && len(tt.requestBody.Title) != 0 {
				mock.On("CreateCard", model.CardInputCreate{
					ListID:      tt.requestBody.ListID,
					Title:       tt.requestBody.Title,
					Description: tt.requestBody.Description,
				}).Return(tt.mockReturn, tt.mockError)
			}

			handler.HandleCards(rw, req)

			require.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectedStatus == http.StatusOK && !tt.expectEncodeErr {
				var resp dto.CardDTO
				err := json.NewDecoder(rec.Body).Decode(&resp)
				require.NoError(t, err)
				require.Equal(t, tt.mockReturn.ID, *resp.ID)
				require.Equal(t, tt.mockReturn.Title, resp.Title)
			}
		})
	}
}

func TestHandleCards_DELETE(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    dto.DeleteCardDTO
		mockReturn     model.Card
		mockError      error
		expectedStatus int
	}{
		{
			name:           "success",
			requestBody:    dto.DeleteCardDTO{CardID: 1, ListID: 2},
			mockReturn:     model.Card{ID: 1, Title: "Deleted Card"},
			mockError:      nil,
			expectedStatus: http.StatusOK,
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
			mock := new(MockCardStorage)
			logger := zap.NewNop()
			handler := NewCardHandler(mock, logger)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodDelete, "/cards", bytes.NewReader(body))
			rec := httptest.NewRecorder()

			if tt.requestBody.CardID != 0 && tt.requestBody.ListID != 0 {
				mock.On("DeleteCard", tt.requestBody.ListID, tt.requestBody.CardID).
					Return(tt.mockReturn, tt.mockError)
			}

			handler.HandleCards(rec, req)

			require.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectedStatus == http.StatusOK {
				var resp dto.CardDTO
				err := json.NewDecoder(rec.Body).Decode(&resp)
				require.NoError(t, err)
				require.Equal(t, tt.mockReturn.ID, *resp.ID)
				require.Equal(t, tt.mockReturn.Title, resp.Title)
			}
		})
	}
}
func TestHandleCards_PUT(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    dto.UpdateCardDTO
		mockReturn     model.Card
		mockError      error
		expectedStatus int
	}{
		{
			name: "success",
			requestBody: dto.UpdateCardDTO{
				ID:          1,
				ListID:      2,
				Title:       "Updated Card",
				Description: "Updated Desc",
			},
			mockReturn: model.Card{
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
			mock := new(MockCardStorage)
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
				}).Return(tt.mockReturn, tt.mockError)
			}

			handler.HandleCards(rec, req)

			require.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectedStatus == http.StatusOK {
				var resp dto.CardDTO
				err := json.NewDecoder(rec.Body).Decode(&resp)
				require.NoError(t, err)
				require.Equal(t, tt.mockReturn.ID, *resp.ID)
				require.Equal(t, tt.mockReturn.Title, resp.Title)
			}
		})
	}
}
func TestHandleCard_MethodNotAllowed(t *testing.T) {
	mockService := new(MockCardStorage)
	logger := zap.NewNop()
	handler := NewCardHandler(mockService, logger)

	req := httptest.NewRequest(http.MethodPatch, "/cards", nil)
	rec := httptest.NewRecorder()

	handler.HandleCards(rec, req)

	require.Equal(t, http.StatusMethodNotAllowed, rec.Code)
	require.Contains(t, rec.Body.String(), "Method not allowed")
}
