package storage

import (
	"awesomeProject2/cmd/model"
	"fmt"
	"maps"
	"slices"
)

type Storage struct {
	Boards  map[int]model.Board
	boardID int
	listID  int
	cardID  int
}

func (s *Storage) GetBoards() []model.Board {
	return slices.Collect(maps.Values(s.Boards))
}
func (s *Storage) CreateBoard(title string) model.Board {
	board := model.Board{
		Title: title,
		ID:    s.boardID,
	}
	s.Boards[s.boardID] = board
	s.boardID++
	return board
}
func (s *Storage) GetLists(listID *int) []model.List {
	if listID == nil {
		var result []model.List
		for _, b := range s.Boards {
			result = append(result, b.Lists...)
		}
		return result

	}
	for _, b := range s.Boards {
		for _, l := range b.Lists {
			if l.ID == *listID {
				return []model.List{l}
			}
		}
	}
	return []model.List{}
}
func (s *Storage) CreateList(title string, boardID int) model.List {
	newList := model.List{
		ID:      s.listID,
		Title:   title,
		BoardID: boardID,
	}
	s.listID++
	board, ok := s.Boards[boardID]
	if !ok {
		return model.List{}
	}
	board.Lists = append(board.Lists, newList)
	s.Boards[boardID] = board
	return newList
}
func (s *Storage) GetCards(cardID *int) []model.Card {
	if cardID == nil {
		var result []model.Card
		for _, b := range s.Boards {
			for _, l := range b.Lists {
				result = append(result, l.Cards...)
			}
		}
		return result
	}
	for _, b := range s.Boards {
		for _, l := range b.Lists {
			for _, c := range l.Cards {
				if c.ID == *cardID {
					return []model.Card{c}
				}
			}
		}
	}
	return []model.Card{}
}

func (s *Storage) CreateCard(title string, boardID int, listID int) model.Card {
	newCard := model.Card{
		Title:   title,
		ID:      s.cardID,
		BoardID: boardID,
		ListID:  listID,
	}
	newCard.ID = s.cardID
	s.cardID++
	for i := range s.Boards {
		if s.Boards[i].ID == boardID {
			for j := range s.Boards[i].Lists {
				if s.Boards[i].Lists[j].ID == listID {
					s.Boards[i].Lists[j].Cards = append(s.Boards[i].Lists[j].Cards, newCard)
					return newCard
				}
			}
		}
	}
	return model.Card{}
}
func (s *Storage) DeleteCard(boardID int, listID int, cardID int) (model.Card, error) {
	for i := range s.Boards {
		if s.Boards[i].ID == boardID {
			for j := range s.Boards[i].Lists {
				if s.Boards[i].Lists[j].ID == listID {
					cards := s.Boards[i].Lists[j].Cards
					for k, c := range cards {
						if c.ID == cardID {
							s.Boards[i].Lists[j].Cards = append(cards[:k], cards[k+1:]...)
							return c, nil
						}
					}
					return model.Card{}, fmt.Errorf("card %d not found in list %d", cardID, listID)
				}
			}
			return model.Card{}, fmt.Errorf("list %d not found in board %d", listID, boardID)
		}
	}
	return model.Card{}, fmt.Errorf("board %d not found", boardID)
}
func (s *Storage) UpdateCard(updated model.Card) (model.Card, error) {
	for i := range s.Boards {
		if s.Boards[i].ID == updated.BoardID {
			for j := range s.Boards[i].Lists {
				if s.Boards[i].Lists[j].ID == updated.ListID {
					for c := range s.Boards[i].Lists[j].Cards {
						if s.Boards[i].Lists[j].Cards[c].ID == updated.ID {
							s.Boards[i].Lists[j].Cards[c].Title = updated.Title
							s.Boards[i].Lists[j].Cards[c].Description = updated.Description
							return s.Boards[i].Lists[j].Cards[c], nil
						}
					}
					return model.Card{}, fmt.Errorf("card %d not found in list %d", updated.ID, updated.ListID)
				}
			}
			return model.Card{}, fmt.Errorf("list %d not found in board %d", updated.ListID, updated.BoardID)
		}
	}
	return model.Card{}, fmt.Errorf("board %d not found", updated.BoardID)
}
