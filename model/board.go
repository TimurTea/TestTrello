package model

import "time"

type Board struct {
	ID         int       `db:"id" json:"id"`
	Title      string    `db:"title" json:"title"`
	Lists      []List    `db:"lists" json:"lists"`
	NextListID int       `db:"next_list_id" json:"next_list_id"`
	NextCardID int       `db:"next_card_id" json:"next_card_id"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}
