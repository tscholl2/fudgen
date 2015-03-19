package recipes

import (
	"fmt"
	"math/rand"

	"../units"
)

var operationList []string

func init() {
	operationList = []string{
		"grill",
		"cook",
		"microwave",
	}
}

func randomTitle(R *Recipe) string {
	return "untited"
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
	return randomPreRecipeHelper(rand.Intn(5), 0)
}

func randomPreRecipeHelper(level int, id int) (pr *preRecipe, err error) {
	if rand.Float32() < 0.3 {
		ndbNo, err := randomNdbNo()
		amount := units.Quantity{Amount: (rand.Float64() * 2.5), Unit: "servings"}
		if err != nil {
			break
		}
		//measurement units.Quantity
		//data map[string]string
		//nutrition map[string]units.Quantity
		//err error
		measurement, data, nutrition, err := searchForFood(ndbNo, amount)
		if err != nil {
			return
		}
		pr.Name = data["Com_Desc"]
		pr.Notes = "randomly generated!"
		pr.Quantity = fmt.Sprintf("%f %s", measurement.Amount, measurement.Unit)
		pr.ID = id
	} else {
		pr.Operation = operationList[rand.Intn(len(operationList))]
		pr.Time, err = randomTime(60, 500, "minutes")
		if err != nil {
			return
		}
	}
	return
}

func randomTime(lowerBound float64, upperBound float64, unit string) (q units.Quantity, err error) {
	x := lowerBound + (upperBound-lowerBound)*rand.Float64()
	q, err := units.Parse(fmt.Sprintf("%f %s", x, unit))
	return
}
