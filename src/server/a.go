package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	//open file
	file, err := os.Open("../../data/WEIGHT.txt")
	if err != nil {
		fmt.Errorf("error: %v", err)
		panic(err)
	}

	//setup line reading stuff
	reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(reader)

	//run over lines
	for i := 0; i < 9430; i++ {
		scanner.Scan()
	}
	s := scanner.Text()
	fmt.Println(s)
	for i := 0; i < len(s); i++ {
		//fmt.Println(s[i])
		if s[i] == 233 {
			fmt.Println("found 233!")
		}
	}
	fmt.Println(strings.ContainsRune(s, 233))

}
