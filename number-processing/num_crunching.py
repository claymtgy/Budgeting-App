# This file contains the main number crunching files

def get_total_incomes(INCOMES_LIST, TOTAL_INCOME):
    total = 0
    for i in INCOMES_LIST:
        print(INCOMES_LIST[i])
        total += int(INCOMES_LIST[i])
    print(total)
    return total

def get_env_total(ENVELOPES):
    total = 0
    for i in ENVELOPES:
        print(ENVELOPES[i])
        total += ENVELOPES[i]
        return total

def get_remaining_total(ENVELOPES,TOTAL_INCOME):
    total = get_env_total(ENVELOPES)
    remaining = TOTAL_INCOME - total
    return remaining