ROOT = .

BOWER := node_modules/bower/bin/bower
BROWSERIFY := node_modules/browserify/bin/browserify
COFFEE := node_modules/coffee-script/bin/coffee



NPM_DEPS := $(BOWER) $(BROWSERIFY) $(COFFEE)
BOWER_DEPS := 

$(NPM_DEPS): 
	npm install

$(BOWER_DEPS): $(NPM_DEPS)
	./$(BOWER) install

setup: $(NPM_DEPS) $(BOWER_DEPS)
	#compile coffee files
	$(COFFEE) --output lib --compile src