package main

import (
	"./recipes"
	//"./units"
	"fmt"
	"io/ioutil"
)

func main() {
	// read whole the file
	b, err := ioutil.ReadFile("../recipes/r3.yml")
	if err != nil {
		panic(err)
	}
	parsed, err := recipes.ParseYaml(string(b))
	if err != nil {
		panic(err)
	}
	//fmt.Println(parsed)
	fmt.Println("START STEPS=======")
	for _, s := range parsed.Steps {
		fmt.Println(s)
	}
	fmt.Println("END STEPS=========")
	fmt.Println(recipes.Schedule(parsed, 2))

	s := parsed.Steps[0]
	fmt.Println(s)
	b, err = s.MarshalJSON()
	fmt.Println(string(b))
	//fmt.Println(recipes.FindNutrition("01009", units.Quantity{Amount: 3, Unit: "slices"}))
}
