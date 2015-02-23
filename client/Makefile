
#npm stuff
BOWER := node_modules/bower/bin/bower
BROWSERIFY := node_modules/browserify/bin/browserify
COFFEE := node_modules/coffee-script/bin/coffee
NPM_DEPS := $(BOWER) $(BROWSERIFY) $(COFFEE)

#bower stuff
BACKBONE := bower_components/backbone/backbone.js
JQUERY := bower_components/jquery/dist/jquery.js
UNDERSCORE := bower_components/underscore/underscore.js
BOWER_DEPS := $(BACKBONE)

#other stuff???


$(NPM_DEPS): 
	npm install

$(BOWER_DEPS): $(BOWER)
	./$(BOWER) install

setup: $(COFFEE) $(BACKBONE)
	#compile coffee files
	./$(COFFEE) --output lib --compile src

update:
	npm install
	$(BOWER) install
	#git pull