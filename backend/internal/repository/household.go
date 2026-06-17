package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/claym/budgeting-app/internal/household"
	"github.com/claym/budgeting-app/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var ErrInvalidJoinCode = errors.New("invalid join code")

func (r *Repository) RegisterUser(ctx context.Context, email, passwordHash, joinCode string) (*model.User, *model.Household, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	joinCode = strings.ToUpper(strings.TrimSpace(joinCode))

	var h model.Household
	if joinCode == "" {
		for range 5 {
			code, err := household.GenerateJoinCode()
			if err != nil {
				return nil, nil, err
			}
			err = tx.QueryRow(ctx,
				`INSERT INTO households (join_code, last_budget_month) VALUES ($1, date_trunc('month', CURRENT_DATE)::date)
				 RETURNING id, join_code, created_at`,
				code,
			).Scan(&h.ID, &h.JoinCode, &h.CreatedAt)
			if err == nil {
				break
			}
			if !isUniqueViolation(err) {
				return nil, nil, fmt.Errorf("create household: %w", err)
			}
		}
		if h.ID == uuid.Nil {
			return nil, nil, fmt.Errorf("create household: could not generate unique join code")
		}
	} else {
		err = tx.QueryRow(ctx,
			`SELECT id, join_code, created_at FROM households WHERE join_code = $1`,
			joinCode,
		).Scan(&h.ID, &h.JoinCode, &h.CreatedAt)
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil, ErrInvalidJoinCode
		}
		if err != nil {
			return nil, nil, fmt.Errorf("find household: %w", err)
		}
	}

	var user model.User
	err = tx.QueryRow(ctx,
		`INSERT INTO users (email, password_hash, household_id) VALUES ($1, $2, $3)
		 RETURNING id, email, password_hash, household_id, created_at`,
		email, passwordHash, h.ID,
	).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.HouseholdID, &user.CreatedAt)
	if err != nil {
		return nil, nil, fmt.Errorf("create user: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, nil, fmt.Errorf("commit transaction: %w", err)
	}

	return &user, &h, nil
}

func (r *Repository) GetHouseholdByUserID(ctx context.Context, userID uuid.UUID) (*model.Household, error) {
	var h model.Household
	err := r.pool.QueryRow(ctx,
		`SELECT h.id, h.join_code, h.created_at
		 FROM households h
		 JOIN users u ON u.household_id = h.id
		 WHERE u.id = $1`,
		userID,
	).Scan(&h.ID, &h.JoinCode, &h.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get household: %w", err)
	}
	return &h, nil
}

func (r *Repository) GetUserHouseholdID(ctx context.Context, userID uuid.UUID) (uuid.UUID, error) {
	var householdID uuid.UUID
	err := r.pool.QueryRow(ctx,
		`SELECT household_id FROM users WHERE id = $1`,
		userID,
	).Scan(&householdID)
	if errors.Is(err, pgx.ErrNoRows) {
		return uuid.Nil, ErrNotFound
	}
	if err != nil {
		return uuid.Nil, fmt.Errorf("get user household: %w", err)
	}
	return householdID, nil
}

func isUniqueViolation(err error) bool {
	return err != nil && strings.Contains(err.Error(), "unique")
}
