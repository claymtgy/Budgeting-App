package handler

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/claym/budgeting-app/internal/model"
	"github.com/claym/budgeting-app/internal/middleware"
	"github.com/claym/budgeting-app/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type IncomeReceiptHandler struct {
	repo *repository.Repository
}

func NewIncomeReceiptHandler(repo *repository.Repository) *IncomeReceiptHandler {
	return &IncomeReceiptHandler{repo: repo}
}

type incomeReceiptRequest struct {
	AmountCents int64  `json:"amount_cents"`
	Description string `json:"description"`
	IncomeDate  string `json:"income_date"`
}

func (h *IncomeReceiptHandler) List(w http.ResponseWriter, r *http.Request) {
	householdID, ok := middleware.HouseholdIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	fromDate := strings.TrimSpace(r.URL.Query().Get("from"))
	toDate := strings.TrimSpace(r.URL.Query().Get("to"))

	if fromDate != "" {
		if _, err := time.Parse("2006-01-02", fromDate); err != nil {
			writeError(w, http.StatusBadRequest, "invalid from date")
			return
		}
	}
	if toDate != "" {
		if _, err := time.Parse("2006-01-02", toDate); err != nil {
			writeError(w, http.StatusBadRequest, "invalid to date")
			return
		}
	}
	if fromDate != "" && toDate != "" && fromDate > toDate {
		writeError(w, http.StatusBadRequest, "from date must be on or before to date")
		return
	}

	receipts, err := h.repo.ListIncomeReceipts(r.Context(), householdID, fromDate, toDate)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not list income receipts")
		return
	}
	if receipts == nil {
		receipts = []model.IncomeReceipt{}
	}
	writeJSON(w, http.StatusOK, receipts)
}

func (h *IncomeReceiptHandler) Create(w http.ResponseWriter, r *http.Request) {
	householdID, ok := middleware.HouseholdIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req incomeReceiptRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.AmountCents <= 0 {
		writeError(w, http.StatusBadRequest, "amount must be positive")
		return
	}
	if req.IncomeDate == "" {
		writeError(w, http.StatusBadRequest, "income_date required")
		return
	}

	receipt, err := h.repo.CreateIncomeReceipt(r.Context(), householdID, req.AmountCents, strings.TrimSpace(req.Description), req.IncomeDate)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not create income receipt")
		return
	}
	writeJSON(w, http.StatusCreated, receipt)
}

func (h *IncomeReceiptHandler) Void(w http.ResponseWriter, r *http.Request) {
	householdID, ok := middleware.HouseholdIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	receipt, err := h.repo.VoidIncomeReceipt(r.Context(), householdID, id)
	if errors.Is(err, repository.ErrNotFound) {
		writeError(w, http.StatusNotFound, "income receipt not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not void income receipt")
		return
	}
	writeJSON(w, http.StatusOK, receipt)
}
