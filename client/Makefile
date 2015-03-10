#configurable stuff
INSTALL_DIR = www

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
COFFEE_FILES := $(shell find src -wholename 'src/*.coffee')

default: build

#
#setup
#

#setinal stuff
npm-setinal := node_modules/.ns
bower-setinal := bower_components/.bs

setup: $(npm-setinal) $(bower-setinal)

$(npm-setinal): $(NPM_DEPS)
	touch $(npm-setinal)

$(bower-setinal): $(BOWER_DEPS)
	touch $(bower-setinal)

$(NPM_DEPS):
	@command -v npm >/dev/null 2>&1 || { echo >&2 "I require npm but it's not installed.  Aborting."; exit 1; }
	npm install

$(BOWER_DEPS): $(BOWER)
	./$(BOWER) install

#
#building
#

build: $(BUNDLE)

$(BUNDLE): $(COFFEE_FILES)
	mkdir -p lib
	@command -v node >/dev/null 2>&1 || { echo >&2 "I require node but it's not installed.  Aborting."; exit 1; }
	node $(BROWSERIFY) -t coffeeify --extension=".coffee" $(MAIN) --debug | ./$(EXORCIST) $(BUNDLE).map > $(BUNDLE)

#for production use this one
#node $(BROWSERIFY) -t coffeeify --extension=".coffee" src/main.coffee > lib/bundle.js

#
#installing
#
JS := $(wildcard lib/*)
CSS := $(wildcard css/*) bower_components/fontawesome/css/font-awesome.min.css bower_components/bootstrap/dist/css/bootstrap.min.css bower_components/bootstrap/dist/css/bootstrap-theme.min.css
IMG := $(wildcard img/*)
HTML := $(wildcard ./*.html)
PREINSTALL := $(JS) $(CSS) $(IMG) $(HTML)
POSTINSTALL := $(addprefix $(INSTALL_DIR)/,$(PREINSTALL))
install: $(POSTINSTALL)

$(INSTALL_DIR)/%: %
	install -m 644 -D $< $@


#
#update dependencies
#
update:
	npm install
	./$(BOWER) install

.PHONY: update compile setup install default