package recipes

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func FindFood(food string) (data map[string]string, err error) {
	data = map[string]string{}

	//connect to db
	db, err := sql.Open("sqlite3", "../../data/db")
	if err != nil {
		return
	}
	defer db.Close()

	//find food
	var desc string
	var id string
	//sql := "select Long_Desc,NDB_No from FOOD_DES (abs(random()) % (select max(rowid)+1 from FOOD_DES where Long_Desc LIKE ?))"
	sql := "select Long_Desc,NDB_No from FOOD_DES where Long_Desc like ? order by random() limit 1"
	rows, err := db.Query(sql, "%"+food+"%")
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&desc, &id)
		data["NDB_No"] = id
		data["Long_Desc"] = desc
	}
	return
}
