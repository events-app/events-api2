package mid

import (
	"context"
	"net/http"

	"github.com/events-app/events-api2/internal/platform/web"
)

func Cors(origin string) web.Middleware {

	// This is the actual middleware function to be executed.
	f := func(before web.Handler) web.Handler {

		// Create the handler that will be attached in the middleware chain.
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			// w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			// w.Header().Set("Access-Control-Allow-Headers",
			// 	`Accept, Content-Type, Content-Length,
			// 	 Accept-Encoding, X-CSRF-Token, Authorization`)
			return before(ctx, w, r)
		}

		return h
	}

	return f
}
