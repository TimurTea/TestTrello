package handler

import (
	"awesomeProject2/cmd/dto"
	"awesomeProject2/cmd/model"
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
)

type CardHandler struct {
	service CardService
	logger  *zap.Logger
}

func NewCardHandler(service CardService, logger *zap.Logger) *CardHandler {
	return &CardHandler{
		service: service,
		logger:  logger,
	}
}
func (h *CardHandler) HandleCards(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		var requestDTO dto.CardDTO
		if r.Body != nil {
			defer r.Body.Close()
			if err := json.NewDecoder(r.Body).Decode(&requestDTO); err != nil {
				h.logger.Error("Ошибка декодирования запроса(GET)", zap.Error(err), zap.Any("requestDTO", requestDTO))
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}
		cards, err := h.service.GetCards(requestDTO.ID)
		if err != nil {
			h.logger.Error("Ошибка получения карточек", zap.Error(err), zap.Any("requestDTO", requestDTO))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var cardDTOs []dto.CardDTO
		for i := range cards {
			cardDTOs = append(cardDTOs, dto.CardToDTO(cards[i]))
		}
		if err := json.NewEncoder(w).Encode(cardDTOs); err != nil {
			h.logger.Error("Ошибка при кодировании ответа", zap.Error(err), zap.Any("requestDTO", requestDTO))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	} else if r.Method == http.MethodPost {
		var input dto.CreateCardDTO
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			h.logger.Error("Ошибка декодирования запрос(POST)", zap.Error(err), zap.Any("requestDTO", input))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if input.ListID == 0 {
			h.logger.Error("Листы отсуствуют", zap.Any("input", input))
			http.Error(w, "list id required", http.StatusBadRequest)
			return
		}
		if len(input.Title) == 0 {
			h.logger.Error("Название отсустует", zap.Any("input", input))
			http.Error(w, "title is required", http.StatusBadRequest)
			return
		}
		card, err := h.service.CreateCard(model.CardInputCreate{
			ListID:      input.ListID,
			Title:       input.Title,
			Description: input.Description,
		})
		if err != nil {
			h.logger.Error("Ошибка создание карточки", zap.Error(err), zap.Any("input", input))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		cardDTOs := dto.CardToDTO(card)
		json.NewEncoder(w).Encode(cardDTOs)
	} else if r.Method == http.MethodDelete {
		var input dto.DeleteCardDTO
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			h.logger.Error("Ошибка декодирования(DELETE)", zap.Error(err), zap.Any("input", input))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if input.ListID == 0 {
			h.logger.Error("Листы отсуствуют", zap.Any("input", input))
			http.Error(w, "list id required", http.StatusBadRequest)
			return
		}
		if input.CardID == 0 {
			h.logger.Error("Карточки отсуствуют", zap.Any("input", input))
			http.Error(w, "card id required", http.StatusBadRequest)
			return
		}
		deletedCard, err := h.service.DeleteCard(input.ListID, input.CardID)
		if err != nil {
			h.logger.Error("Ошибка удаления карточки", zap.Error(err))
			http.Error(w, "Error deleting card", http.StatusInternalServerError)
			return
		}
		cardDTOs := dto.CardToDTO(deletedCard)
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(cardDTOs); err != nil {
			h.logger.Error("Ошибка кодирования запроса(DELETE)", zap.Error(err), zap.Any("cardDTOs", cardDTOs))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	} else if r.Method == http.MethodPut {
		var updatedCardDTO dto.UpdateCardDTO
		if err := json.NewDecoder(r.Body).Decode(&updatedCardDTO); err != nil {
			h.logger.Error("Ошибка декодирования запроса(PUT)", zap.Error(err), zap.Any("updatedCardDTO", updatedCardDTO))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if updatedCardDTO.ListID == 0 {
			h.logger.Error("Листы отсуствуют", zap.Any("input", updatedCardDTO))
			http.Error(w, "list id error", http.StatusBadRequest)
			return
		}
		if updatedCardDTO.ID == 0 {
			h.logger.Error("Карточка отсуствует", zap.Any("input", updatedCardDTO))
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
			h.logger.Error("Ошибка обновление карточки", zap.Error(err), zap.Any("updatedCard", updatedCard))
			http.Error(w, "Error updating card", http.StatusInternalServerError)
			return
		}
		updatedCardDTOResponse := dto.CardToDTO(updatedCard)
		json.NewEncoder(w).Encode(updatedCardDTOResponse)
	} else {
		h.logger.Warn("Метод не поддерживается", zap.Any("method", r.Method))
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
