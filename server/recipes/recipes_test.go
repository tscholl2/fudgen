package recipes

import (
	"fmt"
	"io/ioutil"
	"testing"

	"../units"
)

func init() {
	DB_PATH = "../data/db"
}

func check(t *testing.T, value interface{}, expected interface{}) {
	if value != expected {
		t.Fail()
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

// TestJSON tests recipes go to json correctly
func TestJSON(t *testing.T) {
	// read whole the file
	b, err := ioutil.ReadFile("../../recipes/r3.yml")
	if err != nil {
		t.Error(err)
	}
	parsed, err := ParseYaml(string(b))
	if err != nil {
		t.Error(err)
	}
	//check op first
	s := parsed.Steps[0]
	exp := `{"name":"grilled cheese","desc":"grill","id":0,"time":{"unit":"minute","amount":10,"Type":"time"},"reqs":[1,2],"notes":"flip at some point"}`
	b, err = s.MarshalJSON()
	check(t, string(b), exp)

	//check ingrediant next
}

// TestQueries makes sure query functions
// work properly
func TestQueries(t *testing.T) {
	ndbno, err := findFood(`CHEESE,MEXICAN,QUESO ASADERO`)
	if err != nil {
		t.Error(err)
	}
	check(t, ndbno, `01166`)
}

// TestPrint just runs things
func TestPrint(t *testing.T) {
	// read whole the file
	b, err := ioutil.ReadFile("../../recipes/r3.yml")
	if err != nil {
		t.Error(err)
	}
	parsed, err := ParseYaml(string(b))
	if err != nil {
		t.Error(err)
	}
	//fmt.Println(parsed)
	fmt.Println("START STEPS=======")
	for _, s := range parsed.Steps {
		fmt.Println(s)
	}
	fmt.Println("END STEPS=========")
	fmt.Println(Schedule(parsed, 2))

	s := parsed.Steps[0]
	fmt.Println(s)
	b, err = s.MarshalJSON()
	fmt.Println(string(b))
	//fmt.Println(recipes.FindNutrition("01009", units.Quantity{Amount: 3, Unit: "slices"}))
}
