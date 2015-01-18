ROOT = .

#directories
DATA_DIR := $(ROOT)/data
SETUP_DIR := $(ROOT)/setup

#commands
coffee := node_modules/coffee-script/bin/coffee

#files
SETUP_SCRIPT := dbsetup.coffee
DB := $(DATA_DIR)/db
SOURCE_NAMES := DATA_SRC DATSRCLN DERIV_CD FD_GROUP FOOD_DES FOOTNOTE LANGDESC LANGUAL NUT_DATA NUTR_DEF SRC_CD WEIGHT
ns := node_modules/.npm-sentinal
p := package.json
SOURCE := $(addprefix $(DATA_DIR)/,$(addsuffix .txt,$(SOURCE_NAMES)))
PDF := $(DATA_DIR)/sr27_doc.pdf
ZIP := $(DATA_DIR)/sr27asc.zip
ss := $(DATA_DIR)/.source-sentinal


build: $(DB)

$(DB): $(SOURCE) $(ns)
	touch $(DB)
	(cd setup; ../$(coffee) $(SETUP_SCRIPT))

$(ZIP):
	mkdir -p $(DATA_DIR)
	wget -O $(ZIP) https://www.ars.usda.gov/SP2UserFiles/Place/12354500/Data/SR27/dnload/sr27asc.zip
	
$(SOURCE): $(ss)

$(ss): $(ZIP)
	unzip $(ZIP) -d $(DATA_DIR)
	touch $(ss)
	
$(ns): $(p)
	@command -v npm >/dev/null 2>&1 || { echo >&2 "I require npm but it's not installed.  Aborting."; exit 1; }
	@command -v nodejs >/dev/null 2>&1 || { echo >&2 "I require nodejs but it's not installed.  Aborting."; exit 1; }
	npm install
	touch $(ns)

	
rebuild:
	$(RM) -f $(DB)
	$(MAKE) build
	
clean:
	$(RM) $(PDF)
	$(RM) $(SOURCE)
	$(RM) $(ZIP)
	$(RM) $(ss)

.PHONY: clean rebuild update
.DEFAULT: build
