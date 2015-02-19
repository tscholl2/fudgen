package main

import (
	"./recipes"
	"fmt"
	"io/ioutil"
)

func main() {
	// read whole the file
	b, err := ioutil.ReadFile("../../recipes/r3.yml")
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
}
