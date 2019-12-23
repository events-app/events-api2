package card

import (
	"context"
	"database/sql"
	"time"

	"github.com/events-app/events-api2/internal/platform/auth"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// Predefined errors identify expected failure conditions.
var (
	// ErrNotFound is used when a specific Card is requested but does not exist.
	ErrNotFound = errors.New("card not found")

	// ErrInvalidID is used when an invalid UUID is provided.
	ErrInvalidID = errors.New("ID is not in its proper form")

	// ErrForbidden occurs when a user tries to do something that is forbidden to
	// them according to our access control policies.
	ErrForbidden = errors.New("Attempted action is not allowed")
)

// List gets all Cards from the database.
func List(ctx context.Context, db *sqlx.DB) ([]Card, error) {
	cards := []Card{}
	// const q = `SELECT
	// 		p.*,
	// 		COALESCE(SUM(s.quantity) ,0) AS sold,
	// 		COALESCE(SUM(s.paid), 0) AS revenue
	// 	FROM cards AS c
	// 	LEFT JOIN sales AS s ON c.card_id = c.card_id
	// 	GROUP BY p.card_id`
	const q = `SELECT * FROM cards`
	if err := db.SelectContext(ctx, &cards, q); err != nil {
		return nil, errors.Wrap(err, "selecting cards")
	}

	return cards, nil
}

// Create adds a Card to the database. It returns the created Card with
// fields like ID and DateCreated populated..
func Create(ctx context.Context, db *sqlx.DB, user auth.Claims, nc NewCard, now time.Time) (*Card, error) {
	c := Card{
		ID:          uuid.New().String(),
		Name:        nc.Name,
		Content:     nc.Content,
		UserID:      user.Subject,
		DateCreated: now.UTC(),
		DateUpdated: now.UTC(),
	}

	const q = `
		INSERT INTO cards
		(card_id, user_id, name, content, date_created, date_updated)
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := db.ExecContext(ctx, q,
		c.ID, c.UserID,
		c.Name, c.Content,
		c.DateCreated, c.DateUpdated)
	if err != nil {
		return nil, errors.Wrap(err, "inserting card")
	}

	return &c, nil
}

// Retrieve finds the card identified by a given ID.
func Retrieve(ctx context.Context, db *sqlx.DB, id string) (*Card, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, ErrInvalidID
	}

	var c Card

	const q = `SELECT * FROM cards WHERE card_id = $1`
	if err := db.GetContext(ctx, &c, q, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}

		return nil, errors.Wrap(err, "selecting single card")
	}

	return &c, nil
}

// Update modifies data about a Card. It will error if the specified ID is
// invalid or does not reference an existing Card.
func Update(ctx context.Context, db *sqlx.DB, id string, update UpdateCard, now time.Time) error {
	c, err := Retrieve(ctx, db, id)
	if err != nil {
		return err
	}

	// // If you do not have the admin role ...
	// // and you are not the owner of this card ...
	// // then get outta here!
	// if !user.HasRole(auth.RoleAdmin) && c.UserID != user.Subject {
	// 	return ErrForbidden
	// }

	if update.Name != nil {
		c.Name = *update.Name
	}
	if update.Content != nil {
		c.Content = *update.Content
	}
	c.DateUpdated = now

	const q = `UPDATE cards SET
		"name" = $2,
		"content" = $3,
		"date_updated" = $4
		WHERE card_id = $1`
	_, err = db.ExecContext(ctx, q, id,
		c.Name, c.Content, c.DateUpdated,
	)
	if err != nil {
		return errors.Wrap(err, "updating card")
	}

	return nil
}

// Delete removes the card identified by a given ID.
func Delete(ctx context.Context, db *sqlx.DB, id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return ErrInvalidID
	}

	const q = `DELETE FROM cards WHERE card_id = $1`

	if _, err := db.ExecContext(ctx, q, id); err != nil {
		return errors.Wrapf(err, "deleting card %s", id)
	}

	return nil
}
