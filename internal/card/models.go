package card

import (
	"time"
)

// Card is an datastructure for a Card object.
type Card struct {
	ID          string    `db:"card_id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Content     string    `db:"content" json:"content"`
	UserID      string    `db:"user_id" json:"userID"`
	DateCreated time.Time `db:"date_created" json:"dateCreated"`
	DateUpdated time.Time `db:"date_updated" json:"dateUpdated"`
}

// NewCard is what we require from admin when adding a Card.
type NewCard struct {
	Name     string `json:"name" validate:"required"`
	Content        string       `json:"content"`
}

// UpdateCard defines what information may be provided to modify an
// existing Card. All fields are optional so clients can send just the
// fields they want changed. It uses pointer fields so we can differentiate
// between a field that was not provided and a field that was provided as
// explicitly blank. Normally we do not want to use pointers to basic types but
// we make exceptions around marshalling/unmarshalling.
type UpdateCard struct {
	Name     *string `json:"name"`
	Content        *string       `json:"content"`
}

// // Sale represents one item of a transaction where some amount of a card was
// // sold. Quantity is the number of units sold and Paid is the total price paid.
// // Note that due to haggling the Paid value might not equal Quantity sold *
// // Card cost.
// type Sale struct {
// 	ID          string    `db:"sale_id" json:"id"`
// 	CardID   string    `db:"card_id" json:"card_id"`
// 	Quantity    int       `db:"quantity" json:"quantity"`
// 	Paid        int       `db:"paid" json:"paid"`
// 	DateCreated time.Time `db:"date_created" json:"date_created"`
// }

// // NewSale is what we require from clients for recording new transactions.
// type NewSale struct {
// 	Quantity int `json:"quantity" validate:"gte=0"`
// 	Paid     int `json:"paid" validate:"gte=0"`
// }
