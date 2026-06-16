package repository

import (
	"context"
	"errors"
	"fmt"

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
	rows, err := r.pool.Query(ctx,
		`SELECT id, household_id, name, allocated_cents, created_at, updated_at
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
		if err := rows.Scan(&envelope.ID, &envelope.HouseholdID, &envelope.Name, &envelope.AllocatedCents, &envelope.CreatedAt, &envelope.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan envelope: %w", err)
		}
		envelopes = append(envelopes, envelope)
	}
	return envelopes, rows.Err()
}

func (r *Repository) CreateEnvelope(ctx context.Context, householdID uuid.UUID, name string, allocatedCents int64) (*model.Envelope, error) {
	var envelope model.Envelope
	err := r.pool.QueryRow(ctx,
		`INSERT INTO envelopes (household_id, name, allocated_cents)
		 VALUES ($1, $2, $3)
		 RETURNING id, household_id, name, allocated_cents, created_at, updated_at`,
		householdID, name, allocatedCents,
	).Scan(&envelope.ID, &envelope.HouseholdID, &envelope.Name, &envelope.AllocatedCents, &envelope.CreatedAt, &envelope.UpdatedAt)
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
		 RETURNING id, household_id, name, allocated_cents, created_at, updated_at`,
		name, allocatedCents, id, householdID,
	).Scan(&envelope.ID, &envelope.HouseholdID, &envelope.Name, &envelope.AllocatedCents, &envelope.CreatedAt, &envelope.UpdatedAt)
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

func (r *Repository) ListExpenses(ctx context.Context, householdID uuid.UUID) ([]model.Expense, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, household_id, envelope_id, amount_cents, description, expense_date::text, voided, created_at, updated_at
		 FROM expenses WHERE household_id = $1 ORDER BY expense_date DESC, created_at DESC`,
		householdID,
	)
	if err != nil {
		return nil, fmt.Errorf("list expenses: %w", err)
	}
	defer rows.Close()

	var expenses []model.Expense
	for rows.Next() {
		var expense model.Expense
		if err := rows.Scan(&expense.ID, &expense.HouseholdID, &expense.EnvelopeID, &expense.AmountCents, &expense.Description, &expense.ExpenseDate, &expense.Voided, &expense.CreatedAt, &expense.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan expense: %w", err)
		}
		expenses = append(expenses, expense)
	}
	return expenses, rows.Err()
}

func (r *Repository) CreateExpense(ctx context.Context, householdID, envelopeID uuid.UUID, amountCents int64, description, expenseDate string) (*model.Expense, error) {
	var expense model.Expense
	err := r.pool.QueryRow(ctx,
		`INSERT INTO expenses (household_id, envelope_id, amount_cents, description, expense_date)
		 SELECT $1, $2, $3, $4, $5::date
		 WHERE EXISTS (SELECT 1 FROM envelopes WHERE id = $2 AND household_id = $1)
		 RETURNING id, household_id, envelope_id, amount_cents, description, expense_date::text, voided, created_at, updated_at`,
		householdID, envelopeID, amountCents, description, expenseDate,
	).Scan(&expense.ID, &expense.HouseholdID, &expense.EnvelopeID, &expense.AmountCents, &expense.Description, &expense.ExpenseDate, &expense.Voided, &expense.CreatedAt, &expense.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("create expense: %w", err)
	}
	return &expense, nil
}

func (r *Repository) VoidExpense(ctx context.Context, householdID, id uuid.UUID) (*model.Expense, error) {
	var expense model.Expense
	err := r.pool.QueryRow(ctx,
		`UPDATE expenses SET voided = TRUE, updated_at = NOW()
		 WHERE id = $1 AND household_id = $2
		 RETURNING id, household_id, envelope_id, amount_cents, description, expense_date::text, voided, created_at, updated_at`,
		id, householdID,
	).Scan(&expense.ID, &expense.HouseholdID, &expense.EnvelopeID, &expense.AmountCents, &expense.Description, &expense.ExpenseDate, &expense.Voided, &expense.CreatedAt, &expense.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("void expense: %w", err)
	}
	return &expense, nil
}

func (r *Repository) GetSummary(ctx context.Context, householdID uuid.UUID) (*model.BudgetSummary, error) {
	incomes, err := r.ListIncomes(ctx, householdID)
	if err != nil {
		return nil, err
	}

	envelopes, err := r.ListEnvelopes(ctx, householdID)
	if err != nil {
		return nil, err
	}

	expenses, err := r.ListExpenses(ctx, householdID)
	if err != nil {
		return nil, err
	}

	incomeAmounts := make([]int64, len(incomes))
	for i, income := range incomes {
		incomeAmounts[i] = income.AmountCents
	}

	allocatedAmounts := make([]int64, len(envelopes))
	for i, envelope := range envelopes {
		allocatedAmounts[i] = envelope.AllocatedCents
	}

	spentByEnvelope := make(map[uuid.UUID]int64)
	for _, expense := range expenses {
		if expense.Voided {
			continue
		}
		spentByEnvelope[expense.EnvelopeID] += expense.AmountCents
	}

	envelopeSummaries := make([]model.EnvelopeSummary, 0, len(envelopes))
	var totalSpent int64
	for _, envelope := range envelopes {
		spent := spentByEnvelope[envelope.ID]
		totalSpent += spent
		envelopeSummaries = append(envelopeSummaries, model.EnvelopeSummary{
			ID:             envelope.ID,
			Name:           envelope.Name,
			AllocatedCents: envelope.AllocatedCents,
			SpentCents:     spent,
			AvailableCents: envelope.AllocatedCents - spent,
		})
	}

	totalIncome := sum(incomeAmounts)
	totalAllocated := sum(allocatedAmounts)

	return &model.BudgetSummary{
		TotalIncomeCents:    totalIncome,
		TotalAllocatedCents: totalAllocated,
		UnallocatedCents:    totalIncome - totalAllocated,
		TotalSpentCents:     totalSpent,
		Envelopes:           envelopeSummaries,
	}, nil
}

func sum(values []int64) int64 {
	var total int64
	for _, v := range values {
		total += v
	}
	return total
}
