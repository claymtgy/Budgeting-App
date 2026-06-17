package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func firstOfMonth(t time.Time) time.Time {
	utc := t.UTC()
	return time.Date(utc.Year(), utc.Month(), 1, 0, 0, 0, 0, time.UTC)
}

func currentMonthBounds(now time.Time) (start, end string) {
	startTime := firstOfMonth(now)
	endTime := startTime.AddDate(0, 1, 0)
	return startTime.Format("2006-01-02"), endTime.Format("2006-01-02")
}

func (r *Repository) EnsureMonthlyFunding(ctx context.Context, householdID uuid.UUID) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := r.ensureMonthlyFundingTx(ctx, tx, householdID); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *Repository) ensureMonthlyFundingTx(ctx context.Context, tx pgx.Tx, householdID uuid.UUID) error {
	var lastFunded *time.Time
	err := tx.QueryRow(ctx,
		`SELECT last_budget_month FROM households WHERE id = $1 FOR UPDATE`,
		householdID,
	).Scan(&lastFunded)
	if err != nil {
		return fmt.Errorf("lock household budget month: %w", err)
	}

	currentMonth := firstOfMonth(time.Now())

	if lastFunded == nil {
		_, err = tx.Exec(ctx,
			`UPDATE households SET last_budget_month = $1 WHERE id = $2`,
			currentMonth, householdID,
		)
		if err != nil {
			return fmt.Errorf("init budget month: %w", err)
		}
		return nil
	}

	last := *lastFunded
	for last.Before(currentMonth) {
		next := last.AddDate(0, 1, 0)
		_, err = tx.Exec(ctx,
			`UPDATE envelopes
			 SET balance_cents = balance_cents + allocated_cents, updated_at = NOW()
			 WHERE household_id = $1`,
			householdID,
		)
		if err != nil {
			return fmt.Errorf("fund envelopes: %w", err)
		}

		_, err = tx.Exec(ctx,
			`UPDATE households SET last_budget_month = $1 WHERE id = $2`,
			next, householdID,
		)
		if err != nil {
			return fmt.Errorf("advance budget month: %w", err)
		}
		last = next
	}

	return nil
}
