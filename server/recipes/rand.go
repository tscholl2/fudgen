package recipes

import (
	"fmt"
	"math/rand"

	"../units"
)

type randOperation struct {
	Name string
	Upper: float64
	Lower: float64
	Unit: string
	Type string
}

var operationList []string
var maxDepth int
var maxBreadth int

func init() {
	operationList = []randOperation{
		randOperation{
			Name: "grill"
			Upper: 30
			Lower: 2
			Unit: "min"
			Type: "cook"
		},
		randOperation{
			Name: "slice"
			Upper: 10
			Lower: 0.5
			Unit: "min"
			Type: "prep"
		}
	}
	maxDepth = 5
	maxBreadth = 3
}

func randomTitle(R *Recipe) string {
	return "untited"
}


//RandomRecipe returns a recipe and error
//public version of random rcipe function below
func RandomRecipe() (*Recipe, error) {
	return randomRecipe()
}

func randomRecipe() (r *Recipe, err error) {
	p, err := randomPreRecipe()
	if err != nil {
		return
	}
	s, err := preRecipe2Steps(p)
	if err != nil {
		return
	}
	return steps2recipe(s)
}

/*
type preRecipe struct {
	Name        string      //name of food/recipe/step
	Operation   string      //name of operation to make this step, nil for ingrediants
	Notes       string      //random notes to keep track of
	Time        string      //length of step, nil for ingrediants
	Quantity    string      //how much of ingrediant, e.g. "1/2 cup" or "3 slices"
	ID          int         //for keeping track
	Ingrediants []preRecipe //if empty then this is raw ingrediant
}
*/
//generates a random prerecipe using random foods from
//common table and also random operations from some list
//somewhere
func randomPreRecipe() (pr *preRecipe, err error) {
	p, err := randomPreRecipeHelper(maxDepth)
	if err != nil {
		return
	}
	pr = &p
	//first go through and set id's
	counter := 0
	var setID func(*preRecipe)
	setID = func(ptr *preRecipe) {
		ptr.ID = counter
		counter++
		//recurse into dependencies
		for k := 0; k < len(ptr.Ingrediants); k++ {
			setID(&(ptr.Ingrediants[k]))
		}
	}
	setID(pr)
	return
}

func randomPreRecipeHelper(level int) (pr preRecipe, err error) {
	if (rand.Float32() < 0.3 || level < 1) && level != maxDepth {
		ndbNo, err2 := randomNdbNo()
		amount := units.Quantity{Amount: (rand.Float64() * 2.5), Unit: "servings"}
		if err2 != nil {
			err = err2
			return
		}
		measurement, data, _, err2 := searchForFood(ndbNo, amount)
		if err2 != nil {
			err = err2
			return
		}
		//generate raw ingrediant
		pr.Name = data["Com_Desc"]
		pr.Quantity = fmt.Sprintf("%f %s", measurement.Amount, measurement.Unit)
	} else {
		//generate operation
		o := operationList[rand.Intn(len(operationList))]
		pr.Name = o.Name
		pr.Operation = o.Name
		t, err2 := randomTime(o.Lower, o.Upper, o.Unit)
		if err2 != nil {
			err = err2
			return
		}
		pr.Time = fmt.Sprintf("%f %s", t.Amount, t.Unit)
		n := rand.Intn(maxBreadth) + 1
		pr.Ingrediants = make([]preRecipe, n)
		for i := 0; i < n; i++ {
			pr2, err2 := randomPreRecipeHelper(level - 1)
			if err2 != nil {
				err = err2
				return
			}
			pr.Ingrediants[i] = pr2
		}
	}
	return
}

func randomTime(lowerBound float64, upperBound float64, unit string) (q units.Quantity, err error) {
	x := lowerBound + (upperBound-lowerBound)*rand.Float64()
	q, err = units.Parse(fmt.Sprintf("%f %s", x, unit))
	return
}
