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
		name            string
		serviceResponse []model.Board
		mockError       error
		expectedStatus  int
		expectedCount   int
	}{
		{
			name: "success",
			serviceResponse: []model.Board{
				{ID: 1, Title: "Board 1"},
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedCount:  1,
		},
		{
			name:            "error from service",
			serviceResponse: nil,
			mockError:       errors.New("storage error"),
			expectedStatus:  http.StatusInternalServerError,
			expectedCount:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockBoardService)
			logger := zap.NewNop()
			handler := NewBoardHandler(mockService, logger)

			mockService.On("GetBoards").Return(tt.serviceResponse, tt.mockError)

			req := httptest.NewRequest(http.MethodGet, "/boards", nil)
			rec := httptest.NewRecorder()

			handler.HandleBoards(rec, req)

			require.Equal(t, tt.expectedStatus, rec.Code)
			expected := make([]dto.BoardDTO, 0, len(tt.serviceResponse))
			for _, b := range tt.serviceResponse {
				expected = append(expected, dto.BoardToDTO(b))
			}
			if tt.expectedStatus == http.StatusOK {
				var response []dto.BoardDTO
				err := json.NewDecoder(rec.Body).Decode(&response)
				require.NoError(t, err)
				require.Equal(t, expected, response)
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
		serviceResponse model.Board
		mockError       error
		expectedStatus  int
		expectEncodeErr bool
		setupMock       func(service *MockBoardService)
	}{
		{
			name:            "success",
			requestBody:     dto.CreateBoardDTO{Title: "New Board"},
			serviceResponse: model.Board{ID: 1, Title: "New Board"},
			mockError:       nil,
			expectedStatus:  http.StatusOK,
			setupMock: func(s *MockBoardService) {
				s.On("CreateBoard", "New Board").
					Return(model.Board{ID: 1, Title: "New Board"}, nil)
			},
		},
		{
			name:           "empty title",
			requestBody:    dto.CreateBoardDTO{Title: ""},
			expectedStatus: http.StatusBadRequest,
			setupMock:      func(s *MockBoardService) {}, // мок не нужен, но нужен stub
		},
		{
			name:           "service error",
			requestBody:    dto.CreateBoardDTO{Title: "Boom"},
			mockError:      errors.New("service error"),
			expectedStatus: http.StatusInternalServerError,
			setupMock: func(s *MockBoardService) {
				s.On("CreateBoard", "Boom").
					Return(model.Board{}, errors.New("service error"))
			},
		},
		{
			name:           "error decode",
			rawBody:        `{"title":123}`,
			expectedStatus: http.StatusBadRequest,
			setupMock:      func(s *MockBoardService) {},
		},
		{
			name:            "error from encode",
			requestBody:     dto.CreateBoardDTO{Title: "EncodeFail"},
			serviceResponse: model.Board{ID: 123, Title: "EncodeFail"},
			expectedStatus:  http.StatusBadRequest,
			expectEncodeErr: true,
			setupMock: func(s *MockBoardService) {
				s.On("CreateBoard", "EncodeFail").
					Return(model.Board{ID: 123, Title: "EncodeFail"}, nil)
			},
		},
		{
			name:            "success with setupMock",
			requestBody:     dto.CreateBoardDTO{Title: "Test Board"},
			serviceResponse: model.Board{ID: 1, Title: "Test Board"},
			expectedStatus:  http.StatusOK,
			setupMock: func(s *MockBoardService) {
				s.On("CreateBoard", "Test Board").
					Return(model.Board{ID: 1, Title: "Test Board"}, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockBoardService)
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

			tt.setupMock(mockService)

			rw := http.ResponseWriter(rec)
			if tt.expectEncodeErr {
				rw = &badResponseWriter{rec}
			}

			handler.HandleBoards(rw, req)
			require.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectedStatus == http.StatusOK && !tt.expectEncodeErr {
				var response dto.BoardDTO
				err := json.NewDecoder(rec.Body).Decode(&response)
				expected := dto.BoardDTO{
					ID:    &tt.serviceResponse.ID,
					Title: tt.serviceResponse.Title,
				}
				require.NoError(t, err)
				require.Equal(t, expected, response)
			}
		})
	}
}
func TestHandleBoards_MethodNotAllowed(t *testing.T) {
	mockService := new(MockBoardService)
	logger := zap.NewNop()
	handler := NewBoardHandler(mockService, logger)

	req := httptest.NewRequest(http.MethodPut, "/boards", nil)
	rec := httptest.NewRecorder()

	handler.HandleBoards(rec, req)

	require.Equal(t, http.StatusMethodNotAllowed, rec.Code)
	require.Contains(t, rec.Body.String(), "Method not allowed")
}
