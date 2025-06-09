package dto

import (
	"awesomeProject2/cmd/model"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBoardToDTO(t *testing.T) {
	pointer := 1
	tests := []struct {
		name  string
		board model.Board
		want  BoardDTO
	}{
		{
			name: "normal board",
			board: model.Board{
				ID:    pointer,
				Title: "Sprint 1",
			},
			want: BoardDTO{
				ID:    &pointer,
				Title: "Sprint 1",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BoardToDTO(tt.board)
			require.NotNil(t, got.ID)
			require.NotNil(t, tt.want)
			require.Equal(t, *got.ID, *tt.want.ID)
			require.Equal(t, tt.want.Title, got.Title)
		})
	}
}

func TestListToDTO(t *testing.T) {
	pointer := 1
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
				ID:      &pointer,
				Title:   "To Do",
				BoardID: 10,
			},
		},
		{
			name: "empty title",
			list: model.List{
				ID:      pointer,
				Title:   "",
				BoardID: 20,
			},
			want: ListDTO{
				ID:      &pointer,
				Title:   "",
				BoardID: 20,
			},
		},
		{
			name: "zero board id",
			list: model.List{
				ID:      pointer,
				Title:   "In Progress",
				BoardID: 0,
			},
			want: ListDTO{
				ID:      &pointer,
				Title:   "In Progress",
				BoardID: 0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ListToDTO(tt.list)
			require.NotNil(t, got.ID)
			require.NotNil(t, tt.want.ID)
			require.Equal(t, *tt.want.ID, *got.ID)
			require.Equal(t, tt.want.Title, got.Title)
			require.Equal(t, tt.want.BoardID, got.BoardID)
		})
	}
}
func TestCardToDTO(t *testing.T) {
	pointer := 1
	tests := []struct {
		name string
		card model.Card
		want CardDTO
	}{
		{
			name: "basic conversion",
			card: model.Card{
				ID:          pointer,
				Title:       "Fix bug",
				Description: "Null pointer exception",
				ListID:      2,
			},
			want: CardDTO{
				ID:          &pointer,
				Title:       "Fix bug",
				Description: "Null pointer exception",
				ListID:      2,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CardToDTO(tt.card)
			require.NotNil(t, got.ID)
			require.NotNil(t, tt.want.ID)
			require.Equal(t, *tt.want.ID, *got.ID)
			require.Equal(t, tt.want.Title, got.Title)
			require.Equal(t, tt.want.BoardID, got.BoardID)
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
