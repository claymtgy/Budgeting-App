import json, os

def save_jsons(INCOMES_LIST, ENVELOPES):
    with open('../income.json', 'w') as f1:
        json.dump(INCOMES_LIST, f1)
        f1.close()
    with open('../envelopes.json', 'w') as f2:
        json.dump(ENVELOPES, f2)
        f2.close()

def delete_jsons():
    if os.path.exists("../income.json"):
        os.remove("../income.json")
    else:
        print("No income file found. Skipping deletion.")
    if os.path.exists("../envelopes.json"):
        os.remove("../envelopes.json")
    else:
        print("No envelopes file found. Skipping deletion.")