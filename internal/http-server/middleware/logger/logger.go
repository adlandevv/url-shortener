package logger

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func New(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log = log.With(
			slog.String("component", "middleware/logger"),
		)

		log.Info("logger middleware enabled")

		// handler code
		fn := func(w http.ResponseWriter, r *http.Request) {
			// collecting initial information about the request
			entry := log.With(
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
				slog.String("request_id", middleware.GetReqID(r.Context())),
			)

			// creating a wrapper around http.ResponseWriter
			// to expose the response
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			// Moment the request is received to calculate the processing time
			t1 := time.Now()

			// The entry will be written to the log in a defer statement
			// by that time, the request will have already been processed
			defer func() {
				entry.Info("request completed",
					slog.Int("status", ww.Status()),
					slog.Int("bytes", ww.BytesWritten()),
					slog.String("duration", time.Since(t1).String()),
				)
			}()

			// Pass control to the next handler in the middleware chain
			next.ServeHTTP(ww, r)
		}

		// Return handler created above, casting it to the `http.HandlerFunc` type
		return http.HandlerFunc(fn)
	}
}
