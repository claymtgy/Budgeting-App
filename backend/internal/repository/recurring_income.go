package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/claym/budgeting-app/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func monthlyIncomeAmount(income model.Income) int64 {
	switch income.Period {
	case "weekly":
		return income.AmountCents * 52 / 12
	case "biweekly":
		return income.AmountCents * 26 / 12
	default:
		return income.AmountCents
	}
}

func (r *Repository) listIncomesTx(ctx context.Context, tx pgx.Tx, householdID uuid.UUID) ([]model.Income, error) {
	rows, err := tx.Query(ctx,
		`SELECT id, household_id, name, amount_cents, period, created_at, updated_at
		 FROM incomes WHERE household_id = $1 ORDER BY created_at DESC`,
		householdID,
	)
	if err != nil {
		return nil, fmt.Errorf("list incomes: %w", err)
	}
	defer rows.Close()

	var incomes []model.Income
	for rows.Next() {
		var income model.Income
		if err := rows.Scan(&income.ID, &income.HouseholdID, &income.Name, &income.AmountCents, &income.Period, &income.CreatedAt, &income.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan income: %w", err)
		}
		incomes = append(incomes, income)
	}
	return incomes, rows.Err()
}

func (r *Repository) ensureRecurringIncomeReceiptsForMonthTx(ctx context.Context, tx pgx.Tx, householdID uuid.UUID, month time.Time) error {
	incomes, err := r.listIncomesTx(ctx, tx, householdID)
	if err != nil {
		return err
	}
	if len(incomes) == 0 {
		return nil
	}

	monthStart := firstOfMonth(month).Format("2006-01-02")

	for _, income := range incomes {
		amount := monthlyIncomeAmount(income)
		if amount <= 0 {
			continue
		}

		var exists bool
		err := tx.QueryRow(ctx,
			`SELECT EXISTS (
				SELECT 1 FROM income_receipts
				WHERE household_id = $1 AND income_id = $2 AND auto_generated
				  AND date_trunc('month', income_date) = date_trunc('month', $3::date)
			)`,
			householdID, income.ID, monthStart,
		).Scan(&exists)
		if err != nil {
			return fmt.Errorf("check recurring income receipt: %w", err)
		}
		if exists {
			continue
		}

		_, err = tx.Exec(ctx,
			`INSERT INTO income_receipts (household_id, income_id, amount_cents, description, income_date, auto_generated)
			 VALUES ($1, $2, $3, $4, $5::date, TRUE)`,
			householdID, income.ID, amount, income.Name, monthStart,
		)
		if err != nil {
			return fmt.Errorf("insert recurring income receipt: %w", err)
		}
	}

	return nil
}

func (r *Repository) EnsureRecurringIncomeReceipts(ctx context.Context, householdID uuid.UUID) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	currentMonth := firstOfMonth(time.Now())
	if err := r.ensureRecurringIncomeReceiptsForMonthTx(ctx, tx, householdID, currentMonth); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
