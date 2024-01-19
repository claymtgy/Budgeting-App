# This file contains the necessary functions to initialize the program

def init_for_local(INITIALIZATION):
    if INITIALIZATION:
        init_incomes()
        init_envelopes()
        INITIALIZATION = False
    else:
        print("No init needed")

def init_incomes(INCOMES_LIST):
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
    print("Finished adding incomes")

def init_envelopes(ENVELOPES):
    print("Initialization is required. Please follow the prompts")
    add_envelopes = True
    print("You have no current envelopes listed")
    while add_envelopes:
        inc_key = input("Name the envelopes: \n")
        inc_amount = int(input("How much money is from this envelopes?\n$"))
        ENVELOPES[inc_key] = inc_amount
        add_another = input("Do you want to add another envelope? y/n\n")
        print(f"Your current envelopes are as follows: {ENVELOPES}")
        if add_another == "n" or add_another == "no":
            add_envelopes = False   
            print("Finished adding envelopes") 
