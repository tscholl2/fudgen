#configurable stuff
DEST_DIR = www

#npm stuff
BOWER := node_modules/bower/bin/bower
BROWSERIFY := node_modules/browserify/bin/cmd.js
EXORCIST := node_modules/exorcist/bin/exorcist.js
COFFEEIFY := node_modules/coffeeify/index.js
COFFEE := node_modules/coffee-script/bin/coffee
NPM_DEPS := $(BOWER) $(BROWSERIFY) $(COFFEE) $(COFFEEIFY) $(EXORCIST)

#bower stuff
BACKBONE := bower_components/backbone/backbone.js
JQUERY := bower_components/jquery/dist/jquery.js
UNDERSCORE := bower_components/underscore/underscore.js
BOWER_DEPS := $(BACKBONE) $(JQUERY) $(UNDERSCORE)

#files
BUNDLE := lib/bundle.js
MAIN := src/main.coffee

#setinal stuff
npm-setinal := node_modules/.ns
bower-setinal := bower_components/.bs

#download dependencies

setup: $(npm-setinal) $(bower-setinal)

$(npm-setinal): $(NPM_DEPS)
	touch $(npm-setinal)

$(bower-setinal): $(BOWER_DEPS)
	touch $(bower-setinal)

$(NPM_DEPS):
	npm install

$(BOWER_DEPS): $(BOWER)
	./$(BOWER) install

#installing

COFFEE_FILES := $(shell find src -wholename 'src/*.coffee')

build: $(JS_FILES)
	node $(BROWSERIFY) -t coffeeify --extension=".coffee" src/main.coffee --debug | ./$(EXORCIST) $(BUNDLE).map > $(BUNDLE)
	#for production use this one
	#node $(BROWSERIFY) -t coffeeify --extension=".coffee" src/main.coffee > lib/bundle.js

#update dependencies

update:
	npm install
	./$(BOWER) install

.PHONY: update compile setup install build
