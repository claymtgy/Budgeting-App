# Budgeting-App
A budgeting app based on the envelope system.

The app will utilize Python on the backend, a relational database such as MySQL, as well as React and MUI for the front end. 

## Current Project

### Number Crunching
Though simple, it's still going to take some time to get the number crunching part of the app done. It's being done first to 
just get a proof-of-concept in a CLI. The logic will eventually be re-factored into a Flask app, and made into a web-app. 

There are a few steps that any contributors should be aware of:

1. Currently, all logic and math is done just off of the local filesystem. The current implemented files that are created and 
stored are in the root directory of this repo. They should be under the ./incomes.json and the ./envelopes.json files. These
files will eventually be removed. When the app gets a bit more fleshed out, we will move to a SQLite database for development
and then it'll be any SQL-based database that all the information will be stored off of.
2. The number-crunching part is currently done in the ./number-crunching directory. The logic is still in infancy, and I'm 
still hashing out exactly how I want to have it done. The logic should follow: See if the incomes and expenses are existing 
-> run initialization steps if necessary -> present options of all functions, run them as the user requests -> close the app.
3. Most of the processes of 2 will obviously change as development progresses. The database schema is being worked on currently
and I have a pretty good idea of what I want to happen there. That will be laid out in a db directory soon.
