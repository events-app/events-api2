package menu

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
	// ErrNotFound is used when a specific Menu is requested but does not exist.
	ErrNotFound = errors.New("menu not found")

	// ErrInvalidID is used when an invalid UUID is provided.
	ErrInvalidID = errors.New("ID is not in its proper form")

	// ErrForbidden occurs when a user tries to do something that is forbidden to
	// them according to our access control policies.
	ErrForbidden = errors.New("Attempted action is not allowed")
)

// List gets all Menus from the database.
func List(ctx context.Context, db *sqlx.DB) ([]Menu, error) {
	menus := []Menu{}
	const q = `SELECT * FROM menus`
	if err := db.SelectContext(ctx, &menus, q); err != nil {
		return nil, errors.Wrap(err, "selecting menus")
	}

	return menus, nil
}

// Create adds a Menu to the database. It returns the created Menu with
// fields like ID and DateCreated populated..
func Create(ctx context.Context, db *sqlx.DB, user auth.Claims, nm NewMenu, now time.Time) (*Menu, error) {
	m := Menu{
		ID:          uuid.New().String(),
		Name:        nm.Name,
		CardID:     nm.CardID,
		UserID:      user.Subject,
		DateCreated: now.UTC(),
		DateUpdated: now.UTC(),
	}

	const q = `
		INSERT INTO menus
		(menu_id, user_id, name, card_id, date_created, date_updated)
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := db.ExecContext(ctx, q,
		m.ID, m.UserID,
		m.Name, m.CardID,
		m.DateCreated, m.DateUpdated)
	if err != nil {
		return nil, errors.Wrap(err, "inserting menu")
	}

	return &m, nil
}

// Retrieve finds the menu identified by a given ID.
func Retrieve(ctx context.Context, db *sqlx.DB, id string) (*Menu, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, ErrInvalidID
	}

	var m Menu

	// const c = `SELECT
	// 		p.*,
	// 		COALESCE(SUM(s.quantity), 0) AS sold,
	// 		COALESCE(SUM(s.paid), 0) AS revenue
	// 	FROM menus AS p
	// 	LEFT JOIN sales AS s ON p.menu_id = s.menu_id
	// 	WHERE p.menu_id = $1
	// 	GROUP BY p.menu_id`

	const q = `SELECT * FROM menus WHERE menu_id = $1`
	if err := db.GetContext(ctx, &m, q, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}

		return nil, errors.Wrap(err, "selecting single menu")
	}

	return &m, nil
}

// Update modifies data about a Menu. It will error if the specified ID is
// invalid or does not reference an existing Menu.
func Update(ctx context.Context, db *sqlx.DB, user auth.Claims, id string, update UpdateMenu, now time.Time) error {
	m, err := Retrieve(ctx, db, id)
	if err != nil {
		return err
	}

	// If you do not have the admin role ...
	// and you are not the owner of this menu ...
	// then get outta here!
	if !user.HasRole(auth.RoleAdmin) && m.UserID != user.Subject {
		return ErrForbidden
	}

	if update.Name != nil {
		m.Name = *update.Name
	}

	if update.CardID != nil {
		m.CardID = *update.CardID
	}

	m.DateUpdated = now

	const q = `UPDATE menus SET
		"name" = $2,
		"card_id" = $3,
		"date_updated" = $4
		WHERE menu_id = $1`
	_, err = db.ExecContext(ctx, q, id,
		m.Name, m.CardID, m.DateUpdated,
	)
	if err != nil {
		return errors.Wrap(err, "updating menu")
	}

	return nil
}

// Delete removes the menu identified by a given ID.
func Delete(ctx context.Context, db *sqlx.DB, id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return ErrInvalidID
	}

	const q = `DELETE FROM menus WHERE menu_id = $1`

	if _, err := db.ExecContext(ctx, q, id); err != nil {
		return errors.Wrapf(err, "deleting menu %s", id)
	}

	return nil
}
