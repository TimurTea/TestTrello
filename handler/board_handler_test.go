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

func TestHandleBoards_GET(t *testing.T) {
	tests := []struct {
		name           string
		mockReturn     []model.Board
		mockError      error
		expectedStatus int
		expectedCount  int
	}{
		{
			name: "success",
			mockReturn: []model.Board{
				{ID: 1, Title: "Board 1"},
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedCount:  1,
		},
		{
			name:           "error from service",
			mockReturn:     nil,
			mockError:      errors.New("storage error"),
			expectedStatus: http.StatusInternalServerError,
			expectedCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockBoardStorage)
			logger := zap.NewNop()
			handler := NewBoardHandler(mockService, logger)

			mockService.On("GetBoards").Return(tt.mockReturn, tt.mockError)

			req := httptest.NewRequest(http.MethodGet, "/boards", nil)
			rec := httptest.NewRecorder()

			handler.HandleBoards(rec, req)

			require.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectedStatus == http.StatusOK {
				var response []dto.BoardDTO
				err := json.NewDecoder(rec.Body).Decode(&response)
				require.NoError(t, err)
				require.Len(t, response, tt.expectedCount)
			}

			mockService.AssertExpectations(t)
		})
	}
}

type badResponseWriter struct {
	http.ResponseWriter
}

func (b *badResponseWriter) Write(p []byte) (int, error) {
	return 0, errors.New("encode fail")
}
func TestHandleBoards_POST(t *testing.T) {
	tests := []struct {
		name            string
		requestBody     dto.CreateBoardDTO
		rawBody         string
		mockReturn      model.Board
		mockError       error
		expectedStatus  int
		expectEncodeErr bool
	}{
		{
			name:           "success",
			requestBody:    dto.CreateBoardDTO{Title: "New Board"},
			mockReturn:     model.Board{ID: 1, Title: "New Board"},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "empty title",
			requestBody:    dto.CreateBoardDTO{Title: ""},
			mockReturn:     model.Board{},
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "service error",
			requestBody:    dto.CreateBoardDTO{Title: "Boom"},
			mockReturn:     model.Board{},
			mockError:      errors.New("service error"),
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "error decode",
			rawBody:        `{"title":123}`,
			mockReturn:     model.Board{ID: 123, Title: "EncodeFail"},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:            "error from encode",
			requestBody:     dto.CreateBoardDTO{Title: "EncodeFail"},
			mockReturn:      model.Board{ID: 123, Title: "EncodeFail"},
			expectedStatus:  http.StatusBadRequest,
			expectEncodeErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockBoardStorage)
			logger := zap.NewNop()
			handler := NewBoardHandler(mockService, logger)

			var body []byte
			if tt.rawBody != "" {
				body = []byte(tt.rawBody)
			} else {
				body, _ = json.Marshal(tt.requestBody)
			}

			req := httptest.NewRequest(http.MethodPost, "/boards", bytes.NewReader(body))
			rec := httptest.NewRecorder()

			if len(tt.requestBody.Title) != 0 && !tt.expectEncodeErr {
				mockService.On("CreateBoard", tt.requestBody.Title).Return(tt.mockReturn, tt.mockError)
			}
			if tt.expectEncodeErr {
				mockService.On("CreateBoard", tt.requestBody.Title).Return(model.Board{}, tt.mockError)
			}

			rw := http.ResponseWriter(rec)
			if tt.expectEncodeErr {
				rw = &badResponseWriter{rec}
			}

			handler.HandleBoards(rw, req)

			require.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectedStatus == http.StatusOK && !tt.expectEncodeErr {
				var response dto.BoardDTO
				err := json.NewDecoder(rec.Body).Decode(&response)
				require.NoError(t, err)
				require.Equal(t, tt.mockReturn.ID, *response.ID)
				require.Equal(t, tt.mockReturn.Title, response.Title)
			}
		})
	}
}
func TestHandleBoards_MethodNotAllowed(t *testing.T) {
	mockService := new(MockBoardStorage)
	logger := zap.NewNop()
	handler := NewBoardHandler(mockService, logger)

	req := httptest.NewRequest(http.MethodPut, "/boards", nil)
	rec := httptest.NewRecorder()

	handler.HandleBoards(rec, req)

	require.Equal(t, http.StatusMethodNotAllowed, rec.Code)
	require.Contains(t, rec.Body.String(), "Method not allowed")
}
