package handler

import (
	"awesomeProject2/cmd/dto"
	"awesomeProject2/cmd/model"
	"encoding/json"
	"net/http"
)

type ListHandler struct {
	service ListService
}

func NewListHandler(service ListService) *ListHandler {
	return &ListHandler{service}
}
func (h *ListHandler) HandleLists(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		var requestDTO dto.ListDTO
		defer r.Body.Close()
		if r.Body != nil {
			if err := json.NewDecoder(r.Body).Decode(&requestDTO); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			lists, err := h.service.GetLists(requestDTO.ID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			var listDTOs []dto.ListDTO
			for _, l := range lists {
				listDTOs = append(listDTOs, dto.ListToDTO(l))
			}
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(listDTOs); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	} else if r.Method == http.MethodPost {
		defer r.Body.Close()
		var input dto.CreateListDTO
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if len(input.Title) == 0 {
			http.Error(w, "title is required", http.StatusBadRequest)
			return
		}

		list, err := h.service.CreateList(model.ListInputCreate{
			BoardID: input.BoardID,
			Title:   input.Title,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		dto := dto.ListToDTO(list)
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(dto); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
