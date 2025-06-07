package storage

import (
	"awesomeProject2/cmd/model"
	"reflect"
	"testing"
)

func TestStorage_CreateBoard(t *testing.T) {
	type fields struct {
		Boards  map[int]model.Board
		boardID int
	}
	type args struct {
		title string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   model.Board
	}{
		{
			name:   "success",
			fields: fields{Boards: map[int]model.Board{}, boardID: 1},
			args:   args{title: "test"},
			want: model.Board{
				ID:    1,
				Title: "test",
			},
		},
		{
			name:   "error create board",
			fields: fields{Boards: map[int]model.Board{}, boardID: 333},
			args:   args{title: "bad test"},
			want: model.Board{
				ID:    333,
				Title: "bad test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				Boards:  tt.fields.Boards,
				boardID: tt.fields.boardID,
			}
			if got := s.CreateBoard(tt.args.title); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateBoard() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_CreateCard(t *testing.T) {
	type fields struct {
		Boards  map[int]model.Board
		boardID int
		listID  int
		cardID  int
	}
	type args struct {
		title   string
		boardID int
		listID  int
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantID     int
		wantTitle  string
		wantInList bool
	}{
		{
			name: "success",
			fields: fields{
				Boards: map[int]model.Board{
					1: {
						ID:    1,
						Title: "Board 1",
						Lists: []model.List{
							{
								ID:    1,
								Title: "List 1",
								Cards: []model.Card{},
							},
						},
					},
				},
				boardID: 1,
				listID:  1,
				cardID:  1,
			},
			args: args{
				title:   "Test Card",
				boardID: 1,
				listID:  1,
			},
			wantID:     1,
			wantTitle:  "Test Card",
			wantInList: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				Boards:  tt.fields.Boards,
				boardID: tt.fields.boardID,
				listID:  tt.fields.listID,
				cardID:  tt.fields.cardID,
			}

			got := s.CreateCard(tt.args.title, tt.args.boardID, tt.args.listID)

			if got.ID != tt.wantID || got.Title != tt.wantTitle {
				t.Errorf("CreateCard() = %+v, want ID=%d, Title=%q", got, tt.wantID, tt.wantTitle)
			}

			if tt.wantInList {
				board := s.Boards[tt.args.boardID]
				found := false
				for _, list := range board.Lists {
					if list.ID == tt.args.listID {
						for _, card := range list.Cards {
							if card.ID == got.ID && card.Title == got.Title {
								found = true
								break
							}
						}
					}
				}
				if !found {
					t.Errorf("Created card not found in list: %+v", board.Lists)
				}
			}
		})
	}
}

func TestStorage_CreateList(t *testing.T) {
	type fields struct {
		Boards  map[int]model.Board
		boardID int
		listID  int
	}
	type args struct {
		title   string
		boardID int
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantID      int
		wantTitle   string
		wantInBoard bool
	}{
		{
			name: "success",
			fields: fields{
				Boards: map[int]model.Board{
					1: {
						ID:    1,
						Title: "Test Board",
						Lists: []model.List{},
					},
				},
				boardID: 1,
				listID:  1,
			},
			args: args{
				title:   "test",
				boardID: 1,
			},
			wantID:      1,
			wantTitle:   "test",
			wantInBoard: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				Boards:  tt.fields.Boards,
				boardID: tt.fields.boardID,
				listID:  tt.fields.listID,
			}

			got := s.CreateList(tt.args.title, tt.args.boardID)

			if got.ID != tt.wantID || got.Title != tt.wantTitle {
				t.Errorf("CreateList() = %+v, want ID=%d, Title=%q", got, tt.wantID, tt.wantTitle)
			}

			if tt.wantInBoard {
				board := s.Boards[tt.args.boardID]
				found := false
				for _, list := range board.Lists {
					if list.ID == got.ID && list.Title == got.Title {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Created list not found in board's list slice: got %+v", board.Lists)
				}
			}
		})
	}
}

