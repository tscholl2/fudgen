package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
)

type T struct {
	Files     []string `yaml:"files"`
	Columns   map[string][][]string
	Primary   map[string][]string
	Relations map[string][][]string
}

func main() {
	//load yaml
	//see http://godoc.org/gopkg.in/yaml.v2
	schema_file, err := ioutil.ReadFile("./schema.yml")
	if err != nil {
		fmt.Errorf("error: %v", err)
		panic(err)
	}
	schema := T{}
	err = yaml.Unmarshal([]byte(string(schema_file)), &schema)
	if err != nil {
		fmt.Errorf("error: %v", err)
		panic(err)
	}

	//remove old db
	os.Remove("./db")

	//connect to db
	//see https://github.com/mattn/go-sqlite3/blob/master/_example/simple/simple.go
	//and http://godoc.org/github.com/mattn/go-sqlite3#SQLiteConn.Begin
	db, err := sql.Open("sqlite3", "./db")
	if err != nil {
		fmt.Errorf("error: %v", err)
		panic(err)
	}
	defer db.Close() //puts on stack for execution after function returns

	//load tables
	tx, err := db.Begin()
	if err != nil {
		fmt.Errorf("error: %v", err)
		panic(err)
	}
	for _, f := range schema.Files {
		pk := schema.Primary[f]
		fk := schema.Relations[f]
		cols := []string{}
		for _, c := range schema.Columns[f] {
			cols = append(cols, fmt.Sprintf("%s \t\t %s", c[0], c[1]))
		}
		cmd := fmt.Sprintf("CREATE TABLE %s \n\t(", f)
		cmd += fmt.Sprintf("\n\t\t%s", strings.Join(cols, ",\n\t\t"))
		if len(pk) > 0 {
			cmd += fmt.Sprintf(",\n\t\tPRIMARY KEY\t(%s)", strings.Join(pk, ", "))
		}
		if len(fk) > 0 {
			cmd += ","
		}
		for _, k := range fk {
			cmd += fmt.Sprintf("\n\t\tFOREIGN KEY(%s)\tREFERENCES %s(%s)", k[0], k[1], k[0])
		}
		cmd += "\n\t)"
		_, err = tx.Exec(cmd)
		if err != nil {
			fmt.Errorf("error: %v", err)
			panic(err)
		}
	}
	tx.Commit()

}
