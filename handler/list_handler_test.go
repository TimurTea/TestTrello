package handler

import (
	"awesomeProject2/cmd/dto"
	"awesomeProject2/cmd/helper"
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

func TestHandleLists_GET(t *testing.T) {
	tests := []struct {
		name            string
		requestBody     dto.ListDTO
		serviceResponse []model.List
		mockError       error
		expectedStatus  int
	}{
		{
			name:            "success with cards",
			requestBody:     dto.ListDTO{ID: helper.GetPointer(1)},
			serviceResponse: []model.List{{ID: 1, Title: "Card1"}, {ID: 2, Title: "Card2"}},
			mockError:       nil,
			expectedStatus:  http.StatusOK,
		},
		{
			name:            "empty cards list",
			requestBody:     dto.ListDTO{ID: helper.GetPointer(1)},
			serviceResponse: []model.List{},
			mockError:       nil,
			expectedStatus:  http.StatusOK,
		},
		{
			name:            "service error",
			requestBody:     dto.ListDTO{ID: helper.GetPointer(1)},
			serviceResponse: nil,
			mockError:       errors.New("fail"),
			expectedStatus:  http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := new(MockListStorage)
			logger := zap.NewNop()
			handler := NewListHandler(mock, logger)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodGet, "/lists", bytes.NewReader(body))
			rec := httptest.NewRecorder()

			mock.On("GetLists", tt.requestBody.ID).Return(tt.serviceResponse, tt.mockError)

			handler.HandleLists(rec, req)

			require.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectedStatus == http.StatusOK {
				var resp []dto.ListDTO
				err := json.NewDecoder(rec.Body).Decode(&resp)
				require.NoError(t, err)
				require.Equal(t, len(tt.serviceResponse), len(resp))
				for i, list := range tt.serviceResponse {
					require.Equal(t, list.ID, *resp[i].ID)
					require.Equal(t, list.Title, resp[i].Title)
				}
			}
		})
	}
}
func TestHandleLists_POST(t *testing.T) {
	tests := []struct {
		name            string
		requestBody     dto.CreateListDTO
		serviceResponse model.List
		mockError       error
		expectedStatus  int
		rawBody         string
		expectEncodeErr bool
	}{
		{
			name:            "success",
			requestBody:     dto.CreateListDTO{Title: "New List"},
			serviceResponse: model.List{ID: 1, Title: "New List"},
			mockError:       nil,
			expectedStatus:  http.StatusOK,
		},
		{
			name:            "empty title",
			requestBody:     dto.CreateListDTO{Title: ""},
			serviceResponse: model.List{},
			mockError:       nil,
			expectedStatus:  http.StatusBadRequest,
		},
		{
			name:            "service error",
			requestBody:     dto.CreateListDTO{Title: "Boom"},
			serviceResponse: model.List{},
			mockError:       errors.New("service error"),
			expectedStatus:  http.StatusInternalServerError,
		},
		{
			name:            "invalid json",
			requestBody:     dto.CreateListDTO{},
			serviceResponse: model.List{},
			mockError:       nil,
			expectedStatus:  http.StatusBadRequest,
		},
		{
			name:            "error decode",
			rawBody:         `{"title":123}`,
			serviceResponse: model.List{ID: 123, Title: "EncodeFail"},
			expectedStatus:  http.StatusBadRequest,
		},
		{
			name:            "error from encode",
			requestBody:     dto.CreateListDTO{Title: "EncodeFail"},
			serviceResponse: model.List{ID: 123, Title: "EncodeFail"},
			expectedStatus:  http.StatusInternalServerError,
			expectEncodeErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockListStorage)
			logger := zap.NewNop()
			handler := NewListHandler(mockService, logger)

			var req *http.Request
			if tt.rawBody != "" {
				req = httptest.NewRequest(http.MethodPost, "/lists", bytes.NewBufferString(tt.rawBody))
			} else {
				body, _ := json.Marshal(tt.requestBody)
				req = httptest.NewRequest(http.MethodPost, "/lists", bytes.NewReader(body))
			}

			rec := httptest.NewRecorder()
			rw := http.ResponseWriter(rec)

			if tt.expectEncodeErr {
				rw = &badResponseWriter{rec}
			}

			if tt.requestBody.Title != "" {
				input := model.ListInputCreate{Title: tt.requestBody.Title}
				mockService.On("CreateList", input).Return(tt.serviceResponse, tt.mockError)
			}

			handler.HandleLists(rw, req)

			require.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectedStatus == http.StatusOK && !tt.expectEncodeErr {
				var resp dto.ListDTO
				err := json.NewDecoder(rec.Body).Decode(&resp)
				require.NoError(t, err)

				expected := dto.ListToDTO(tt.serviceResponse)
				require.Equal(t, expected, resp)
			}
		})
	}
}
func TestHandleLists_MethodNotAllowed(t *testing.T) {
	mockService := new(MockListStorage)
	logger := zap.NewNop()
	handler := NewListHandler(mockService, logger)

	req := httptest.NewRequest(http.MethodPut, "/lists", nil)
	rec := httptest.NewRecorder()

	handler.HandleLists(rec, req)

	require.Equal(t, http.StatusMethodNotAllowed, rec.Code)
	require.Contains(t, rec.Body.String(), "Method not allowed")
}
