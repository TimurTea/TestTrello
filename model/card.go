package model

import "time"

type Card struct {
	ID          int       `db:"id" json:"id"`
	BoardID     int       `db:"board_id" json:"board_id"`
	ListID      int       `db:"list_id" json:"list_id"`
	Title       string    `db:"title" json:"title"`
	Description string    `db:"description" json:"description"`
	Status      string    `db:"status" json:"status"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}
type CardInputCreate struct {
	ListID      int    `db:"list_id" json:"list_id"`
	Title       string `db:"title" json:"title"`
	Description string `db:"description" json:"description"`
}
