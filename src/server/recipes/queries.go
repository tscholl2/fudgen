package recipes

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"strings"
)

func searchForFood(name string, servings float64) (measurement Quantity, data map[string]string, nutrition map[string]Quantity, err error) {
	//first find food
	data, err = findFood(name)
	if err != nil {
		return
	}
	//find nutritional information
	ndb_no := data["NDB_No"]
	measurement, nutrition, err = findNutrition(ndb_no, servings)
	if err != nil {
		return
	}
	return
}

func findNutrition(ndb_no string, servings float64) (measurement Quantity, nutrition map[string]Quantity, err error) {
	//connect to db
	db, err := sql.Open("sqlite3", "../../data/db")
	if err != nil {
		return
	}
	defer db.Close()

	//search for food
	sql := `
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
	row := db.QueryRow(sql, ndb_no)
	err = row.Scan(&Msre_Desc, &Amount, &Gm_Wgt)
	if err != nil {
		return
	}

	//record measurement
	measurement.Amount = servings * Amount
	measurement.Unit = Msre_Desc

	//find number of grams
	grams := Gm_Wgt * servings

	//collect nurtition information
	sql = `
			select
				NUT_DATA.Nutr_Val,
				NUTR_DEF.Units,
				NUTR_DEF.NutrDesc
			from NUT_DATA
			join NUTR_DEF on
				NUTR_DEF.Nutr_No=NUT_DATA.Nutr_No
			where NUT_DATA.NDB_No=?
		`
	rows, err := db.Query(sql, ndb_no)
	if err != nil {
		return
	}
	defer rows.Close()

	//initialize map
	nutrition = make(map[string]Quantity)

	//load data into map
	for rows.Next() {
		var nutr_val float64
		var units string
		var nutr_desc string
		rows.Scan(&nutr_val, &units, &nutr_desc)
		nutrition[nutr_desc] = Quantity{Unit: units, Amount: nutr_val * grams * 1.0 / 100}
	}

	return
}

func findFood(food string) (data map[string]string, err error) {
	data = map[string]string{}

	//connect to db
	db, err := sql.Open("sqlite3", "../../data/db")
	if err != nil {
		return
	}
	defer db.Close()

	//find food
	//sql := "select Long_Desc,NDB_No from FOOD_DES (abs(random()) % (select max(rowid)+1 from FOOD_DES where Long_Desc LIKE ?))"
	//sql := "select Long_Desc,NDB_No from FOOD_DES where Long_Desc like ? collate nocase order by random() limit 1"
	sql := `
		select
			FOOD_DES.Long_Desc,
			FOOD_DES.NDB_No
		from FOOD_DES
		where FOOD_DES.NDB_No=
		(
			select NDB_No from
			(
				select FOOD_DES.NDB_No
				from FOOD_DES
				join RANKING on RANKING.NDB_No=FOOD_DES.NDB_No 
				where FOOD_DES.Shrt_Desc like ? collate nocase
				order by -Google_Hits
				limit 15
			)
			order by random() limit 1
		)
	`
	rows, err := db.Query(sql, "%"+strings.ToUpper(food)+"%")
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var lng_desc string
		var id string
		rows.Scan(&lng_desc, &id)
		data["NDB_No"] = id
		data["Long_Desc"] = lng_desc
	}

	//check if didn't work
	_, ok := data["NDB_No"]
	if !ok {
		//search for food
		sql := `
		select
			FOOD_DES.Long_Desc,
			FOOD_DES.NDB_No
		from FOOD_DES
		where FOOD_DES.NDB_No=
		(
			select NDB_No from
			(
				select FOOD_DES.NDB_No
				from FOOD_DES
				join RANKING on RANKING.NDB_No=FOOD_DES.NDB_No
				order by -Google_Hits
				limit 200
			)
			order by random() limit 1
		)
		`
		var lng_desc string
		var id string
		row := db.QueryRow(sql)
		err = row.Scan(&lng_desc, &id)
		if err != nil {
			return
		}
		data["Long_Desc"] = lng_desc
		data["NDB_No"] = id
	}

	return
}
