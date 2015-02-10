package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

func main() {
	os.Remove("./db")

	db, err := sql.Open("sqlite3", "./db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
	create table foo (id integer not null primary key, name text);
	delete from foo;
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		fmt.Println("%q: %s\n", err, sqlStmt)
		return
	}

	db.Exec("insert into foo(id,name) values(?,?)", []string{"a", "b"}...)

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("insert into foo(id, name) values(?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	for i := 0; i < 5; i++ {
		args := []string{"a", fmt.Sprintf("%s", i)}
		_, err = stmt.Exec(args...)
		if err != nil {
			log.Fatal(err)
		}
	}
	tx.Commit()

	fmt.Println("done!")
}
