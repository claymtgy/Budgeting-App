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

type EnvelopeHandler struct {
	repo *repository.Repository
}

func NewEnvelopeHandler(repo *repository.Repository) *EnvelopeHandler {
	return &EnvelopeHandler{repo: repo}
}

type envelopeRequest struct {
	Name           string `json:"name"`
	AllocatedCents int64  `json:"allocated_cents"`
}

func (h *EnvelopeHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	envelopes, err := h.repo.ListEnvelopes(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not list envelopes")
		return
	}
	if envelopes == nil {
		envelopes = []model.Envelope{}
	}
	writeJSON(w, http.StatusOK, envelopes)
}

func (h *EnvelopeHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req envelopeRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Name == "" || req.AllocatedCents < 0 {
		writeError(w, http.StatusBadRequest, "name required and allocation must be non-negative")
		return
	}

	envelope, err := h.repo.CreateEnvelope(r.Context(), userID, req.Name, req.AllocatedCents)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not create envelope")
		return
	}
	writeJSON(w, http.StatusCreated, envelope)
}

func (h *EnvelopeHandler) Update(w http.ResponseWriter, r *http.Request) {
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

	var req envelopeRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Name == "" || req.AllocatedCents < 0 {
		writeError(w, http.StatusBadRequest, "name required and allocation must be non-negative")
		return
	}

	envelope, err := h.repo.UpdateEnvelope(r.Context(), userID, id, req.Name, req.AllocatedCents)
	if errors.Is(err, repository.ErrNotFound) {
		writeError(w, http.StatusNotFound, "envelope not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not update envelope")
		return
	}
	writeJSON(w, http.StatusOK, envelope)
}

func (h *EnvelopeHandler) Delete(w http.ResponseWriter, r *http.Request) {
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

	if err := h.repo.DeleteEnvelope(r.Context(), userID, id); errors.Is(err, repository.ErrNotFound) {
		writeError(w, http.StatusNotFound, "envelope not found")
		return
	} else if err != nil {
		writeError(w, http.StatusInternalServerError, "could not delete envelope")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
