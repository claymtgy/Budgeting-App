package service

func TotalIncomeCents(amounts []int64) int64 {
	var total int64
	for _, a := range amounts {
		total += a
	}
	return total
}

func TotalAllocatedCents(amounts []int64) int64 {
	var total int64
	for _, a := range amounts {
		total += a
	}
	return total
}

func UnallocatedCents(totalIncome, totalAllocated int64) int64 {
	return totalIncome - totalAllocated
}

func EnvelopeSpentCents(expenseAmounts []int64) int64 {
	var total int64
	for _, a := range expenseAmounts {
		total += a
	}
	return total
}

func EnvelopeAvailableCents(allocated int64, spent int64) int64 {
	return allocated - spent
}
