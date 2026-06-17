package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/claym/budgeting-app/internal/model"
	"github.com/google/uuid"
)

func (r *Repository) GetMonthlyTrends(ctx context.Context, householdID uuid.UUID, monthCount int) ([]model.MonthlyTrendPoint, error) {
	if monthCount < 1 {
		monthCount = 12
	}
	if monthCount > 24 {
		monthCount = 24
	}

	now := time.Now().UTC()
	current := firstOfMonth(now)
	start := current.AddDate(0, -(monthCount - 1), 0)
	startDate := start.Format("2006-01-02")
	endDate := current.AddDate(0, 1, 0).Format("2006-01-02")

	spentByMonth, err := r.sumExpensesByMonth(ctx, householdID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	incomeByMonth, err := r.sumIncomeReceiptsByMonth(ctx, householdID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	points := make([]model.MonthlyTrendPoint, 0, monthCount)
	for i := 0; i < monthCount; i++ {
		month := start.AddDate(0, i, 0).Format("2006-01")
		points = append(points, model.MonthlyTrendPoint{
			Month:       month,
			IncomeCents: incomeByMonth[month],
			SpentCents:  spentByMonth[month],
		})
	}

	return points, nil
}

func (r *Repository) sumExpensesByMonth(ctx context.Context, householdID uuid.UUID, startDate, endDate string) (map[string]int64, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT to_char(expense_date, 'YYYY-MM') AS month, COALESCE(SUM(amount_cents), 0)
		 FROM expenses
		 WHERE household_id = $1 AND NOT voided
		   AND expense_date >= $2::date AND expense_date < $3::date
		 GROUP BY 1`,
		householdID, startDate, endDate,
	)
	if err != nil {
		return nil, fmt.Errorf("sum expenses by month: %w", err)
	}
	defer rows.Close()

	result := make(map[string]int64)
	for rows.Next() {
		var month string
		var total int64
		if err := rows.Scan(&month, &total); err != nil {
			return nil, fmt.Errorf("scan expense month: %w", err)
		}
		result[month] = total
	}
	return result, rows.Err()
}

func (r *Repository) sumIncomeReceiptsByMonth(ctx context.Context, householdID uuid.UUID, startDate, endDate string) (map[string]int64, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT to_char(income_date, 'YYYY-MM') AS month, COALESCE(SUM(amount_cents), 0)
		 FROM income_receipts
		 WHERE household_id = $1 AND NOT voided
		   AND income_date >= $2::date AND income_date < $3::date
		 GROUP BY 1`,
		householdID, startDate, endDate,
	)
	if err != nil {
		return nil, fmt.Errorf("sum income receipts by month: %w", err)
	}
	defer rows.Close()

	result := make(map[string]int64)
	for rows.Next() {
		var month string
		var total int64
		if err := rows.Scan(&month, &total); err != nil {
			return nil, fmt.Errorf("scan income month: %w", err)
		}
		result[month] = total
	}
	return result, rows.Err()
}
