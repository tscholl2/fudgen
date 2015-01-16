ROOT = .

DATA_DIR := $(ROOT)/data
SETUP_DIR := $(ROOT)/setup

SETUP_SCRIPT := dbsetup.coffee
DB := $(DATA_DIR)/db
SOURCE_NAMES := DATA_SRC DATSRCLN DERIV_CD FD_GROUP FOOD_DES FOOTNOTE LANGDESC LANGUAL NUT_DATA NUTR_DEF SRC_CD WEIGHT
SOURCE := $(addprefix $(DATA_DIR)/,$(addsuffix .txt,$(SOURCE_NAMES)))
PDF := $(DATA_DIR)/sr27_doc.pdf
ZIP := $(DATA_DIR)/sr27asc.zip
ss := $(DATA_DIR)/.source-sentinal


build: $(DB)

$(DB): $(SOURCE)
	touch $(DB)
	(cd setup; coffee $(SETUP_SCRIPT))

$(ZIP):
	wget -O $(ZIP) https://www.ars.usda.gov/SP2UserFiles/Place/12354500/Data/SR27/dnload/sr27asc.zip
	
$(SOURCE): $(ss)

$(ss): $(ZIP)
	mkdir -p $(DATA_DIR)
	unzip $(ZIP) -d $(DATA_DIR)
	touch $(ss)

rebuild:
	$(RM) $(DB)
	$(MAKE) build
	
clean:
	$(RM) $(PDF)
	$(RM) $(SOURCE)
	$(RM) $(ZIP)
	$(RM) $(ss)

.PHONY: clean rebuild
.DEFAULT: build
