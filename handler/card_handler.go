package handler

import (
	"awesomeProject2/cmd/dto"
	"awesomeProject2/cmd/model"
	"encoding/json"
	"net/http"
)

type CardHandler struct {
	service CardService
}

func NewCardHandler(service CardService) *CardHandler {
	return &CardHandler{service}
}
func (h *CardHandler) HandleCards(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		var requestDTO dto.CardDTO
		if r.Body != nil {
			defer r.Body.Close()
			if err := json.NewDecoder(r.Body).Decode(&requestDTO); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}
		cards, err := h.service.GetCards(requestDTO.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		var cardDTOs []dto.CardDTO
		for i := range cards {
			cardDTOs = append(cardDTOs, dto.CardToDTO(cards[i]))
		}
		if err := json.NewEncoder(w).Encode(cardDTOs); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	} else if r.Method == http.MethodPost {
		var input dto.CreateCardDTO
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if input.ListID == 0 {
			http.Error(w, "list id required", http.StatusBadRequest)
		}
		if len(input.Title) == 0 {
			http.Error(w, "title is required", http.StatusBadRequest)
			return
		}
		card, err := h.service.CreateCard(model.CardInputCreate{
			ListID:      input.ListID,
			Title:       input.Title,
			Description: input.Description,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		cardDTOs := dto.CardToDTO(card)
		json.NewEncoder(w).Encode(cardDTOs)
	} else if r.Method == http.MethodDelete {
		var input dto.DeleteCardDTO
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		if input.ListID == 0 {
			http.Error(w, "list id required", http.StatusBadRequest)
		}
		if input.CardID == 0 {
			http.Error(w, "card id required", http.StatusBadRequest)
		}
		deletedCard, err := h.service.DeleteCard(input.ListID, input.CardID)
		if err != nil {
			http.Error(w, "Error deleting card", http.StatusInternalServerError)
			return
		}
		cardDTOs := dto.CardToDTO(deletedCard)
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(cardDTOs); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	} else if r.Method == http.MethodPut {
		var updatedCardDTO dto.UpdateCardDTO
		if err := json.NewDecoder(r.Body).Decode(&updatedCardDTO); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if updatedCardDTO.ListID == 0 {
			http.Error(w, "list id error", http.StatusBadRequest)
			return
		}
		if updatedCardDTO.ID == 0 {
			http.Error(w, "card id error", http.StatusBadRequest)
			return
		}
		updatedCard := model.Card{
			Title:       updatedCardDTO.Title,
			Description: updatedCardDTO.Description,
			ID:          updatedCardDTO.ID,
			BoardID:     updatedCardDTO.BoardID,
			ListID:      updatedCardDTO.ListID,
		}
		updatedCard, err := h.service.UpdateCard(updatedCard)
		if err != nil {
			http.Error(w, "Error updating card", http.StatusInternalServerError)
			return
		}
		updatedCardDTOResponse := dto.CardToDTO(updatedCard)
		json.NewEncoder(w).Encode(updatedCardDTOResponse)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
