package recipes

import (
	"fmt"
	"io/ioutil"
	"testing"

	"../units"
)

func init() {
	dbPath = "../data/db"
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
	fmt.Println(Schedule(parsed))

	s := parsed.Steps[0]
	fmt.Println(s)
	b, err = s.MarshalJSON()
	fmt.Println(string(b))
	//fmt.Println(recipes.FindNutrition("01009", units.Quantity{Amount: 3, Unit: "slices"}))
}

func TestRand(t *testing.T) {
	fmt.Println("Testing random prerec...")
	pr, err := randomPreRecipe(RandomParameters{})
	check(t, err, nil)
	fmt.Println(pr)
	pr, err = randomPreRecipe(RandomParameters{})
	check(t, err, nil)
	fmt.Println(pr)
	pr, err = randomPreRecipe(RandomParameters{})
	check(t, err, nil)
	fmt.Println(pr)

	fmt.Println("Testing random rec...")
	_, err = randomRecipe(RandomParameters{})
	//fmt.Println(r)
	check(t, err, nil)
	_, err = randomRecipe(RandomParameters{})
	check(t, err, nil)
	//fmt.Println(r)
	_, err = randomRecipe(RandomParameters{})
	check(t, err, nil)
	//fmt.Println(r)
}

func TestSchedule(t *testing.T) {
	fmt.Println("Testing scheudling...")
	v1 := "A"
	v2 := "B"
	v4 := "D"
	v3 := "C"
	v5 := "E"
	V := map[interface{}]int{v1: 10, v2: 10, v3: 10, v4: 100, v5: 10}
	E := [][]interface{}{{v1, v2}, {v2, v3}, {v2, v4}, {v4, v5}, {v3, v5}}
	h, err := optimalSchedule(V, E)
	check(t, err, nil)
	fmt.Println(h)
}
