package handler

import (
	"awesomeProject2/cmd/dto"
	"encoding/json"
	"net/http"
)

type BoardHandler struct {
	service BoardService
}

func NewBoardHandler(service BoardService) *BoardHandler {
	return &BoardHandler{service}
}
func (h *BoardHandler) HandleBoards(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		var requestDTO dto.BoardDTO
		if r.Body != nil {
			defer r.Body.Close()
			if err := json.NewDecoder(r.Body).Decode(&requestDTO); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}
		boards, err := h.service.GetBoards()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var boardDTOs []dto.BoardDTO
		for _, b := range boards {
			boardDTOs = append(boardDTOs, dto.BoardToDTO(b))
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(boardDTOs); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else if r.Method == http.MethodPost {
		var input dto.CreateBoardDTO
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if len(input.Title) == 0 {
			http.Error(w, "title is required", http.StatusBadRequest)
			return
		}
		board, err := h.service.CreateBoard(input.Title)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		dto := dto.BoardToDTO(board)
		if err := json.NewEncoder(w).Encode(dto); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
