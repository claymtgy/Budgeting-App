package service

import "testing"

func TestTotalIncomeCents(t *testing.T) {
	got := TotalIncomeCents([]int64{1000, 2500, 500})
	if got != 4000 {
		t.Fatalf("expected 4000, got %d", got)
	}
}

func TestTotalAllocatedCents(t *testing.T) {
	got := TotalAllocatedCents([]int64{300, 700})
	if got != 1000 {
		t.Fatalf("expected 1000, got %d", got)
	}
}

func TestUnallocatedCents(t *testing.T) {
	got := UnallocatedCents(5000, 3200)
	if got != 1800 {
		t.Fatalf("expected 1800, got %d", got)
	}
}

func TestEnvelopeSpentCents(t *testing.T) {
	got := EnvelopeSpentCents([]int64{100, 200, 50})
	if got != 350 {
		t.Fatalf("expected 350, got %d", got)
	}
}

func TestEnvelopeAvailableCents(t *testing.T) {
	got := EnvelopeAvailableCents(1000, 350)
	if got != 650 {
		t.Fatalf("expected 650, got %d", got)
	}
}
