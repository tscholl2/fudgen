ROOT = .

#directories
DATA_DIR := $(ROOT)/data

#files
SETUP_SCRIPT := $(ROOT)/setup/setup.go
DB := $(DATA_DIR)/db
SOURCE_NAMES := DATA_SRC DATSRCLN DERIV_CD FD_GROUP FOOD_DES FOOTNOTE LANGDESC LANGUAL NUT_DATA NUTR_DEF SRC_CD WEIGHT
SOURCE := $(addprefix $(DATA_DIR)/,$(addsuffix .txt,$(SOURCE_NAMES)))
PDF := $(DATA_DIR)/sr27_doc.pdf
ZIP := $(DATA_DIR)/sr27asc.zip
ss := $(DATA_DIR)/.source-sentinal


build: $(DB)

$(DB): $(SOURCE)
	touch $(DB)
	@command -v go >/dev/null 2>&1 || { echo >&2 "I require go but it's not installed.  Aborting."; exit 1; }
	go get -u github.com/mattn/go-sqlite3
	go get -u gopkg.in/yaml.v2
	go run $(SETUP_SCRIPT)

$(ZIP):
	mkdir -p $(DATA_DIR)
	wget -O $(ZIP) https://www.ars.usda.gov/SP2UserFiles/Place/12354500/Data/SR27/dnload/sr27asc.zip

$(SOURCE): $(ss)

$(ss): $(ZIP)
	unzip $(ZIP) -d $(DATA_DIR)
	touch $(ss)

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
