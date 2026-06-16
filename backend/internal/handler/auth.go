package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/claym/budgeting-app/internal/auth"
	"github.com/claym/budgeting-app/internal/middleware"
	"github.com/claym/budgeting-app/internal/repository"
)

type AuthHandler struct {
	repo   *repository.Repository
	tokens *auth.TokenService
}

func NewAuthHandler(repo *repository.Repository, tokens *auth.TokenService) *AuthHandler {
	return &AuthHandler{repo: repo, tokens: tokens}
}

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	JoinCode string `json:"join_code"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authResponse struct {
	Token string      `json:"token"`
	User  userPayload `json:"user"`
}

type userPayload struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	JoinCode string `json:"join_code"`
}

func userPayloadFrom(userID, email, joinCode string) userPayload {
	return userPayload{ID: userID, Email: email, JoinCode: joinCode}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	email := strings.TrimSpace(strings.ToLower(req.Email))
	if email == "" || len(req.Password) < 8 {
		writeError(w, http.StatusBadRequest, "email required and password must be at least 8 characters")
		return
	}

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not create user")
		return
	}

	user, household, err := h.repo.RegisterUser(r.Context(), email, hash, req.JoinCode)
	if errors.Is(err, repository.ErrInvalidJoinCode) {
		writeError(w, http.StatusBadRequest, "invalid join code")
		return
	}
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
			writeError(w, http.StatusConflict, "email already registered")
			return
		}
		writeError(w, http.StatusInternalServerError, "could not create user")
		return
	}

	token, err := h.tokens.GenerateToken(user.ID, user.Email)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not create session")
		return
	}

	writeJSON(w, http.StatusCreated, authResponse{
		Token: token,
		User:  userPayloadFrom(user.ID.String(), user.Email, household.JoinCode),
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	email := strings.TrimSpace(strings.ToLower(req.Email))
	user, err := h.repo.GetUserByEmail(r.Context(), email)
	if errors.Is(err, repository.ErrNotFound) {
		writeError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not log in")
		return
	}

	if err := auth.CheckPassword(user.PasswordHash, req.Password); err != nil {
		writeError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}

	household, err := h.repo.GetHouseholdByUserID(r.Context(), user.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not log in")
		return
	}

	token, err := h.tokens.GenerateToken(user.ID, user.Email)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not create session")
		return
	}

	writeJSON(w, http.StatusOK, authResponse{
		Token: token,
		User:  userPayloadFrom(user.ID.String(), user.Email, household.JoinCode),
	})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	user, err := h.repo.GetUserByID(r.Context(), userID)
	if errors.Is(err, repository.ErrNotFound) {
		writeError(w, http.StatusNotFound, "user not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not load user")
		return
	}

	household, err := h.repo.GetHouseholdByUserID(r.Context(), user.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not load household")
		return
	}

	writeJSON(w, http.StatusOK, userPayloadFrom(user.ID.String(), user.Email, household.JoinCode))
}
