package model

import "time"

type List struct {
	ID        int       `db:"id" json:"id"`
	BoardID   int       `db:"board_id" json:"board_id"`
	Title     string    `db:"title" json:"title"`
	Cards     []Card    `db:"cards" json:"cards"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
type ListInputCreate struct {
	ListID int    `db:"board_id" json:"board_id"`
	Title  string `db:"title" json:"title"`
}
