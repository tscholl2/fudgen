package recipes

import (
	"../units"
	"fmt"
	"testing"
)

func check(t *testing.T, value interface{}, expected interface{}) {
	if value != expected {
		t.Error("Expected ", expected, ", got ", value)
	}
}

func search(name string) (s string) {
	_, d, _, e := searchForFood(name, units.Quantity{})
	if e != nil {
		return e.Error()
	}
	return d["Long_Desc"]
}

func TestSearch(t *testing.T) {
	s := "cheese"
	fmt.Printf("Randomly searching for %s...\n", s)
	fmt.Println(search(s))
	fmt.Println(search(s))
	fmt.Println(search(s))
	s = "milk"
	fmt.Printf("Randomly searching for %s...\n", s)
	fmt.Println(search(s))
	fmt.Println(search(s))
	fmt.Println(search(s))
	s = "almonds"
	fmt.Printf("Randomly searching for %s...\n", s)
	fmt.Println(search(s))
	fmt.Println(search(s))
	fmt.Println(search(s))
}
