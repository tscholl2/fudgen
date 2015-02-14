package recipes

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"strings"
)

func FindNutrition(steps []Step) {

}

func FindFood(food string) (data map[string]string, err error) {
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
			FOOD_DES.NDB_No,
			WEIGHT.Amount,
			WEIGHT.GM_Wgt,
			WEIGHT.Msre_Desc
		from FOOD_DES
		join WEIGHT on FOOD_DES.NDB_No=WEIGHT.NDB_No
		where FOOD_DES.NDB_No=
		(
			select * from
			(
				select FOOD_DES.NDB_No
				from FOOD_DES
				join RANKING on RANKING.NDB_No=FOOD_DES.NDB_No 
				where Shrt_Desc like ? collate nocase
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
		var amt float32
		var msre_des string
		var gm_wgt float32
		rows.Scan(&lng_desc, &id, &amt, &msre_des, &gm_wgt)
		data["NDB_No"] = id
		data["Long_Desc"] = lng_desc
		data["Amount"] = amt
		data["Msre_Desc"] = msre_des
		data["GM_Wgt"] = gm_wgt
	}
	return
}
