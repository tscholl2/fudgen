package main

import (
	"./rparser"
	"fmt"
	"io/ioutil"
)

func main() {
	// read whole the file
	b, err := ioutil.ReadFile("r.txt")
	if err != nil {
		panic(err)
	}
	var str = string(b)
	fmt.Println("read file: \n", str)
	var parsed = rparser.Parse(str)
	fmt.Println(parsed)
}
