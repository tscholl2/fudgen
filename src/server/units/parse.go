package units

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	//"unicode/utf8"
)

var number_re *regexp.Regexp

func init() {
	number_re, _ = regexp.Compile("[.0-9/]+")
}

func Parse(s string) (q Quantity, err error) {
	if len(s) > 100 { //nope
		return
	}
	//remove extra space
	s = strings.TrimSpace(s)
	//find numbers
	index := number_re.FindStringIndex(s)
	number_part := s[index[0]:index[1]]
	//parse number
	x, err := ParseNumber(number_part)
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

func ParseNumber(s string) (x float64, err error) {
	s = strings.TrimSpace(s)
	i := strings.Index(s, "/")
	if i >= 0 && i < len(s)-1 {
		a, err := strconv.ParseFloat(s[:i], 64)
		b, err := strconv.ParseFloat(s[i+1:], 64)
		if err != nil || b == 0 {
			err = errors.New(fmt.Sprintf("Error paring number: [%s],[%s]", err, s))
		} else {
			x = a * 1.0 / b
		}
	} else {
		x, err = strconv.ParseFloat(s, 64)
	}
	return
}

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
