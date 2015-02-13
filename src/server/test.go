package main

import (
	"./recipes"
	"fmt"
	"io/ioutil"
)

func main() {
	// read whole the file
	b, err := ioutil.ReadFile("../../recipes/r2.yml")
	if err != nil {
		panic(err)
	}
	parsed, err := recipes.ParseYaml(string(b))
	if err != nil {
		panic(err)
	}
	fmt.Println("START STEPS=======")
	for _, s := range parsed {
		fmt.Println(s)
	}
	fmt.Println("END STEPS=========")

	//test queries
	f, err := recipes.FindFood("butter, salted")
	fmt.Println(err)
	fmt.Println(f)
}
