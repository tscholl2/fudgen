package main

import (
	"database/sql"
	"fmt"
	//"github.com/mattn/go-sqlite3"
)

func main() {
	fmt.Println("opening db")
	db, err := sql.Open("sqlite3", "./db")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	fmt.Println("opened!")
}
