package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/claym/budgeting-app/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) ListIncomeReceipts(ctx context.Context, householdID uuid.UUID, fromDate, toDate string) ([]model.IncomeReceipt, error) {
	query := `SELECT id, household_id, amount_cents, description, income_date::text, voided, created_at, updated_at
		 FROM income_receipts WHERE household_id = $1`
	args := []any{householdID}
	argNum := 2

	if fromDate != "" {
		query += fmt.Sprintf(" AND income_date >= $%d::date", argNum)
		args = append(args, fromDate)
		argNum++
	}
	if toDate != "" {
		query += fmt.Sprintf(" AND income_date <= $%d::date", argNum)
		args = append(args, toDate)
		argNum++
	}

	query += " ORDER BY income_date DESC, created_at DESC"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list income receipts: %w", err)
	}
	defer rows.Close()

	var receipts []model.IncomeReceipt
	for rows.Next() {
		var receipt model.IncomeReceipt
		if err := rows.Scan(&receipt.ID, &receipt.HouseholdID, &receipt.AmountCents, &receipt.Description, &receipt.IncomeDate, &receipt.Voided, &receipt.CreatedAt, &receipt.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan income receipt: %w", err)
		}
		receipts = append(receipts, receipt)
	}
	return receipts, rows.Err()
}

func (r *Repository) CreateIncomeReceipt(ctx context.Context, householdID uuid.UUID, amountCents int64, description, incomeDate string) (*model.IncomeReceipt, error) {
	var receipt model.IncomeReceipt
	err := r.pool.QueryRow(ctx,
		`INSERT INTO income_receipts (household_id, amount_cents, description, income_date)
		 VALUES ($1, $2, $3, $4::date)
		 RETURNING id, household_id, amount_cents, description, income_date::text, voided, created_at, updated_at`,
		householdID, amountCents, description, incomeDate,
	).Scan(&receipt.ID, &receipt.HouseholdID, &receipt.AmountCents, &receipt.Description, &receipt.IncomeDate, &receipt.Voided, &receipt.CreatedAt, &receipt.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("create income receipt: %w", err)
	}
	return &receipt, nil
}

func (r *Repository) VoidIncomeReceipt(ctx context.Context, householdID, id uuid.UUID) (*model.IncomeReceipt, error) {
	var receipt model.IncomeReceipt
	err := r.pool.QueryRow(ctx,
		`UPDATE income_receipts SET voided = TRUE, updated_at = NOW()
		 WHERE id = $1 AND household_id = $2 AND voided = FALSE
		 RETURNING id, household_id, amount_cents, description, income_date::text, voided, created_at, updated_at`,
		id, householdID,
	).Scan(&receipt.ID, &receipt.HouseholdID, &receipt.AmountCents, &receipt.Description, &receipt.IncomeDate, &receipt.Voided, &receipt.CreatedAt, &receipt.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("void income receipt: %w", err)
	}
	return &receipt, nil
}

func (r *Repository) totalIncomeReceiptCents(ctx context.Context, householdID uuid.UUID) (int64, error) {
	var total int64
	err := r.pool.QueryRow(ctx,
		`SELECT COALESCE(SUM(amount_cents), 0) FROM income_receipts WHERE household_id = $1 AND NOT voided`,
		householdID,
	).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("sum income receipts: %w", err)
	}
	return total, nil
}
