package handler

import (
	"errors"
	"net/http"

	"github.com/claym/budgeting-app/internal/model"
	"github.com/claym/budgeting-app/internal/middleware"
	"github.com/claym/budgeting-app/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type ExpenseHandler struct {
	repo *repository.Repository
}

func NewExpenseHandler(repo *repository.Repository) *ExpenseHandler {
	return &ExpenseHandler{repo: repo}
}

type expenseRequest struct {
	EnvelopeID  string `json:"envelope_id"`
	AmountCents int64  `json:"amount_cents"`
	Description string `json:"description"`
	ExpenseDate string `json:"expense_date"`
}

func (h *ExpenseHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	expenses, err := h.repo.ListExpenses(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not list expenses")
		return
	}
	if expenses == nil {
		expenses = []model.Expense{}
	}
	writeJSON(w, http.StatusOK, expenses)
}

func (h *ExpenseHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req expenseRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	envelopeID, err := uuid.Parse(req.EnvelopeID)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid envelope_id")
		return
	}
	if req.AmountCents <= 0 {
		writeError(w, http.StatusBadRequest, "amount must be positive")
		return
	}
	if req.ExpenseDate == "" {
		writeError(w, http.StatusBadRequest, "expense_date required")
		return
	}

	expense, err := h.repo.CreateExpense(r.Context(), userID, envelopeID, req.AmountCents, req.Description, req.ExpenseDate)
	if errors.Is(err, repository.ErrNotFound) {
		writeError(w, http.StatusNotFound, "envelope not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not create expense")
		return
	}
	writeJSON(w, http.StatusCreated, expense)
}

func (h *ExpenseHandler) Void(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	expense, err := h.repo.VoidExpense(r.Context(), userID, id)
	if errors.Is(err, repository.ErrNotFound) {
		writeError(w, http.StatusNotFound, "expense not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not void expense")
		return
	}
	writeJSON(w, http.StatusOK, expense)
}
