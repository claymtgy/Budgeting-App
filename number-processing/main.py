"""
This is the main python number crunching file.

If you're planning on making commits to this, please name global variables at the top,
and functions below that.

To-Do:
    - The Envelopes need to be a dict. The dictionary should be appended every time a new envelope
            is created.
    - Create global variable that sets initiating to true by default. Then changes that to false if all required files exist
    - If initiation = false, then run the else statements for stuff.
"""

import os, json, json_creation

TOTAL_INCOME = 0
ENVELOPES = {}
DEFAULT_PATH = ''


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

if os