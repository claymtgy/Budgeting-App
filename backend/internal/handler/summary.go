package handler

import (
	"net/http"

	"github.com/claym/budgeting-app/internal/model"
	"github.com/claym/budgeting-app/internal/middleware"
	"github.com/claym/budgeting-app/internal/repository"
)

type SummaryHandler struct {
	repo *repository.Repository
}

func NewSummaryHandler(repo *repository.Repository) *SummaryHandler {
	return &SummaryHandler{repo: repo}
}

func (h *SummaryHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	summary, err := h.repo.GetSummary(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not load summary")
		return
	}
	if summary.Envelopes == nil {
		summary.Envelopes = []model.EnvelopeSummary{}
	}
	writeJSON(w, http.StatusOK, summary)
}
