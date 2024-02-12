"""
This is the main python number crunching file.

If you're planning on making commits to this, please name global variables at the top,
and functions below that.

TO-DO:
    - Initialization needs a check
    - envelope_totals is not returning a value, it's returning a dict. This needs to be changed
    - Rename variables (remaining-total is currently for unallocated income funds, not accounted for by an envelope)
    - Create functionality to edit/delete individual envelopes
    - Implement default path?
        (may not be useful/practical when running on a web-server)
    - Ensure everything is a functio that can be called to start
    - Wrap the main loop in a function
    - Restructure files for better organization
    - Implement expenses
        - Expenses need to be housed in a dictionary themselves.
        - Any time a change is made to Expenses, there needs to be a call to calculate remaining total
        - Expenses details to be in a nested dictionary
        - Expenses need to be categorized from a categories list
        - Expenses need to be able to be voided. 
        - Expenses need to have a unique identifier
        - Expenses need to have an amount tied to them as well.

        ex:
        expenses: {
            1: [void: False, category: 4, name: 'expense-1', amount: -49, envelope: 1, transaction-id: 1]
            2: [void: True, categorry: 1, name: 'expense-2', amount: 50, envelope: 0, transaction-id: 2]
        }

        The expenses are in a nested dictionary.

        1 is the identifier, void is false, meaning run calculations on that particular category
        2 is the identifier, void is true, so don't calculate based on that.

        The string names will be nice for the front end to display a name for the transaction

        A negative number should indicate a negative subtraction from the total. In the examples above, expense 1 should ADD money
        to the envelope it's tied to

Write Tests
    - Test JSON saving
    - Test reading JSON on startup
    - Test INITIALIZATION is actually working
    - Test JSON deleting


LONGER-TERM:
    - Move all logic to a web server
"""

import os, json
from start_file import init_incomes, init_envelopes
from num_crunching import get_remaining_total, get_total_incomes, get_env_total
from json_creation import save_jsons, delete_jsons

# Global variables here
# Many of these could probably be changed from caps depending on what should change throughout the process. 
CONTINUE = True
INCOMES_LIST = {}
TOTAL_INCOME = 0
ENVELOPES = {}
DEFAULT_PATH = ''
INITIALIZATION = True

# This is not a constant variable, but it is global. 
envelope_totals = 0


HAS_GOTTEN_TOTAL_INCOME = False
HAS_GOTTEN_ENV_TOTAL = False
HAS_GOTTEN_REMANING_TOTAL = False


# Sets initialization to true or false,
def initialize():
    if os.path.isfile('../income.json') and os.path.isfile("../envelopes.json"):
        print("Initialization skipped, moving on.") 
        INITIALIZATION = False
        with open('../income.json', 'r') as f1:
            data_from_incomes = json.load(f1)
            INCOMES_LIST = data_from_incomes
            f1.close()
        with open('../envelopes.json', 'r') as f2:
            data_from_envelopes = json.load(f2)
            ENVELOPES = data_from_envelopes
            f2.close()


if INITIALIZATION:
    init_incomes(INCOMES_LIST)
    init_envelopes(ENVELOPES)

# Sets the main loop
while CONTINUE:
    selection = int(input('Press 1 to get total income, 2 to get envelope totals, or 3 to get remaining. 4 re-run initialization, 5 will save incomes. 6 will stop the program'))
    if selection == 1:
        print("Getting total income")
        TOTAL_INCOME = get_total_incomes(INCOMES_LIST)
    elif selection == 2:
        print("Getting envelope total")
        get_env_total(ENVELOPES)
    elif selection == 3:
        print("Getting remaining total")
        get_remaining_total(ENVELOPES,TOTAL_INCOME)
    elif selection == 4:
        print("Re-running  initialization")
        initialize()
    elif selection == 5:
        print("Saving income")
        save_jsons(INCOMES_LIST, ENVELOPES)
    elif selection == 6:
        CONTINUE = False
        print("Exiting")
    else:
        print("whoops, invalid choice")



TOTAL_INCOME = get_total_incomes(INCOMES_LIST)

# The get env totals is takign in incomes instead of envelopes
#get_env_total(INCOMES_LIST)

envelope_totals = get_env_total(ENVELOPES)

# This needs to be redone. There is too much happening in the maths side. Each equation should only be written once. We should
# be calling to retotal each thing instad of doing math more than once. 
# Any math or equation being done here, or nested in a funciton in num_crunching needs to be its own function in num_crunching
# right now I'm getting envelope totals too many times.
#get_remaining_total(ENVELOPES, TOTAL_INCOME)
'''
print(f"Your total unallocated funds are: \n{get_remaining_total(ENVELOPES, TOTAL_INCOME)}")

save_or_no = input("Would you like to save your incomes and envelopes? y/n\n")

if save_or_no == "y" or save_or_no == "yes":
    save_jsons(INCOMES_LIST, ENVELOPES)
elif save_or_no == "delete":
    delete_jsons()
'''
