package recipes

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"strings"
)

func findNutrition(id string, msre string, quantity float32) (data map[string]string) {
	//connect to db
	db, err := sql.Open("sqlite3", "../../data/db")
	if err != nil {
		return
	}
	defer db.Close()

	//search for food
	sql := `select Msre_Desc,Amount,Gm_Wgt from WEIGHT where NDB_No=?`
	rows, err := db.Query(sql, id)
	if err != nil {
		return
	}
	defer rows.Close()

	//find matching unit
	var match sql.Row
	for rows.Next() {
		var Msre_Desc string
		var Amount float32
		var Gm_Wgt float32
		rows.Scan(&Msre_Desc, &Amount, &Gm_Wgt)
		//check for match
		if strings(Msre_Desc, msre) != -1 {
			grams := q * Gm_Wgt / 100.0
			//TODO
			/*
			def getNutrition(foodId,multiplier,nutrition):

			  with con:

			    cur = con.cursor()

			    cur.execute('select nutr_no,nutr_val from nutrition_data where ndb_no match "'+foodId+'"')

			    rows = cur.fetchall()

			    for row in rows:
			      id = int(row[0])
			      val = float(row[1])
			      cur2 = con.cursor()
			      cur2.execute('select units,NutrDesc from nutr_def where nutr_no == "'+str(id)+'"')
			      rows2 = cur2.fetchone()
			      units = rows2[0]
			      name = rows2[1]
			      if ord(units[0])==65533:
			        units = 'microgram'
			      if units == 'IU':
			        units = 'dimensionless'
			      if name in nutrition.keys():
			        nutrition[name.encode('utf-8')] = str(val*ureg.parse_expression(units)+ureg.parse_expression(nutrition[name.encode('utf-8')]))
			      else:
			        nutrition[name.encode('utf-8')] =str(val*ureg.parse_expression(units))


			  return nutrition
			*/
		}
	}
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
	return
}
