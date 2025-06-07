package dto

import (
	"awesomeProject2/cmd/model"
	"testing"
)

func TestBoardToDTO(t *testing.T) {
	ptr := 1
	tests := []struct {
		name  string
		board model.Board
		want  BoardDTO
	}{
		{
			name: "normal board",
			board: model.Board{
				ID:    1,
				Title: "Sprint 1",
			},
			want: BoardDTO{
				ID:    &ptr,
				Title: "Sprint 1",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BoardToDTO(tt.board)
			if got.ID == nil || *got.ID != *tt.want.ID || got.Title != tt.want.Title {
				t.Errorf("BoardToDTO() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestListToDTO(t *testing.T) {
	ptr := 1
	tests := []struct {
		name string
		list model.List
		want ListDTO
	}{
		{
			name: "basic conversion",
			list: model.List{
				ID:      1,
				Title:   "To Do",
				BoardID: 10,
			},
			want: ListDTO{
				ID:      &ptr,
				Title:   "To Do",
				BoardID: 10,
			},
		},
		{
			name: "empty title",
			list: model.List{
				ID:      1,
				Title:   "",
				BoardID: 20,
			},
			want: ListDTO{
				ID:      &ptr,
				Title:   "",
				BoardID: 20,
			},
		},
		{
			name: "zero board id",
			list: model.List{
				ID:      1,
				Title:   "In Progress",
				BoardID: 0,
			},
			want: ListDTO{
				ID:      &ptr,
				Title:   "In Progress",
				BoardID: 0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ListToDTO(tt.list)

			if got.ID == nil || *got.ID != *tt.want.ID ||
				got.Title != tt.want.Title ||
				got.BoardID != tt.want.BoardID {
				t.Errorf("ListToDTO() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
func TestCardToDTO(t *testing.T) {
	ptr := 1
	tests := []struct {
		name string
		card model.Card
		want CardDTO
	}{
		{
			name: "basic conversion",
			card: model.Card{
				ID:          1,
				Title:       "Fix bug",
				Description: "Null pointer exception",
				ListID:      2,
			},
			want: CardDTO{
				ID:          &ptr,
				Title:       "Fix bug",
				Description: "Null pointer exception",
				ListID:      2,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CardToDTO(tt.card)
			if got.ID == nil && tt.want.ID != nil {
				t.Errorf("CardToDTO() ID = nil, want %v", *tt.want.ID)
				return
			}
			if got.ID != nil && tt.want.ID == nil {
				t.Errorf("CardToDTO() ID = %v, want nil", *got.ID)
				return
			}
			if got.ID != nil && tt.want.ID != nil && *got.ID != *tt.want.ID {
				t.Errorf("CardToDTO() ID = %v, want %v", *got.ID, *tt.want.ID)
			}

			if got.Title != tt.want.Title {
				t.Errorf("CardToDTO() Title = %v, want %v", got.Title, tt.want.Title)
			}
			if got.Description != tt.want.Description {
				t.Errorf("CardToDTO() Description = %v, want %v", got.Description, tt.want.Description)
			}
			if got.ListID != tt.want.ListID {
				t.Errorf("CardToDTO() ListID = %v, want %v", got.ListID, tt.want.ListID)
			}

		})
	}
}
