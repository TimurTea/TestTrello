package dto

import (
	"awesomeProject2/cmd/helper"
	"awesomeProject2/cmd/model"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBoardToDTO(t *testing.T) {
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
				ID:    helper.GetPointer(1),
				Title: "Sprint 1",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BoardToDTO(tt.board)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestListToDTO(t *testing.T) {
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
				ID:      helper.GetPointer(1),
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
				ID:      helper.GetPointer(1),
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
				ID:      helper.GetPointer(1),
				Title:   "In Progress",
				BoardID: 0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ListToDTO(tt.list)
			require.Equal(t, tt.want, got)
		})
	}
}
func TestCardToDTO(t *testing.T) {
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
				ID:          helper.GetPointer(1),
				Title:       "Fix bug",
				Description: "Null pointer exception",
				ListID:      2,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CardToDTO(tt.card)
			require.Equal(t, tt.want, got)
		})
	}
}
