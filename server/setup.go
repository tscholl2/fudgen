package main

import (
	"bufio"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"sync"
)

var DATA_PATH string
var DB_PATH string
var SCHEMA_PATH string

//allows use of all cores
func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	DATA_PATH = "./data/"
	SCHEMA_PATH = DATA_PATH + "schema.yml"
	DB_PATH = DATA_PATH + "db"
}

type T struct {
	Files     []string
	Columns   map[string][][]string
	Primary   map[string][]string
	Relations map[string][][]string
}

func cleanField(raw string) string {
	raw = strings.TrimSpace(raw) //remove extra spacing
	var clean string = ""
	for i := 0; i < len(raw); i++ {
		switch raw[i] {
		default:
			clean += string(raw[i])
		case 233: //remove non-ascii character é with ascii code 130, go code 233
			clean += "e"
			// case 181: //remove non-ascii character µ with ascii code 230, go code 181
			// 	clean += "u"
		}
	}
	if len(clean) > 1 && clean[0] == "~"[0] && clean[len(clean)-1] == "~"[0] {
		return string(clean[1 : len(clean)-1])
	} else {
		return clean
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
	schema_file, err := ioutil.ReadFile(SCHEMA_PATH)
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
	os.Remove(DB_PATH)

	//connect to db
	//see https://github.com/mattn/go-sqlite3/blob/master/_example/simple/simple.go
	//and http://godoc.org/github.com/mattn/go-sqlite3#SQLiteConn.Begin
	db, err := sql.Open("sqlite3", DB_PATH)
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
			cols = append(cols, fmt.Sprintf("%s\t\t%s", c[0], c[1]))
		}
		cmd := fmt.Sprintf("CREATE TABLE %s\n\t(", f)
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
	var wg sync.WaitGroup
	for _, f := range schema.Files {
		wg.Add(1)
		go func() {
			//open file
			file, err := os.Open(DATA_PATH + f + ".txt")
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
			fmt.Printf("loading file %s...\n", f)

			//build statement
			stupid_variable_go_needs := make([]string, len(schema.Columns[f]))
			for i := 0; i < len(stupid_variable_go_needs); i++ {
				stupid_variable_go_needs[i] = "?"
			}
			insrt_cmd := fmt.Sprintf("INSERT INTO %s VALUES (%s)", f, strings.Join(stupid_variable_go_needs, ","))
			stmt, err := tx.Prepare(insrt_cmd)
			if err != nil {
				fmt.Errorf("error: %v", err)
				panic(err)
			}

			//run over lines
			for scanner.Scan() {
				vals := extractColumns(scanner.Text())
				args := make([]interface{}, len(vals))
				for i, v := range vals {
					args[i] = interface{}(v)
				}
				_, err = stmt.Exec(args...)
				if err != nil {
					fmt.Errorf("error: %v", err)
					panic(err)
				}
			}

			//cleanup
			stmt.Close()
			tx.Commit()
			file.Close()
			wg.Done()
		}()
		wg.Wait()
	}
}
