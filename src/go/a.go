package main

import (
	"./rparser"
	"fmt"
	"io/ioutil"
)

func main() {
	// read whole the file
	b, err := ioutil.ReadFile("../../recipes/r.yml")
	if err != nil {
		panic(err)
	}
	parsed, err := rparser.ParseYaml(string(b))
	if err != nil {
		panic(err)
	}
	fmt.Println(parsed)
}
