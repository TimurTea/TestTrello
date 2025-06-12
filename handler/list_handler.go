package handler

import (
	"awesomeProject2/cmd/dto"
	"awesomeProject2/cmd/model"
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
)

type ListHandler struct {
	service ListService
	logger  *zap.Logger
}

func NewListHandler(service ListService, logger *zap.Logger) *ListHandler {
	return &ListHandler{
		service: service,
		logger:  logger,
	}
}
func (h *ListHandler) HandleLists(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		var requestDTO dto.ListDTO
		defer r.Body.Close()
		if r.Body != nil {
			if err := json.NewDecoder(r.Body).Decode(&requestDTO); err != nil {
				h.logger.Error("Ошибка декодирования запроса(GET)", zap.Error(err), zap.Any("dto", requestDTO))
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			lists, err := h.service.GetLists(requestDTO.ID)
			if err != nil {
				h.logger.Error("Ошибка получения листов", zap.Error(err), zap.Any("dto", requestDTO))
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			var listDTOs []dto.ListDTO
			for _, l := range lists {
				listDTOs = append(listDTOs, dto.ListToDTO(l))
			}
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(listDTOs); err != nil {
				h.logger.Error("Ошибка кодирования запроса(GET)", zap.Error(err), zap.Any("listDTOs", listDTOs))
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	} else if r.Method == http.MethodPost {
		defer r.Body.Close()
		var input dto.CreateListDTO
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			h.logger.Error("Ошибка декодирования запроса(POST)", zap.Error(err), zap.Any("dto", input))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if len(input.Title) == 0 {
			h.logger.Error("Отсуствуют названия", zap.Any("input", input))
			http.Error(w, "title is required", http.StatusBadRequest)
			return
		}

		list, err := h.service.CreateList(model.ListInputCreate{
			BoardID: input.BoardID,
			Title:   input.Title,
		})
		if err != nil {
			h.logger.Error("Ошибка создания листов", zap.Error(err), zap.Any("input", input))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		dto := dto.ListToDTO(list)
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(dto); err != nil {
			h.logger.Error("Ошибка кодирования запроса(POST)", zap.Error(err), zap.Any("dto", dto))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	} else {
		h.logger.Warn("Метод отсуствует", zap.Any("method", r.Method))
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}
