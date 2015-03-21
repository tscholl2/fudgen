package units

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	//"unicode/utf8"
)

var numberRe *regexp.Regexp

func init() {
	numberRe, _ = regexp.Compile("[.0-9/]+")
}

//Parse takes a string like "3 cups" and returns
//a Quantity object. If it can't parse will throw
//an error
func Parse(s string) (q Quantity, err error) {
	if len(s) > 100 { //nope
		return
	}
	//remove extra space
	s = strings.TrimSpace(s)
	//find numbers
	index := numberRe.FindStringIndex(s)
	if index == nil {
		fmt.Errorf("Unable to parse string: %s", s)
		return
	}
	numberPart := s[index[0]:index[1]]
	//parse number
	x, err := ParseNumber(numberPart)
	if err != nil {
		return
	}
	//parse unit
	u, err := ParseUnit(strings.TrimSpace(s[:index[0]] + s[index[1]:]))
	if err != nil {
		return
	}
	//set type
	_, ok := unit_type[u]
	if ok {
		q.Type = unit_type[u]
	} else {
		q.Type = ""
	}
	//return quantity
	q.Unit = u
	q.Amount = x
	return
}

//ParseNumber takes in a string `123.22` and returns
//a number. It is diff than strconv because it works
//with `1/2` or `3.2/4`
func ParseNumber(s string) (x float64, err error) {
	s = strings.TrimSpace(s)
	i := strings.Index(s, "/")
	if i >= 0 && i < len(s)-1 {
		a, err := strconv.ParseFloat(s[:i], 64)
		b, err := strconv.ParseFloat(s[i+1:], 64)
		if err != nil || b == 0 {
			err = fmt.Errorf("Error paring number: [%s],[%s]", err, s)
		} else {
			x = a * 1.0 / b
		}
	} else {
		x, err = strconv.ParseFloat(s, 64)
	}
	return
}

//ParseUnit takes a string like `cup` or `Cups`
//and returns the best guess for the unit of measurement
//see the table in basics.go for more info
//throws an error if it can't figure it out
func ParseUnit(s string) (u string, err error) {
	//clean up
	s = strings.TrimSpace(strings.ToLower(s))
	//check if exists in  alternatives
	_, ok := alternatives[s]
	if ok {
		u = alternatives[s]
	}
	//pretend dimensionless and return
	return
}
