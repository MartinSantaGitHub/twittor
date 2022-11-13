package middlewares

import (
	"context"
	"helpers"
	"net/http"
)

/* ValidateQueryId Validates that the user sends a valid query param id */
func ValidateQueryId(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")

		if len(id) < 1 {
			http.Error(w, "The id param is required", http.StatusBadRequest)

			return
		}

		ctx := context.WithValue(r.Context(), helpers.RequestQueryIdKey{}, id)
		r = r.Clone(ctx)

		next.ServeHTTP(w, r)
	}
}
