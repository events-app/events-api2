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
	const q = `SELECT
			p.*,
			COALESCE(SUM(s.quantity) ,0) AS sold,
			COALESCE(SUM(s.paid), 0) AS revenue
		FROM cards AS c
		LEFT JOIN sales AS s ON c.card_id = c.card_id
		GROUP BY p.card_id`

	if err := db.SelectContext(ctx, &cards, q); err != nil {
		return nil, errors.Wrap(err, "selecting cards")
	}

	return cards, nil
}

// Create adds a Card to the database. It returns the created Card with
// fields like ID and DateCreated populated..
func Create(ctx context.Context, db *sqlx.DB, user auth.Claims, np NewCard, now time.Time) (*Card, error) {
	p := Card{
		ID:          uuid.New().String(),
		Name:        np.Name,
		Cost:        np.Cost,
		Quantity:    np.Quantity,
		UserID:      user.Subject,
		DateCreated: now.UTC(),
		DateUpdated: now.UTC(),
	}

	const q = `
		INSERT INTO cards
		(card_id, user_id, name, cost, quantity, date_created, date_updated)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := db.ExecContext(ctx, q,
		p.ID, p.UserID,
		p.Name, p.Cost, p.Quantity,
		p.DateCreated, p.DateUpdated)
	if err != nil {
		return nil, errors.Wrap(err, "inserting card")
	}

	return &p, nil
}

// Retrieve finds the card identified by a given ID.
func Retrieve(ctx context.Context, db *sqlx.DB, id string) (*Card, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, ErrInvalidID
	}

	var p Card

	const q = `SELECT
			p.*,
			COALESCE(SUM(s.quantity), 0) AS sold,
			COALESCE(SUM(s.paid), 0) AS revenue
		FROM cards AS p
		LEFT JOIN sales AS s ON p.card_id = s.card_id
		WHERE p.card_id = $1
		GROUP BY p.card_id`

	if err := db.GetContext(ctx, &p, q, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}

		return nil, errors.Wrap(err, "selecting single card")
	}

	return &p, nil
}

// Update modifies data about a Card. It will error if the specified ID is
// invalid or does not reference an existing Card.
func Update(ctx context.Context, db *sqlx.DB, user auth.Claims, id string, update UpdateCard, now time.Time) error {
	p, err := Retrieve(ctx, db, id)
	if err != nil {
		return err
	}

	// If you do not have the admin role ...
	// and you are not the owner of this card ...
	// then get outta here!
	if !user.HasRole(auth.RoleAdmin) && p.UserID != user.Subject {
		return ErrForbidden
	}

	if update.Name != nil {
		p.Name = *update.Name
	}
	if update.Cost != nil {
		p.Cost = *update.Cost
	}
	if update.Quantity != nil {
		p.Quantity = *update.Quantity
	}
	p.DateUpdated = now

	const q = `UPDATE cards SET
		"name" = $2,
		"cost" = $3,
		"quantity" = $4,
		"date_updated" = $5
		WHERE card_id = $1`
	_, err = db.ExecContext(ctx, q, id,
		p.Name, p.Cost,
		p.Quantity, p.DateUpdated,
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
