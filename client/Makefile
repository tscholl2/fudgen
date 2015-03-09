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
	npm install

$(BOWER_DEPS): $(BOWER)
	./$(BOWER) install

#
#building
#

build: $(BUNDLE)

$(BUNDLE): $(COFFEE_FILES)
	mkdir -p lib
	node $(BROWSERIFY) -t coffeeify --extension=".coffee" $(MAIN) --debug | ./$(EXORCIST) $(BUNDLE).map > $(BUNDLE)

#for production use this one
#node $(BROWSERIFY) -t coffeeify --extension=".coffee" src/main.coffee > lib/bundle.js

#
#installing
#
JS := $(wildcard lib/*)
POST_JS := $(addprefix $(DEST_DIR)/,$(JS))
CSS := $(wildcard css/*)
POST_CSS := $(addprefix $(DEST_DIR)/,$(CSS)) $(DEST_DIR)/css/font-awesome.min.css $(DEST_DIR)/css/bootstrap.min.css $(DEST_DIR)/css/bootstrap-theme.min.css
POST_IMG := $(addprefix $(DEST_DIR)/,$(IMG))
HTML := $(wildcard html/*)
POST_HTML := $(subst html/,$(DEST_DIR)/,$(HTML))
PREINSTALL := $(JS) $(CSS) $(IMG) $(HTML)
POSTINSTALL := $(POST_JS) $(POST_CSS) $(POST_IMG) $(POST_HTML)
install: $(POSTINSTALL)

#css files
$(DEST_DIR)/css/font-awesome.min.css: bower_components/fontawesome/css/font-awesome.min.css
	install -m 644 -D $< $@
$(DEST_DIR)/css/bootstrap.min.css: bower_components/bootstrap/dist/css/bootstrap.min.css
	install -m 644 -D $< $@
$(DEST_DIR)/css/bootstrap-theme.min.css: bower_components/bootstrap/dist/css/bootstrap-theme.min.css
	install -m 644 -D $< $@
#the bottom is like this because we can't just use $(DESTDIR)/css/% as not all rules are the same
$(addprefix $(DEST_DIR)/,$(CSS)): $(CSS)
	install -m 644 -D $(patsubst $(DEST_DIR)/%,%,$<) $@

#javascript files
$(addprefix $(DEST_DIR)/,$(JS)): $(JS)
	install -m 644 -D $(patsubst $(DEST_DIR)/%,%,$<) $@

#images
$(addprefix $(DEST_DIR)/,$(IMG)): $(IMG)
	install -m 644 -D $(patsubst $(DEST_DIR)/%,%,$<) $@

#HTML
#????????????????????
#??????????
#??????????????
#?????????????????

test: $(POST_CSS)

#$(DEST_DIR)/%.html: html/%
#	install -m 644 -D $< $@


#
#update dependencies
#
update:
	npm install
	./$(BOWER) install

.PHONY: update compile setup install default