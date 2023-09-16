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
"""

import os, json, json_creation

INCOMES_LIST = {}
TOTAL_INCOME = 0
ENVELOPES = {}
DEFAULT_PATH = ''
ITITIALIZATION = True

# Sets initialization to true or false,
if os.path.isfile('../income.json') and os.path.isfile("../envelopes.json"):
    print("Initialization skipped, moving on.") 
    INITIALIZATION = False
    with open('../income.json', 'r') as f:
        # Implement reading the JSON
        income_data = json.load(f)
        income = 0
    with open('../envelopes.json') as f:
        # Implement reading the JSON
        income_data = json.load(f)
        ENVELOPES = 0
else:
    print("Initialization is required. Please follow the prompts")
    add_incomes = True
    print("You have no current incomes listed")
    while add_incomes:
        inc_key = input("Name the income: \n")
        inc_amount = int(input("How much money is from this income?\n$"))
        INCOMES_LIST[inc_key] = inc_amount
        add_another = input("Do you want to add another income? y/n\n")
        print(f"Your current incomes are as follows: {INCOMES_LIST}")
        if add_another == "n" or add_another == "no":
            add_incomes = False
        


'''
# Checks if the income.json exists, if not, initializes it.
if os.path.isfile("../income.json"):
    print("This file exists")
    with open("../income.json", "r") as f:
        data = json.load(f)
    print(data)
else:
    incomes += int(input("No previous income declaration was found. What is your total monthly income?"))
    TOTAL_INCOME += incomes


if os.path.isfile("../envelopes.json"):
    print("This file exists.")
    with open("../envelopes.json") as f:
        envs = json.load(f)
    print(envs)
else:
    print("We haven't found any envelopes. Would you like to create any?")
'''