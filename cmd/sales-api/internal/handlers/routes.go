package handlers

import (
	"log"
	"net/http"

	"github.com/events-app/events-api2/internal/mid"
	"github.com/events-app/events-api2/internal/platform/auth"
	"github.com/events-app/events-api2/internal/platform/web"
	"github.com/jmoiron/sqlx"
)

// API constructs an http.Handler with all application routes defined.
func API(db *sqlx.DB, log *log.Logger, authenticator *auth.Authenticator) http.Handler {

	// Construct the web.App which holds all routes as well as common Middleware.
	app := web.NewApp(log, mid.Logger(log), mid.Errors(log), mid.Metrics())

	{
		// Register health check handler. This route is not authenticated.
		c := Check{db: db}
		app.Handle(http.MethodGet, "/v1/health", c.Health)
	}

	{
		// Register user handlers.
		u := Users{db: db, authenticator: authenticator}

		// The token route can't be authenticated because they need this route to
		// get the token in the first place.
		app.Handle(http.MethodGet, "/v1/users/token", u.Token)
	}

	{
		// Register Card handlers. Ensure all routes are authenticated.
		c := Cards{db: db, log: log}

		app.Handle(http.MethodGet, "/v1/cards", c.List, mid.Authenticate(authenticator))
		app.Handle(http.MethodGet, "/v1/cards/{id}", c.Retrieve, mid.Authenticate(authenticator))
		app.Handle(http.MethodPost, "/v1/cards", c.Create, mid.Authenticate(authenticator))
		app.Handle(http.MethodPut, "/v1/cards/{id}", c.Update, mid.Authenticate(authenticator))
		app.Handle(http.MethodDelete, "/v1/cards/{id}", c.Delete, mid.Authenticate(authenticator), mid.HasRole(auth.RoleAdmin))

		// app.Handle(http.MethodPost, "/v1/cards/{id}/sales", c.AddSale, mid.Authenticate(authenticator), mid.HasRole(auth.RoleAdmin))
		// app.Handle(http.MethodGet, "/v1/cards/{id}/sales", c.ListSales, mid.Authenticate(authenticator))
	}

	return app
}
