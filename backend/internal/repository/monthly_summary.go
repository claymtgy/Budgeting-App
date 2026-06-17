package repository

import (
	"context"

	"github.com/claym/budgeting-app/internal/model"
	"github.com/google/uuid"
)

func (r *Repository) GetMonthlySummary(ctx context.Context, householdID uuid.UUID, month string) (*model.MonthlySummary, error) {
	start, end, err := monthBoundsFromYYYYMM(month)
	if err != nil {
		return nil, err
	}

	isCurrent := month == currentMonthYYYYMM()

	var envelopes []model.Envelope
	if isCurrent {
		if err := r.EnsureMonthlyFunding(ctx, householdID); err != nil {
			return nil, err
		}
		envelopes, err = r.listEnvelopesWithoutFunding(ctx, householdID)
	} else {
		envelopes, err = r.listEnvelopesWithoutFunding(ctx, householdID)
	}
	if err != nil {
		return nil, err
	}

	expenses, err := r.ListExpenses(ctx, householdID, start, end)
	if err != nil {
		return nil, err
	}

	spentByEnvelope := make(map[uuid.UUID]int64)
	var totalSpent int64
	for _, expense := range expenses {
		if expense.Voided {
			continue
		}
		spentByEnvelope[expense.EnvelopeID] += expense.AmountCents
		totalSpent += expense.AmountCents
	}

	envelopeSummaries := make([]model.EnvelopeSummary, 0, len(envelopes))
	var totalAllocated int64
	for _, envelope := range envelopes {
		spent := spentByEnvelope[envelope.ID]
		totalAllocated += envelope.AllocatedCents

		available := envelope.AllocatedCents - spent
		if isCurrent {
			available = envelope.BalanceCents
		}

		envelopeSummaries = append(envelopeSummaries, model.EnvelopeSummary{
			ID:             envelope.ID,
			Name:           envelope.Name,
			AllocatedCents: envelope.AllocatedCents,
			SpentCents:     spent,
			AvailableCents: available,
		})
	}

	return &model.MonthlySummary{
		Month:               month,
		IsCurrentMonth:      isCurrent,
		TotalAllocatedCents: totalAllocated,
		TotalSpentCents:     totalSpent,
		Envelopes:           envelopeSummaries,
	}, nil
}
