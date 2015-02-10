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
	var str = string(b)
	fmt.Println("read file: \n", str)
	parsed, err := rparser.Parse(str)
	fmt.Println(parsed)
	fmt.Println(parsed.Name)
	fmt.Println(parsed.Recipe[0].Attributes)
	fmt.Println(parsed.Recipe[0].Ingrediants)
}
