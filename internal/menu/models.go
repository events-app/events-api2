package menu

import (
	"time"
)

// Menu is an datastructure for a Menu object.
type Menu struct {
	ID          string    `db:"menu_id" json:"id"`
	Name        string    `db:"name" json:"name"`
	CardID      string    `db:"card_id" json:"cardID"`
	UserID      string    `db:"user_id" json:"userID"`
	DateCreated time.Time `db:"date_created" json:"dateCreated"`
	DateUpdated time.Time `db:"date_updated" json:"dateUpdated"`
}

// NewMenu is what we require from admin when adding a Menu.
type NewMenu struct {
	Name   string `json:"name" validate:"required"`
	CardID string `json:"cardID"`
}

// UpdateMenu defines what information may be provided to modify an
// existing Menu. All fields are optional so clients can send just the
// fields they want changed. It uses pointer fields so we can differentiate
// between a field that was not provided and a field that was provided as
// explicitly blank. Normally we do not want to use pointers to basic types but
// we make exceptions around marshalling/unmarshalling.
type UpdateMenu struct {
	Name   *string `json:"name"`
	CardID *string `json:"cardID"`
}
