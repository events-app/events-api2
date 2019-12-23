package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/events-app/events-api2/internal/card"
	"github.com/events-app/events-api2/internal/platform/auth"
	"github.com/events-app/events-api2/internal/platform/web"
	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// Cards defines all of the handlers related to Cards. It holds the
// application state needed by the handler methods.
type Cards struct {
	db  *sqlx.DB
	log *log.Logger
}

// List gets all Cards from the service layer.
func (c *Cards) List(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	list, err := card.List(ctx, c.db)
	if err != nil {
		return errors.Wrap(err, "getting card list")
	}

	return web.Respond(ctx, w, list, http.StatusOK)
}

// Create decodes the body of a request to create a new card. The full
// card with generated fields is sent back in the response.
func (c *Cards) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return errors.New("claims missing from context")
	}

	var nc card.NewCard
	if err := web.Decode(r, &nc); err != nil {
		return errors.Wrap(err, "decoding new card")
	}

	car, err := card.Create(ctx, c.db, claims, nc, time.Now())
	if err != nil {
		return errors.Wrap(err, "creating new card")
	}

	return web.Respond(ctx, w, &car, http.StatusCreated)
}

// Retrieve finds a single card identified by an ID in the request URL.
func (c *Cards) Retrieve(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	car, err := card.Retrieve(ctx, c.db, id)
	if err != nil {
		switch err {
		case card.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case card.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "getting card %q", id)
		}
	}

	return web.Respond(ctx, w, car, http.StatusOK)
}

// Update decodes the body of a request to update an existing card. The ID
// of the card is part of the request URL.
func (c *Cards) Update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	var update card.UpdateCard
	if err := web.Decode(r, &update); err != nil {
		return errors.Wrap(err, "decoding card update")
	}

	// claims, ok := ctx.Value(auth.Key).(auth.Claims)
	// if !ok {
	// 	return errors.New("claims missing from context")
	// }

	if err := card.Update(ctx, c.db, id, update, time.Now()); err != nil {
		switch err {
		case card.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case card.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "updating card %q", id)
		}
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

// Delete removes a single card identified by an ID in the request URL.
func (c *Cards) Delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	if err := card.Delete(ctx, c.db, id); err != nil {
		switch err {
		case card.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "deleting card %q", id)
		}
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}
