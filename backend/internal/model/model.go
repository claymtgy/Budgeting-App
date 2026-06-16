package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

type Income struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	Name        string    `json:"name"`
	AmountCents int64     `json:"amount_cents"`
	Period      string    `json:"period"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Envelope struct {
	ID             uuid.UUID `json:"id"`
	UserID         uuid.UUID `json:"user_id"`
	Name           string    `json:"name"`
	AllocatedCents int64     `json:"allocated_cents"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type Expense struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	EnvelopeID  uuid.UUID `json:"envelope_id"`
	AmountCents int64     `json:"amount_cents"`
	Description string    `json:"description"`
	ExpenseDate string    `json:"expense_date"`
	Voided      bool      `json:"voided"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type EnvelopeSummary struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	AllocatedCents int64     `json:"allocated_cents"`
	SpentCents     int64     `json:"spent_cents"`
	AvailableCents int64     `json:"available_cents"`
}

type BudgetSummary struct {
	TotalIncomeCents    int64             `json:"total_income_cents"`
	TotalAllocatedCents int64             `json:"total_allocated_cents"`
	UnallocatedCents    int64             `json:"unallocated_cents"`
	TotalSpentCents     int64             `json:"total_spent_cents"`
	Envelopes           []EnvelopeSummary `json:"envelopes"`
}
