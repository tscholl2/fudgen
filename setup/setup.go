package main

import (
	"bufio"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
)

type T struct {
	Files     []string
	Columns   map[string][][]string
	Primary   map[string][]string
	Relations map[string][][]string
}

func cleanField(raw string) string {
	raw = strings.TrimSpace(raw)
	if len(raw) > 1 && raw[0] == "~"[0] && raw[len(raw)-1] == "~"[0] {
		return string(raw[1 : len(raw)-1])
	} else {
		return raw
	}
}

func extractColumns(line string) (output []string) {
	output = strings.Split(line, "^")
	for i := 0; i < len(output); i++ {
		output[i] = cleanField(output[i])
	}
	return output
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
	os.Remove("../data/db")

	//connect to db
	//see https://github.com/mattn/go-sqlite3/blob/master/_example/simple/simple.go
	//and http://godoc.org/github.com/mattn/go-sqlite3#SQLiteConn.Begin
	db, err := sql.Open("sqlite3", "../data/db")
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

	//import files
	for _, f := range schema.Files {
		//open file
		file, err := os.Open("../data/" + f + ".txt")
		if err != nil {
			fmt.Errorf("error: %v", err)
			panic(err)
		}

		//setup line reading/db stuff
		reader := bufio.NewReader(file)
		scanner := bufio.NewScanner(reader)
		tx, err := db.Begin()
		if err != nil {
			fmt.Errorf("error: %v", err)
			panic(err)
		}
		fmt.Println("loading file %s could take a while...", f)

		//build statement
		stupid_variable_go_needs := make([]string, len(schema.Columns[f]))
		for i := 0; i < len(stupid_variable_go_needs); i++ {
			stupid_variable_go_needs[i] = "?"
		}
		insrt_cmd := fmt.Sprintf("INSERT INTO %s VALUES (%s)", f, strings.Join(stupid_variable_go_needs, ","))
		fmt.Println(insrt_cmd)
		stmt, err := tx.Prepare(insrt_cmd)
		if err != nil {
			fmt.Errorf("error: %v", err)
			panic(err)
		}

		//run over lines
		for scanner.Scan() {
			args := extractColumns(scanner.Text())
			fmt.Println(args)
			fmt.Println(len(args))
			//
			/*
				WHAT IS GOING ON HERE????????????????????????

				??????????????

				????????????????????????

				????????????????????????
			*/
			//
			_, err = stmt.Exec(args)
			if err != nil {
				fmt.Errorf("error: %v", err)
				panic(err)
			}
		}

		//cleanup
		stmt.Close()
		tx.Commit()
		file.Close()
	}
}
