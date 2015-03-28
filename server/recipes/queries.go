package recipes

import (
	"database/sql"
	"fmt"
	"strconv"

	"../units"
	_ "github.com/mattn/go-sqlite3" //allows 'sqlite' driver in std sqldb lib
)

var dbPath string

func init() {
	dbPath = "./data/db"
}

func searchForFood(name string, amount units.Quantity) (measurement units.Quantity, data map[string]string, nutrition map[string]units.Quantity, err error) {

	//first find food
	ndbNo, err := findFood(name)
	if err != nil {
		return
	}

	//find nutritional information
	measurement, nutrition, err = findNutrition(ndbNo, amount)
	if err != nil {
		return
	}

	//gather basic details
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return
	}
	defer db.Close()
	sqlCmd := `
	select
		FOOD_DES.Long_desc,
		FOOD_DES.Shrt_Desc,
		COMMON.Com_Desc,
		FOOD_DES.ManufacName,
		PRICE.Mean_Price
	from FOOD_DES
	left join COMMON
		on COMMON.NDB_No=FOOD_DES.NDB_No
	left join PRICE
		on FOOD_DES.NDB_No = PRICE.NDB_No
	where FOOD_DES.NDB_No=?
	limit 1
	`
	var longDesc string
	var shrtDesc string
	var preComName []byte //might not always exist
	var comName string
	var manName string
	var price interface{}
	row := db.QueryRow(sqlCmd, ndbNo)
	err = row.Scan(&longDesc, &shrtDesc, &preComName, &manName, &price)
	if err != nil && err != sql.ErrNoRows {
		return
	}
	comName = string(preComName) //gives "" even if pre is nil

	//initialize and fill data map
	data = make(map[string]string)
	data["NDB_No"] = ndbNo
	data["Long_Desc"] = longDesc
	data["Shrt_Desc"] = shrtDesc
	data["Com_Desc"] = comName
	data["ManufacName"] = manName

	//compute price
	var p float64
	if price != nil {
		//see if the amount is standardized
		if amount.Type != "" {
			//see findNutrition below, same idea
			basicQuantity := amount.ToBasic()
			p = basicQuantity.Amount * price.(float64) //grams * price per gram
		} else {
			//see if we stored the total grams
			gmWgt, ok := nutrition["Gm_Wgt"]
			if ok {
				p = gmWgt.Amount * price.(float64)
			} else { //take a guess by hitting with serving
				p = amount.Amount * price.(float64)
			}
		}
	}
	data["price"] = fmt.Sprintf("%f", p)

	//finally return all the findings
	return
}

func findNutrition(ndbNo string, givenAmount units.Quantity) (measurement units.Quantity, nutrition map[string]units.Quantity, err error) {
	//connect to db
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return
	}
	defer db.Close()

	//search db for the following
	var MsreDesc string
	var Amount float64
	var GmWgt float64

	//search for food with same units
	sqlCmd := `
	select
		Msre_Desc,
		Amount,
		Gm_Wgt
	from WEIGHT
	where NDB_No=? and Msre_Desc=?
	limit 1
	`
	row := db.QueryRow(sqlCmd, ndbNo, givenAmount.Unit)
	err = row.Scan(&MsreDesc, &Amount, &GmWgt)
	//if couldn't find with same unit
	//then search for anything with same ndbno
	if err == sql.ErrNoRows {
		//search for food
		sqlCmd = `
		select
			Msre_Desc,
			Amount,
			Gm_Wgt
		from
		(
			select * from WEIGHT
			where NDB_No=?
		)
		order by random()
		limit 1
		`
		row := db.QueryRow(sqlCmd, ndbNo)
		err = row.Scan(&MsreDesc, &Amount, &GmWgt)
		if MsreDesc == "" || Amount == 0 {
			err = fmt.Errorf("[Q] No weight found for %s", ndbNo)
		}
		if err != nil {
			return
		}
	}

	//calculate nutrition
	var grams float64 //the number of grams in the ingrediant total
	//check if we can convert given_amount to grams
	//if so then we can just go with it
	if givenAmount.Type == "mass" || givenAmount.Type == "volume" {
		measurement = givenAmount
		basicAmount := givenAmount.ToBasic()
		grams = basicAmount.Amount //here we are assuming 1 gram = 1 ml
	} else {
		//else treat as servings and take what we found in the database
		//so we assume given_amount.Amount = # of servings
		//and multiply that against the amount in the database "Amount"
		//because that should be roughly one serving
		measurement.Amount = givenAmount.Amount * Amount
		measurement.Unit = MsreDesc
		grams = measurement.Amount * GmWgt
	}

	//collect nurtition information
	sqlCmd = `
			select
				NUT_DATA.Nutr_Val,
				NUTR_DEF.Units,
				NUTR_DEF.NutrDesc
			from NUT_DATA
			join NUTR_DEF on
				NUTR_DEF.Nutr_No=NUT_DATA.Nutr_No
			where NUT_DATA.NDB_No=?
			and NUTR_DEF.NutrDesc in (

"Alanine", "Alcohol, ethyl", "Arginine", "Ash", "Aspartic acid",
"Caffeine", "Calcium, Ca", "Carbohydrate, by difference",
"Carotene, alpha", "Carotene, beta", "Cholesterol", "Choline, total",
"Copper, Cu", "Cryptoxanthin, beta", "Cystine", "Energy",
"Fatty acids, total monounsaturated", "Fatty acids, total polyunsaturated",
"Fatty acids, total saturated", "Fiber, total dietary", "Folate, DFE",
"Folate, food", "Folate, total", "Folic acid", "Glutamic acid", "Glycine",
"Histidine", "Iron, Fe", "Isoleucine", "Leucine", "Lutein + zeaxanthin",
"Lycopene", "Lysine", "Magnesium, Mg", "Manganese, Mn", "Methionine",
"Niacin", "Pantothenic acid", "Phenylalanine", "Phosphorus, P", "Potassium, K",
"Proline", "Protein", "Retinol", "Riboflavin", "Selenium, Se", "Serine",
"Sodium, Na", "Sugars, total", "Theobromine", "Thiamin", "Threonine",
"Total lipid (fat)", "Tryptophan", "Tyrosine", "Valine", "Vitamin A, IU",
"Vitamin A, RAE", "Vitamin B-12", "Vitamin B-12, added", "Vitamin B-6",
"Vitamin C, total ascorbic acid", "Vitamin D", "Vitamin D (D2 + D3)",
"Vitamin D3 (cholecalciferol)", "Vitamin E (alpha-tocopherol)",
"Vitamin E, added", "Vitamin K (phylloquinone)", "Water", "Zinc, Zn"

				)
		`
	rows, err := db.Query(sqlCmd, ndbNo)
	if err != nil {
		return
	}
	defer rows.Close()

	//initialize map
	nutrition = make(map[string]units.Quantity)

	//load data into map
	for rows.Next() {
		var nutrVal float64
		var u string
		var nutrDesc string
		rows.Scan(&nutrVal, &u, &nutrDesc)
		//calculate nutritional density * amount of ingrediant, see pg 26 of sr27_doc.pdf
		nutrition[nutrDesc] = units.Quantity{Unit: u, Amount: nutrVal * grams * 1.0 / 100}
	}
	//store total weight
	nutrition["Gm_Wgt"] = units.Quantity{Unit: "gram", Amount: grams, Type: "mass"}
	//store # of servings for this thing
	nutrition["servings"] = units.Quantity{Unit: "serving", Amount: grams / GmWgt}
	return
}

