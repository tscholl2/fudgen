#configurable stuff
INSTALL_DIR = www

#npm stuff
BOWER := node_modules/bower/bin/bower
BROWSERIFY := node_modules/browserify/bin/cmd.js
EXORCIST := node_modules/exorcist/bin/exorcist.js
COFFEE := node_modules/coffee-script/bin/coffee
NPM_DEPS := $(BOWER) $(BROWSERIFY) $(COFFEE) $(EXORCIST)

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

setup: $(npm-setinal)

$(npm-setinal): $(NPM_DEPS)
	touch $(npm-setinal)

$(NPM_DEPS):
	@command -v npm >/dev/null 2>&1 || { echo >&2 "I require npm but it's not installed.  Aborting."; exit 1; }
	npm install

#
#building
#

build: $(BUNDLE)

$(BUNDLE): $(COFFEE_FILES) $(COFFEE) $(BROWSERIFY) $(EXORCIST)
	@mkdir -p lib
	@command -v node >/dev/null 2>&1 || { echo >&2 "I require node but it's not installed.  Aborting."; exit 1; }
	node $(BROWSERIFY) -t coffeeify $(MAIN) --debug | ./$(EXORCIST) $(BUNDLE).map > $(BUNDLE)
#for production use this one
#node $(BROWSERIFY) -t coffeeify --extension=".coffee" src/main.coffee > lib/bundle.js

#
#installing
#
JS := $(wildcard lib/*)
CSS := $(wildcard css/*) node_modules/font-awesome/css/font-awesome.min.css node_modules/bootstrap/dist/css/bootstrap.min.css node_modules/bootstrap/dist/css/bootstrap-theme.min.css
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

.PHONY: update compile setup install default