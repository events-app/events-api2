package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/events-app/events-api2/internal/platform/auth"
	"github.com/events-app/events-api2/internal/platform/web"
	"github.com/events-app/events-api2/internal/menu"
	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// Menus defines all of the handlers related to Menus. It holds the
// application state needed by the handler methods.
type Menus struct {
	db  *sqlx.DB
	log *log.Logger
}

// List gets all Menus from the service layer.
func (m *Menus) List(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	list, err := menu.List(ctx, m.db)
	if err != nil {
		return errors.Wrap(err, "getting menu list")
	}

	return web.Respond(ctx, w, list, http.StatusOK)
}

// Create decodes the body of a request to create a new menu. The full
// menu with generated fields is sent back in the response.
func (m *Menus) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return errors.New("claims missing from context")
	}

	var nm menu.NewMenu
	if err := web.Decode(r, &nm); err != nil {
		return errors.Wrap(err, "decoding new menu")
	}

	men, err := menu.Create(ctx, m.db, claims, nm, time.Now())
	if err != nil {
		return errors.Wrap(err, "creating new menu")
	}

	return web.Respond(ctx, w, &men, http.StatusCreated)
}

// Retrieve finds a single menu identified by an ID in the request URL.
func (m *Menus) Retrieve(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	men, err := menu.Retrieve(ctx, m.db, id)
	if err != nil {
		switch err {
		case menu.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case menu.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "getting menu %q", id)
		}
	}

	return web.Respond(ctx, w, men, http.StatusOK)
}

// Update decodes the body of a request to update an existing menu. The ID
// of the menu is part of the request URL.
func (m *Menus) Update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	var update menu.UpdateMenu
	if err := web.Decode(r, &update); err != nil {
		return errors.Wrap(err, "decoding menu update")
	}

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return errors.New("claims missing from context")
	}

	if err := menu.Update(ctx, m.db, claims, id, update, time.Now()); err != nil {
		switch err {
		case menu.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case menu.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		case menu.ErrForbidden:
			return web.NewRequestError(err, http.StatusForbidden)
		default:
			return errors.Wrapf(err, "updating menu %q", id)
		}
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

// Delete removes a single menu identified by an ID in the request URL.
func (m *Menus) Delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	if err := menu.Delete(ctx, m.db, id); err != nil {
		switch err {
		case menu.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "deleting menu %q", id)
		}
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

