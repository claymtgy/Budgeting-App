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

type IncomeHandler struct {
	repo *repository.Repository
}

func NewIncomeHandler(repo *repository.Repository) *IncomeHandler {
	return &IncomeHandler{repo: repo}
}

type incomeRequest struct {
	Name        string `json:"name"`
	AmountCents int64  `json:"amount_cents"`
	Period      string `json:"period"`
}

func (h *IncomeHandler) List(w http.ResponseWriter, r *http.Request) {
	householdID, ok := middleware.HouseholdIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	incomes, err := h.repo.ListIncomes(r.Context(), householdID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not list incomes")
		return
	}
	if incomes == nil {
		incomes = []model.Income{}
	}
	writeJSON(w, http.StatusOK, incomes)
}

func (h *IncomeHandler) Create(w http.ResponseWriter, r *http.Request) {
	householdID, ok := middleware.HouseholdIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req incomeRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Name == "" || req.AmountCents < 0 {
		writeError(w, http.StatusBadRequest, "name required and amount must be non-negative")
		return
	}
	if req.Period == "" {
		req.Period = "monthly"
	}

	income, err := h.repo.CreateIncome(r.Context(), householdID, req.Name, req.AmountCents, req.Period)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not create income")
		return
	}
	writeJSON(w, http.StatusCreated, income)
}

func (h *IncomeHandler) Update(w http.ResponseWriter, r *http.Request) {
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

	var req incomeRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Name == "" || req.AmountCents < 0 {
		writeError(w, http.StatusBadRequest, "name required and amount must be non-negative")
		return
	}
	if req.Period == "" {
		req.Period = "monthly"
	}

	income, err := h.repo.UpdateIncome(r.Context(), householdID, id, req.Name, req.AmountCents, req.Period)
	if errors.Is(err, repository.ErrNotFound) {
		writeError(w, http.StatusNotFound, "income not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not update income")
		return
	}
	writeJSON(w, http.StatusOK, income)
}

func (h *IncomeHandler) Delete(w http.ResponseWriter, r *http.Request) {
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

	if err := h.repo.DeleteIncome(r.Context(), householdID, id); errors.Is(err, repository.ErrNotFound) {
		writeError(w, http.StatusNotFound, "income not found")
		return
	} else if err != nil {
		writeError(w, http.StatusInternalServerError, "could not delete income")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