func TestStorage_GetCards(t *testing.T) {
	type fields struct {
		Boards  map[int]model.Board
		boardID int
		listID  int
		cardID  int
	}
	type args struct {
		cardID *int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []model.Card
	}{
		{
			name: "success",
			fields: fields{
				Boards: map[int]model.Board{
					1: {
						ID: 1,
						Lists: []model.List{
							{
								ID: 1,
								Cards: []model.Card{
									{ID: 1, Title: "Test Card"},
									{ID: 2, Title: "Test Card"},
								},
							},
						},
					},
				},
				boardID: 1,
				listID:  1,
			},
			args: args{cardID: nil},
			want: []model.Card{
				{ID: 1, Title: "Test Card"},
				{ID: 2, Title: "Test Card"},
			},
		},
		{
			name: "error",
			fields: fields{
				Boards: map[int]model.Board{
					1: {
						ID: 1,
						Lists: []model.List{
							{
								ID: 1,
								Cards: []model.Card{
									{ID: 1, Title: "Test Card"},
									{ID: 2, Title: "Test Card"},
								},
							},
						},
					},
				},
				boardID: 1,
				listID:  1,
			},
			args: func() args {
				id := 999
				return args{cardID: &id}
			}(),
			want: []model.Card{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				Boards:  tt.fields.Boards,
				boardID: tt.fields.boardID,
				listID:  tt.fields.listID,
				cardID:  tt.fields.cardID,
			}
			if got := s.GetCards(tt.args.cardID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCards() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_GetLists(t *testing.T) {
	type fields struct {
		Boards  map[int]model.Board
		boardID int
		listID  int
		cardID  int
	}
	type args struct {
		listID *int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []model.List
	}{
		{
			name: "success",
			fields: fields{
				Boards: map[int]model.Board{
					1: {
						ID: 1,
						Lists: []model.List{
							{ID: 1, Title: "Test List"},
							{ID: 2, Title: "Test List"},
						},
					},
				},
				boardID: 1,
			},
			args: args{listID: nil},
			want: []model.List{
				{ID: 1, Title: "Test List"},
				{ID: 2, Title: "Test List"},
			},
		},
		{
			name: "error",
			fields: fields{
				Boards: map[int]model.Board{
					1: {
						ID: 1,
						Lists: []model.List{
							{ID: 1, Title: "Test List"},
							{ID: 2, Title: "Test List"},
						},
					},
				},
				boardID: 1,
			},
			args: func() args {
				id := 999
				return args{listID: &id}
			}(),
			want: []model.List{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				Boards:  tt.fields.Boards,
				boardID: tt.fields.boardID,
				listID:  tt.fields.listID,
				cardID:  tt.fields.cardID,
			}
			if got := s.GetLists(tt.args.listID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetLists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_GetBoards(t *testing.T) {
	type fields struct {
		Boards  map[int]model.Board
		boardID int
		listID  int
		cardID  int
	}
	tests := []struct {
		name   string
		fields fields
		want   []model.Board
	}{
		{
			name: "success",
			fields: fields{
				Boards: map[int]model.Board{
					1: {ID: 1, Title: "Test Board"},
				},
			},
			want: []model.Board{
				{ID: 1, Title: "Test Board"},
			},
		},
		{
			name: "error",
			fields: fields{
				Boards: map[int]model.Board{
					555: {ID: 555, Title: "Test Board"},
				},
			},
			want: []model.Board{
				{ID: 555, Title: "Test Board"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				Boards:  tt.fields.Boards,
				boardID: tt.fields.boardID,
				listID:  tt.fields.listID,
				cardID:  tt.fields.cardID,
			}
			if got := s.GetBoards(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetBoards() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_DeleteCard(t *testing.T) {
	type fields struct {
		Boards  map[int]model.Board
		boardID int
		listID  int
		cardID  int
	}
	type args struct {
		boardID int
		listID  int
		cardID  int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    model.Card
		wantErr bool
	}{
		{
			name: "successfully deletes card",
			fields: fields{
				Boards: map[int]model.Board{
					1: {
						ID: 1,
						Lists: []model.List{
							{
								ID: 1,
								Cards: []model.Card{
									{ID: 1, Title: "Test Card"},
									{ID: 2, Title: "Second Card"},
								},
							},
						},
					},
				},
			},
			args: args{
				boardID: 1,
				listID:  1,
				cardID:  1,
			},
			want:    model.Card{ID: 1, Title: "Test Card"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				Boards:  tt.fields.Boards,
				boardID: tt.fields.boardID,
				listID:  tt.fields.listID,
				cardID:  tt.fields.cardID,
			}
			got, err := s.DeleteCard(tt.args.boardID, tt.args.listID, tt.args.cardID)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteCard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteCard() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_UpdateCard(t *testing.T) {
	type fields struct {
		Boards  map[int]model.Board
		boardID int
		listID  int
		cardID  int
	}
	type args struct {
		updated model.Card
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    model.Card
		wantErr bool
	}{
		{
			name: "successfully updates card",
			fields: fields{
				Boards: map[int]model.Board{
					1: {
						ID: 1,
						Lists: []model.List{
							{
								ID: 1,
								Cards: []model.Card{
									{ID: 1, Title: "Old Title", Description: "Old Desc", ListID: 1, BoardID: 1},
								},
							},
						},
					},
				},
			},
			args: args{
				updated: model.Card{ID: 1, Title: "New Title", Description: "New Desc", ListID: 1, BoardID: 1},
			},
			want:    model.Card{ID: 1, Title: "New Title", Description: "New Desc", ListID: 1, BoardID: 1},
			wantErr: false,
		},
		{
			name: "card not found",
			fields: fields{
				Boards: map[int]model.Board{
					1: {
						ID: 1,
						Lists: []model.List{
							{
								ID: 1,
								Cards: []model.Card{
									{ID: 2, Title: "Other Card", ListID: 1, BoardID: 1},
								},
							},
						},
					},
				},
			},
			args: args{
				updated: model.Card{ID: 999, Title: "Doesn't Exist", ListID: 1, BoardID: 1},
			},
			want:    model.Card{},
			wantErr: true,
		},
		{
			name: "list not found",
			fields: fields{
				Boards: map[int]model.Board{
					1: {
						ID:    1,
						Lists: []model.List{},
					},
				},
			},
			args: args{
				updated: model.Card{ID: 1, Title: "Any", ListID: 999, BoardID: 1},
			},
			want:    model.Card{},
			wantErr: true,
		},
		{
			name: "board not found",
			fields: fields{
				Boards: map[int]model.Board{},
			},
			args: args{
				updated: model.Card{ID: 1, Title: "Any", ListID: 1, BoardID: 999},
			},
			want:    model.Card{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				Boards:  tt.fields.Boards,
				boardID: tt.fields.boardID,
				listID:  tt.fields.listID,
				cardID:  tt.fields.cardID,
			}
			got, err := s.UpdateCard(tt.args.updated)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateCard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateCard() got = %v, want %v", got, tt.want)
			}
		})
	}
}
