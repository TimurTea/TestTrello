package handler

import (
	"awesomeProject2/cmd/dto"
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
)

type BoardHandler struct {
	service BoardService
	logger  *zap.Logger
}

func NewBoardHandler(service BoardService, logger *zap.Logger) *BoardHandler {
	return &BoardHandler{
		service: service,
		logger:  logger,
	}
}
func (h *BoardHandler) HandleBoards(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		var requestDTO dto.BoardDTO
		if r.Body != nil {
			defer r.Body.Close()
			if err := json.NewDecoder(r.Body).Decode(&requestDTO); err != nil {
				h.logger.Error("Ошибка декодирования запроса(GET)", zap.Error(err))
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}
		boards, err := h.service.GetBoards()
		if err != nil {
			h.logger.Error("Ошибка получение досок", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var boardDTOs []dto.BoardDTO
		for _, b := range boards {
			boardDTOs = append(boardDTOs, dto.BoardToDTO(b))
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(boardDTOs); err != nil {
			h.logger.Error("Ошибка кодирования запроса(GET)", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else if r.Method == http.MethodPost {
		var input dto.CreateBoardDTO
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			h.logger.Error("Ошибка декодирования запроса(POST)", zap.Error(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if len(input.Title) == 0 {
			h.logger.Error("Отсуствуют названия", zap.Any("input", input))
			http.Error(w, "title is required", http.StatusBadRequest)
			return
		}
		board, err := h.service.CreateBoard(input.Title)
		if err != nil {
			h.logger.Error("Ошибка создания доски", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		dto := dto.BoardToDTO(board)
		if err := json.NewEncoder(w).Encode(dto); err != nil {
			h.logger.Warn("Ошибка кодирования запроса(POST)", zap.Error(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
