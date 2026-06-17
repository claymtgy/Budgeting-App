package handler

import (
	"fmt"
	"net/http"
	"strings"
	"time"

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
	householdID, ok := middleware.HouseholdIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	summary, err := h.repo.GetSummary(r.Context(), householdID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not load summary")
		return
	}
	if summary.Envelopes == nil {
		summary.Envelopes = []model.EnvelopeSummary{}
	}
	writeJSON(w, http.StatusOK, summary)
}

func (h *SummaryHandler) GetMonthly(w http.ResponseWriter, r *http.Request) {
	householdID, ok := middleware.HouseholdIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	month := strings.TrimSpace(r.URL.Query().Get("month"))
	if month == "" {
		writeError(w, http.StatusBadRequest, "month required (YYYY-MM)")
		return
	}
	if _, err := time.Parse("2006-01", month); err != nil {
		writeError(w, http.StatusBadRequest, "invalid month format")
		return
	}

	currentMonth := time.Now().UTC().Format("2006-01")
	if month > currentMonth {
		writeError(w, http.StatusBadRequest, "month cannot be in the future")
		return
	}

	summary, err := h.repo.GetMonthlySummary(r.Context(), householdID, month)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not load monthly summary")
		return
	}
	if summary.Envelopes == nil {
		summary.Envelopes = []model.EnvelopeSummary{}
	}
	writeJSON(w, http.StatusOK, summary)
}

func (h *SummaryHandler) GetTrends(w http.ResponseWriter, r *http.Request) {
	householdID, ok := middleware.HouseholdIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	months := 12
	if raw := strings.TrimSpace(r.URL.Query().Get("months")); raw != "" {
		var parsed int
		if _, err := fmt.Sscanf(raw, "%d", &parsed); err != nil || parsed < 1 {
			writeError(w, http.StatusBadRequest, "invalid months parameter")
			return
		}
		months = parsed
	}

	trends, err := h.repo.GetMonthlyTrends(r.Context(), householdID, months)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not load monthly trends")
		return
	}
	if trends == nil {
		trends = []model.MonthlyTrendPoint{}
	}
	writeJSON(w, http.StatusOK, trends)
}
