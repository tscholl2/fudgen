package recipes

import (
	"../units"
	"fmt"
	"testing"
)

func init() {
	DB_PATH = "../data/db"
}

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

func TestPrice(t *testing.T) {
	_, d, _, e := searchForFood("Cream cheese", units.Quantity{Unit: "gram", Amount: 1, Type: "mass"})
	check(t, e, nil)
	check(t, d["NDB_No"], "01017")
	check(t, d["price"], fmt.Sprintf("%f", 0.026638))

	_, d, _, e = searchForFood("Cream cheese", units.Quantity{Unit: "gram", Amount: 2, Type: "mass"})
	check(t, e, nil)
	check(t, d["NDB_No"], "01017")
	check(t, d["price"], fmt.Sprintf("%f", 0.026638*2))
}
