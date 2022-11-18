package middlewares

import (
	"context"
	"helpers"
	"net/http"
	"strconv"
)

/* ValidatePageLimit Validates that the user sends a valid page and limit query param */
func ValidatePageLimit(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pageQuery := r.URL.Query().Get("page")

		if len(pageQuery) < 1 {
			pageQuery = "1"
		}

		page, err := strconv.ParseInt(pageQuery, 10, 64)

		if err != nil || page <= 0 {
			http.Error(w, "The page param is invalid. It must be a number and greater to zero", http.StatusBadRequest)

			return
		}

		limitQuery := r.URL.Query().Get("limit")

		if len(limitQuery) < 1 {
			limitQuery = helpers.GetEnvVariable("RECORDS_LIMIT")
		}

		limit, err := strconv.ParseInt(limitQuery, 10, 64)

		if err != nil {
			http.Error(w, "The limit param is invalid", http.StatusBadRequest)

			return
		}

		ctx := context.WithValue(r.Context(), helpers.RequestPageKey{}, page)
		r = r.Clone(ctx)

		ctx = context.WithValue(r.Context(), helpers.RequestLimitKey{}, limit)
		r = r.Clone(ctx)

		next.ServeHTTP(w, r)
	}
}
