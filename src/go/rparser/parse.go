package rparser

import (
	"fmt"
	"strings"
)

type ingrediant struct {
	raw        string
	attributes []string
}

type operation struct {
	raw        string
	attributes []string
}

type expression struct {
	name        string
	operation   operation
	ingrediants []ingrediant
	components  []expression
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

func getLineType(str string) int {
	return strings.TrimSpace(str)[0]
}

func isOperator(str string) bool {
	return getLineType(str) == 43 // + = 43
}

func isIngrediant(str string) bool {
	return getLineType(str) == 45 // - = 45
}

func isAttribute(str string) bool {
	return getLineType(str) == 46 // . = 46
}

func isName(str string) bool {
	return !isIngrediant(str) && !isOperator(str) && !isAttribute(str)
}

func parseIngrediant(lines []string) (i ingrediant) {
	i.raw = strings.TrimSpace(string(lines[0][1:]))
	j = 1
	for j < len(lines) && getLineType(lines[j]) == 46 { // . = 46
		append(i.attributes, strings.TrimSpace(string(lines[j])))
	}
	return
}

func parseOperator(lines []string) (o operation) {
	o.raw = strings.TrimSpace(string(lines[0][1:]))
	j = 1
	for j < len(lines) && getLineType(lines[j]) == 46 { // . = 46
		append(o.attributes, strings.TrimSpace(string(lines[j])))
	}
	return
}

func parseComponent(lines []string) (e expression) {
	var current int = 0
	//check for name
	if isName(string(lines[0])) {
		e.name = strings.TrimSpace(top)
		current += 1
	} else {
		e.name = "step"
	}
	//move through lines
	for current < len(lines) {
		if isOperator(string(lines[current])) {

			continue
		}
	}
	if isIngrediant(top) {
		var end int = endOfBlock(lines, 0)
		append(e.ingrediants, parseIngrediant(lines[0:end]))
	}
	if isOperator(top) {

	}
	if isAttribute(top) {

	}
}

func Parse(s string) (recipe expression) {
	var lines []string = strings.Split(s, "\n")
	fmt.Println("split lines:\n", lines)
	var name string = lines[0]
	var indent int = countIndent(lines[1])
	recipe.operation.raw = name
	fmt.Println(indent)

	//gather rest of lines

	return
}
