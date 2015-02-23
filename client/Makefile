
#npm stuff
BOWER := node_modules/bower/bin/bower
BROWSERIFY := node_modules/browserify/bin/cmd.js
COFFEE := node_modules/coffee-script/bin/coffee
NPM_DEPS := $(BOWER) $(BROWSERIFY) $(COFFEE)

#bower stuff
BACKBONE := bower_components/backbone/backbone.js
JQUERY := bower_components/jquery/dist/jquery.js
UNDERSCORE := bower_components/underscore/underscore.js
BOWER_DEPS := $(BACKBONE)

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

#compiling

compile: $(COFFEE)
	./$(COFFEE) --output lib --compile src

#update dependencies

update:
	npm install
	./$(BOWER) install

.PHONY: update compile setup