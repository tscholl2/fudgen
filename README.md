#Food Recipe Generator

###Goals:

- ~~Transform SR27 food database from usda into sqlite db~~
- Ability to generate custom recipes from database based on various parameters
- Turn into an API server allowing users to submit calls

###Steps:

1. ~~Write code to import SR27 into sqlite~~
2. ~~Write basic scheduling program~~
3. Finalize recipe formats
4. Write recipe parser
5. Randomize parsed recipes
6. Compile nutrition info
7. Build small web ui

###Requirements:


######Server:

- make
- golang. I use 1.4.
- npm <--- coming soon. Forever is an awesome script to manage daemons like this

######Client:

- npm
- bower
- browserify

###Installation

- Run make file. Takes seconds on a real computer and about 50 minutes on a raspberry pi.

###Notes

- Stupid non-ascii characters. Use `grep --color='auto' -P -n "[^\x00-\x7F]" . -r` to find them. Their is one in "Entree" in the weight file. There is also a weird character in `µ` in nutr_def. These are swapped for something else when they go in the database.