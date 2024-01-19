# This file contains the main number crunching functions 

def get_total_incomes(INCOMES_LIST, TOTAL_INCOME):
    total = 0
    for i in INCOMES_LIST:
        print(f'Printing incomes list: {INCOMES_LIST[i]}')
        total += int(INCOMES_LIST[i])
    print(f'Returning total: {total}')
    return total

def get_env_total(ENVELOPES):
    print("getting envelope total"
    total = 0
    for i in ENVELOPES:
        print(ENVELOPES[i])
        total += ENVELOPES[i]
    print(f'Returning total: {total}')
    return total

def get_remaining_total(ENVELOPES,TOTAL_INCOME):
    total = get_env_total(ENVELOPES)
    print(f'Running the following equation: {TOTAL_INCOME} - {total}')
    remaining = TOTAL_INCOME - total
    return remaining
