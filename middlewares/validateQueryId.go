package middlewares

import (
	"context"
	"helpers"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

/* ValidateQueryId Validates that the user sends a valid query param id */
func ValidateQueryId(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")

		if len(id) < 1 {
			http.Error(w, "The id param is required", http.StatusBadRequest)

			return
		}

		objId, err := primitive.ObjectIDFromHex(id)

		if err != nil {
			http.Error(w, "Invalid id param", http.StatusBadRequest)

			return
		}

		ctx := context.WithValue(r.Context(), helpers.RequestQueryIdKey{}, objId)
		r = r.Clone(ctx)

		next.ServeHTTP(w, r)
	}
}