func _queriesStrIsNdbNo(s string) bool {
	if len(s) != 5 {
		return false
	}
	_, err := strconv.Atoi(s)
	if err != nil {
		return false
	}
	return true
}

func findFood(food string) (ndbNo string, err error) {
	//if given an nbd_no immediately return it
	//this assumes no food is labeled something like
	//00123 or something else stupid
	if len(food) == 5 && _queriesStrIsNdbNo(food) {
		ndbNo = food
		return
	}

	//connect to db
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return
	}
	defer db.Close()

	//look for food in common first
	sqlCmd := `
		select NDB_No from
		(
			select COMMON.NDB_No, RANKING.Shrt_Hits
			from COMMON
			join RANKING on RANKING.NDB_No=COMMON.NDB_No
			where COMMON.Com_Desc like ? collate nocase
			order by -RANKING.Shrt_Hits
			limit 5
		)
		order by random()
		limit 1
	`
	row := db.QueryRow(sqlCmd, "%"+food+"%")
	err = row.Scan(&ndbNo)
	if err != nil && err != sql.ErrNoRows {
		return
	}
	if ndbNo != "" {
		return
	}

	//if can't find anything then try again
	//using whole database
	//find food by searching through db
	//using the rankings to find closest matches
	sqlCmd = `
		select NDB_No from
		(
			select FOOD_DES.NDB_No
			from FOOD_DES
			join RANKING on RANKING.NDB_No=FOOD_DES.NDB_No
			where FOOD_DES.Shrt_Desc like ? collate nocase
			order by -Shrt_Hits
			limit 15
		)
		order by random()
		limit 1
	`
	row = db.QueryRow(sqlCmd, "%"+food+"%")
	err = row.Scan(&ndbNo)
	if err != nil && err != sql.ErrNoRows {
		return
	}
	if ndbNo != "" {
		return
	}

	//check if didn't work
	//in case we need to expand the search
	//to pretty much anything really
	//search for food
	sqlCmd = `
	select NDB_No from
	(
		select FOOD_DES.NDB_No
		from FOOD_DES
		join RANKING on RANKING.NDB_No=FOOD_DES.NDB_No
		order by -Shrt_Hits
		limit 200
	)
	order by random()
	limit 1
	`
	row = db.QueryRow(sqlCmd)
	err = row.Scan(&ndbNo)
	if err != nil {
		return
	}
	return
}

func randomNdbNo() (ndbNo string, err error) {
	//connect to db
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return
	}
	defer db.Close()

	//search for food
	sqlCmd := `
		select NDB_No from
		(
			select NDB_No
			from COMMON
		)
		order by random()
		limit 1
	`
	row := db.QueryRow(sqlCmd)
	err = row.Scan(&ndbNo)
	if err != nil && err != sql.ErrNoRows {
		return
	}
	if ndbNo != "" {
		return
	}
	err = fmt.Errorf("[Q] Unable to generate random ndb_no")
	return
}
