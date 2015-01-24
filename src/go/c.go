package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func main() {
	// read whole the file
	b, err := ioutil.ReadFile("r.txt")
	if err != nil {
		panic(err)
	}
	var lines []string = strings.Split(string(b), "\n")
	var start int = 4
	var end int = endOfBlock(lines, start)
	fmt.Println(start, end)
}

func countIndent(str string) (n int) {
	n = 0
	for i := 0; i < len(str); i++ {
		// tab = 9 , ' ' = 32
		if str[i] == 9 || str[i] == 32 {
			n += 1
		} else {
			return
		}
	}
	return
}

func endOfBlock(lines []string, start int) (i int) {
	var indent int = countIndent(lines[start])
	i = start
	for i < len(lines) && countIndent(lines[i]) >= indent {
		i += 1
	}
	return
}
