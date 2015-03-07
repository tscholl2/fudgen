package recipes

import (
	"database/sql"
	"errors"
	"fmt"

	"../units"
	_ "github.com/mattn/go-sqlite3"
)

var DB_PATH string

func init() {
	DB_PATH = "./data/db"
}

func searchForFood(name string, amount units.Quantity) (measurement units.Quantity, data map[string]string, nutrition map[string]units.Quantity, err error) {

	//first find food
	ndb_no, err := findFood(name)
	if err != nil {
		return
	}

	//find nutritional information
	measurement, nutrition, err = findNutrition(ndb_no, amount)
	if err != nil {
		return
	}

	//gather basic details
	db, err := sql.Open("sqlite3", DB_PATH)
	if err != nil {
		return
	}
	defer db.Close()
	sql_cmd := `
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
	var long_desc string
	var shrt_desc string
	var pre_com_name []byte //might not always exist
	var com_name string
	var man_name string
	var price interface{}
	row := db.QueryRow(sql_cmd, ndb_no)
	err = row.Scan(&long_desc, &shrt_desc, &pre_com_name, &man_name, &price)
	if err != nil && err != sql.ErrNoRows {
		return
	}
	com_name = string(pre_com_name) //gives "" even if pre is nil

	//initialize and fill data map
	data = make(map[string]string)
	data["NDB_No"] = ndb_no
	data["Long_Desc"] = long_desc
	data["Shrt_Desc"] = shrt_desc
	data["ComName"] = com_name
	data["ManufacName"] = man_name

	//compute price
	var p float64 = 0
	if price != nil {
		//see if the amount is standardized
		if amount.Type != "" {
			//see findNutrition below, same idea
			basic_quantity := amount.ToBasic()
			p = basic_quantity.Amount * price.(float64) //price per gram
		} else {
			p = amount.Amount * price.(float64)
		}

	}
	data["price"] = fmt.Sprintf("%f", p)

	//finally return all the findings
	return
}
func FindNutrition(s string, u units.Quantity) (units.Quantity, map[string]units.Quantity, error) {
	return findNutrition(s, u)
}
func findNutrition(ndb_no string, given_amount units.Quantity) (measurement units.Quantity, nutrition map[string]units.Quantity, err error) {
	//connect to db
	db, err := sql.Open("sqlite3", DB_PATH)
	if err != nil {
		return
	}
	defer db.Close()

	//search for food
	sql_cmd := `
	select
		Msre_Desc,
		Amount,
		Gm_Wgt
	from
	(
		select * from WEIGHT where NDB_No=?
	)
	order by random()
	limit 1
	`
	var Msre_Desc string
	var Amount float64
	var Gm_Wgt float64
	row := db.QueryRow(sql_cmd, ndb_no)
	err = row.Scan(&Msre_Desc, &Amount, &Gm_Wgt)
	if Msre_Desc == "" {
		err = errors.New("No food found!: " + ndb_no)
	}
	if err != nil {
		return
	}

	//calculate nutrition
	var grams float64 //the number of grams in the ingrediant total
	//check if we can convert given_amount to grams
	//if so then we can just go with it
	if given_amount.Type == "mass" || given_amount.Type == "volume" {
		measurement = given_amount
		basic_amount := given_amount.ToBasic()
		grams = basic_amount.Amount //here we are assuming 1 gram = 1 ml
	} else {
		//else treat as servings and take what we found in the database
		//so we assume given_amount.Amount = # of servings
		//and multiply that against the amount in the database "Amount"
		//because that should be roughly one serving
		measurement.Amount = given_amount.Amount * Amount
		measurement.Unit = Msre_Desc
		grams = measurement.Amount * Gm_Wgt
	}

	//collect nurtition information
	sql_cmd = `
			select
				NUT_DATA.Nutr_Val,
				NUTR_DEF.Units,
				NUTR_DEF.NutrDesc
			from NUT_DATA
			join NUTR_DEF on
				NUTR_DEF.Nutr_No=NUT_DATA.Nutr_No
			where NUT_DATA.NDB_No=?
		`
	rows, err := db.Query(sql_cmd, ndb_no)
	if err != nil {
		return
	}
	defer rows.Close()

	//initialize map
	nutrition = make(map[string]units.Quantity)

	//load data into map
	for rows.Next() {
		var nutr_val float64
		var u string
		var nutr_desc string
		rows.Scan(&nutr_val, &u, &nutr_desc)
		//calculate nutritional density * amount of ingrediant, see pg 26 of sr27_doc.pdf
		nutrition[nutr_desc] = units.Quantity{Unit: u, Amount: nutr_val * grams * 1.0 / 100}
	}

	return
}

func findFood(food string) (ndb_no string, err error) {
	//connect to db
	db, err := sql.Open("sqlite3", DB_PATH)
	if err != nil {
		return
	}
	defer db.Close()

	//look for food in common first
	sql_cmd := `
		select NDB_No from
		(
			select COMMON.NDB_No, RANKING.Shrt_Hits
			from COMMON
			join RANKING on RANKING.NDB_No=COMMON.NDB_No
			where COMMON.Com_Desc like ? collate nocase
			order by RANKING.Shrt_Hits
			limit 3
		)
		order by random()
		limit 1
	`
	row := db.QueryRow(sql_cmd, "%"+food+"%")
	err = row.Scan(&ndb_no)
	if err != nil && err != sql.ErrNoRows {
		return
	}
	if ndb_no != "" {
		return
	}

	//if can't find anything then try again
	//using whole database
	//find food by searching through db
	//using the rankings to find closest matches
	sql_cmd = `
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
	row = db.QueryRow(sql_cmd, "%"+food+"%")
	err = row.Scan(&ndb_no)
	if err != nil && err != sql.ErrNoRows {
		return
	}
	if ndb_no != "" {
		return
	}

	//check if didn't work
	//in case we need to expand the search
	//to pretty much anything really
	//search for food
	sql_cmd = `
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
	row = db.QueryRow(sql_cmd)
	err = row.Scan(&ndb_no)
	if err != nil {
		return
	}
	return
}
