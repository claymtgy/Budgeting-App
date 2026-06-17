package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/claym/budgeting-app/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("not found")

type Repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) CreateUser(ctx context.Context, email, passwordHash string) (*model.User, error) {
	user, _, err := r.RegisterUser(ctx, email, passwordHash, "")
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.pool.QueryRow(ctx,
		`SELECT id, email, password_hash, household_id, created_at FROM users WHERE email = $1`,
		email,
	).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.HouseholdID, &user.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get user by email: %w", err)
	}
	return &user, nil
}

func (r *Repository) GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	var user model.User
	err := r.pool.QueryRow(ctx,
		`SELECT id, email, password_hash, household_id, created_at FROM users WHERE id = $1`,
		id,
	).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.HouseholdID, &user.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	return &user, nil
}

func (r *Repository) ListIncomes(ctx context.Context, householdID uuid.UUID) ([]model.Income, error) {
	rows, err := r.pool.Query(ctx,
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

func (r *Repository) CreateIncome(ctx context.Context, householdID uuid.UUID, name string, amountCents int64, period string) (*model.Income, error) {
	var income model.Income
	err := r.pool.QueryRow(ctx,
		`INSERT INTO incomes (household_id, name, amount_cents, period)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, household_id, name, amount_cents, period, created_at, updated_at`,
		householdID, name, amountCents, period,
	).Scan(&income.ID, &income.HouseholdID, &income.Name, &income.AmountCents, &income.Period, &income.CreatedAt, &income.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("create income: %w", err)
	}
	return &income, nil
}

func (r *Repository) UpdateIncome(ctx context.Context, householdID, id uuid.UUID, name string, amountCents int64, period string) (*model.Income, error) {
	var income model.Income
	err := r.pool.QueryRow(ctx,
		`UPDATE incomes SET name = $1, amount_cents = $2, period = $3, updated_at = NOW()
		 WHERE id = $4 AND household_id = $5
		 RETURNING id, household_id, name, amount_cents, period, created_at, updated_at`,
		name, amountCents, period, id, householdID,
	).Scan(&income.ID, &income.HouseholdID, &income.Name, &income.AmountCents, &income.Period, &income.CreatedAt, &income.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("update income: %w", err)
	}
	return &income, nil
}

func (r *Repository) DeleteIncome(ctx context.Context, householdID, id uuid.UUID) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM incomes WHERE id = $1 AND household_id = $2`, id, householdID)
	if err != nil {
		return fmt.Errorf("delete income: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) ListEnvelopes(ctx context.Context, householdID uuid.UUID) ([]model.Envelope, error) {
	if err := r.EnsureMonthlyFunding(ctx, householdID); err != nil {
		return nil, err
	}

	rows, err := r.pool.Query(ctx,
		`SELECT id, household_id, name, allocated_cents, balance_cents, created_at, updated_at
		 FROM envelopes WHERE household_id = $1 ORDER BY created_at DESC`,
		householdID,
	)
	if err != nil {
		return nil, fmt.Errorf("list envelopes: %w", err)
	}
	defer rows.Close()

	var envelopes []model.Envelope
	for rows.Next() {
		var envelope model.Envelope
		if err := rows.Scan(&envelope.ID, &envelope.HouseholdID, &envelope.Name, &envelope.AllocatedCents, &envelope.BalanceCents, &envelope.CreatedAt, &envelope.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan envelope: %w", err)
		}
		envelopes = append(envelopes, envelope)
	}
	return envelopes, rows.Err()
}

func (r *Repository) CreateEnvelope(ctx context.Context, householdID uuid.UUID, name string, allocatedCents int64) (*model.Envelope, error) {
	var envelope model.Envelope
	err := r.pool.QueryRow(ctx,
		`INSERT INTO envelopes (household_id, name, allocated_cents, balance_cents)
		 VALUES ($1, $2, $3, $3)
		 RETURNING id, household_id, name, allocated_cents, balance_cents, created_at, updated_at`,
		householdID, name, allocatedCents,
	).Scan(&envelope.ID, &envelope.HouseholdID, &envelope.Name, &envelope.AllocatedCents, &envelope.BalanceCents, &envelope.CreatedAt, &envelope.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("create envelope: %w", err)
	}
	return &envelope, nil
}

func (r *Repository) UpdateEnvelope(ctx context.Context, householdID, id uuid.UUID, name string, allocatedCents int64) (*model.Envelope, error) {
	var envelope model.Envelope
	err := r.pool.QueryRow(ctx,
		`UPDATE envelopes SET name = $1, allocated_cents = $2, updated_at = NOW()
		 WHERE id = $3 AND household_id = $4
		 RETURNING id, household_id, name, allocated_cents, balance_cents, created_at, updated_at`,
		name, allocatedCents, id, householdID,
	).Scan(&envelope.ID, &envelope.HouseholdID, &envelope.Name, &envelope.AllocatedCents, &envelope.BalanceCents, &envelope.CreatedAt, &envelope.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("update envelope: %w", err)
	}
	return &envelope, nil
}

func (r *Repository) DeleteEnvelope(ctx context.Context, householdID, id uuid.UUID) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM envelopes WHERE id = $1 AND household_id = $2`, id, householdID)
	if err != nil {
		return fmt.Errorf("delete envelope: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) ListExpenses(ctx context.Context, householdID uuid.UUID, fromDate, toDate string) ([]model.Expense, error) {
	query := `SELECT id, household_id, envelope_id, amount_cents, description, place, expense_date::text, voided, created_at, updated_at
		 FROM expenses WHERE household_id = $1`
	args := []any{householdID}
	argNum := 2

	if fromDate != "" {
		query += fmt.Sprintf(" AND expense_date >= $%d::date", argNum)
		args = append(args, fromDate)
		argNum++
	}
	if toDate != "" {
		query += fmt.Sprintf(" AND expense_date <= $%d::date", argNum)
		args = append(args, toDate)
		argNum++
	}

	query += " ORDER BY expense_date DESC, created_at DESC"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list expenses: %w", err)
	}
	defer rows.Close()

	var expenses []model.Expense
	for rows.Next() {
		var expense model.Expense
		if err := rows.Scan(&expense.ID, &expense.HouseholdID, &expense.EnvelopeID, &expense.AmountCents, &expense.Description, &expense.Place, &expense.ExpenseDate, &expense.Voided, &expense.CreatedAt, &expense.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan expense: %w", err)
		}
		expenses = append(expenses, expense)
	}
	return expenses, rows.Err()
}

func (r *Repository) ListExpensePlaces(ctx context.Context, householdID uuid.UUID) ([]string, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT DISTINCT place FROM expenses
		 WHERE household_id = $1 AND place <> ''
		 ORDER BY place ASC`,
		householdID,
	)
	if err != nil {
		return nil, fmt.Errorf("list expense places: %w", err)
	}
	defer rows.Close()

	var places []string
	for rows.Next() {
		var place string
		if err := rows.Scan(&place); err != nil {
			return nil, fmt.Errorf("scan place: %w", err)
		}
		places = append(places, place)
	}
	if places == nil {
		places = []string{}
	}
	return places, rows.Err()
}

func (r *Repository) CreateExpense(ctx context.Context, householdID, envelopeID uuid.UUID, amountCents int64, description, place, expenseDate string) (*model.Expense, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := r.ensureMonthlyFundingTx(ctx, tx, householdID); err != nil {
		return nil, err
	}

	var expense model.Expense
	err = tx.QueryRow(ctx,
		`INSERT INTO expenses (household_id, envelope_id, amount_cents, description, place, expense_date)
		 SELECT $1, $2, $3, $4, $5, $6::date
		 WHERE EXISTS (SELECT 1 FROM envelopes WHERE id = $2 AND household_id = $1)
		 RETURNING id, household_id, envelope_id, amount_cents, description, place, expense_date::text, voided, created_at, updated_at`,
		householdID, envelopeID, amountCents, description, place, expenseDate,
	).Scan(&expense.ID, &expense.HouseholdID, &expense.EnvelopeID, &expense.AmountCents, &expense.Description, &expense.Place, &expense.ExpenseDate, &expense.Voided, &expense.CreatedAt, &expense.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("create expense: %w", err)
	}

	tag, err := tx.Exec(ctx,
		`UPDATE envelopes SET balance_cents = balance_cents - $1, updated_at = NOW()
		 WHERE id = $2 AND household_id = $3`,
		amountCents, envelopeID, householdID,
	)
	if err != nil {
		return nil, fmt.Errorf("decrement envelope balance: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return nil, ErrNotFound
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit expense: %w", err)
	}
	return &expense, nil
}

func (r *Repository) VoidExpense(ctx context.Context, householdID, id uuid.UUID) (*model.Expense, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var expense model.Expense
	err = tx.QueryRow(ctx,
		`UPDATE expenses SET voided = TRUE, updated_at = NOW()
		 WHERE id = $1 AND household_id = $2 AND voided = FALSE
		 RETURNING id, household_id, envelope_id, amount_cents, description, place, expense_date::text, voided, created_at, updated_at`,
		id, householdID,
	).Scan(&expense.ID, &expense.HouseholdID, &expense.EnvelopeID, &expense.AmountCents, &expense.Description, &expense.Place, &expense.ExpenseDate, &expense.Voided, &expense.CreatedAt, &expense.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("void expense: %w", err)
	}

	_, err = tx.Exec(ctx,
		`UPDATE envelopes SET balance_cents = balance_cents + $1, updated_at = NOW()
		 WHERE id = $2 AND household_id = $3`,
		expense.AmountCents, expense.EnvelopeID, householdID,
	)
	if err != nil {
		return nil, fmt.Errorf("restore envelope balance: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit void expense: %w", err)
	}
	return &expense, nil
}

func (r *Repository) GetSummary(ctx context.Context, householdID uuid.UUID) (*model.BudgetSummary, error) {
	if err := r.EnsureMonthlyFunding(ctx, householdID); err != nil {
		return nil, err
	}

	incomes, err := r.ListIncomes(ctx, householdID)
	if err != nil {
		return nil, err
	}

	envelopes, err := r.listEnvelopesWithoutFunding(ctx, householdID)
	if err != nil {
		return nil, err
	}

	expenses, err := r.ListExpenses(ctx, householdID, "", "")
	if err != nil {
		return nil, err
	}

	monthStart, monthEnd := currentMonthBounds(time.Now())

	incomeAmounts := make([]int64, len(incomes))
	for i, income := range incomes {
		incomeAmounts[i] = income.AmountCents
	}

	allocatedAmounts := make([]int64, len(envelopes))
	for i, envelope := range envelopes {
		allocatedAmounts[i] = envelope.AllocatedCents
	}

	spentThisMonth := make(map[uuid.UUID]int64)
	var totalSpent int64
	for _, expense := range expenses {
		if expense.Voided {
			continue
		}
		totalSpent += expense.AmountCents
		if expense.ExpenseDate >= monthStart && expense.ExpenseDate < monthEnd {
			spentThisMonth[expense.EnvelopeID] += expense.AmountCents
		}
	}

	envelopeSummaries := make([]model.EnvelopeSummary, 0, len(envelopes))
	for _, envelope := range envelopes {
		envelopeSummaries = append(envelopeSummaries, model.EnvelopeSummary{
			ID:             envelope.ID,
			Name:           envelope.Name,
			AllocatedCents: envelope.AllocatedCents,
			SpentCents:     spentThisMonth[envelope.ID],
			AvailableCents: envelope.BalanceCents,
		})
	}

	totalIncome := sum(incomeAmounts)
	totalAllocated := sum(allocatedAmounts)

	receiptTotal, err := r.totalIncomeReceiptCents(ctx, householdID)
	if err != nil {
		return nil, err
	}

	return &model.BudgetSummary{
		TotalIncomeCents:    totalIncome,
		TotalAllocatedCents: totalAllocated,
		UnallocatedCents:    totalIncome - totalAllocated + receiptTotal,
		TotalSpentCents:     totalSpent,
		Envelopes:           envelopeSummaries,
	}, nil
}

func (r *Repository) listEnvelopesWithoutFunding(ctx context.Context, householdID uuid.UUID) ([]model.Envelope, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, household_id, name, allocated_cents, balance_cents, created_at, updated_at
		 FROM envelopes WHERE household_id = $1 ORDER BY created_at DESC`,
		householdID,
	)
	if err != nil {
		return nil, fmt.Errorf("list envelopes: %w", err)
	}
	defer rows.Close()

	var envelopes []model.Envelope
	for rows.Next() {
		var envelope model.Envelope
		if err := rows.Scan(&envelope.ID, &envelope.HouseholdID, &envelope.Name, &envelope.AllocatedCents, &envelope.BalanceCents, &envelope.CreatedAt, &envelope.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan envelope: %w", err)
		}
		envelopes = append(envelopes, envelope)
	}
	return envelopes, rows.Err()
}

func sum(values []int64) int64 {
	var total int64
	for _, v := range values {
		total += v
	}
	return total
}
