"""
This is the main python number crunching file.

If you're planning on making commits to this, please name global variables at the top,
and functions below that.

TO-DO:
    - The Envelopes need to be a dict. The dictionary should be appended every time a new envelope
            is created.
    - Create global variable that sets initiating to true by default. Then changes that to false if all required files exist
    - If initiation = false, then run the else statements for stuff.
    - Incomes should be kept in a dictionary
    - Derive the total income from the incomes list
    - Implememt reading the JSONs
    - Implement saving the JSON
    - Properly format the JSON
    - Implement save JSON function

"""

import os, json
from start_file import init_incomes, init_envelopes
from num_crunching import get_remaining_total, get_total_incomes, get_env_total
from json_creation import save_jsons, delete_jsons

INCOMES_LIST = {}
TOTAL_INCOME = 0
ENVELOPES = {}
DEFAULT_PATH = ''
INITIALIZATION = True
envelope_totals = 0

# Sets initialization to true or false,
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

TOTAL_INCOME = get_total_incomes(INCOMES_LIST, TOTAL_INCOME)

get_env_total(INCOMES_LIST)

envelope_totals = get_remaining_total(ENVELOPES, TOTAL_INCOME)

print(f"Your total unallocated funds are: \n{get_remaining_total(ENVELOPES, TOTAL_INCOME)}")

save_or_no = input("Would you like to save your incomes and envelopes? y/n\n")

if save_or_no == "y" or save_or_no == "yes":
    save_jsons(INCOMES_LIST, ENVELOPES)
elif save_or_no == "delete":
    delete_jsons()