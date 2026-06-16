package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/claym/budgeting-app/internal/repository"
	"github.com/google/uuid"
)

const HouseholdIDKey contextKey = "householdID"

func Household(repo *repository.Repository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := UserIDFromContext(r.Context())
			if !ok {
				http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
				return
			}

			householdID, err := repo.GetUserHouseholdID(r.Context(), userID)
			if errors.Is(err, repository.ErrNotFound) {
				http.Error(w, `{"error":"household not found"}`, http.StatusForbidden)
				return
			}
			if err != nil {
				http.Error(w, `{"error":"could not load household"}`, http.StatusInternalServerError)
				return
			}

			ctx := context.WithValue(r.Context(), HouseholdIDKey, householdID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func HouseholdIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(HouseholdIDKey).(uuid.UUID)
	return id, ok
}
