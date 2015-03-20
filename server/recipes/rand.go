package recipes

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"

	"../units"
)

type randOperation struct {
	Name     string
	MinPerGm float64
	HighMin  float64
	LowMin   float64
	Type     string
	Notes    string
}

var operationList []randOperation
var defaultParameters RandomParameters

func init() {
	operationList = []randOperation{
		//prepping operations
		randOperation{
			Name:     "chop",
			Notes:    "very fine",
			MinPerGm: 0.6,
			LowMin:   1,
			HighMin:  5,
			Type:     "prep",
		},
		randOperation{
			Name:     "smash",
			MinPerGm: 1,
			LowMin:   5,
			HighMin:  15,
			Type:     "prep",
		},
		randOperation{
			Name:     "slice",
			Notes:    "thinly",
			MinPerGm: 0.3,
			LowMin:   1,
			HighMin:  5,
			Type:     "prep",
		},
		randOperation{
			Name:     "mix",
			MinPerGm: 0.2,
			LowMin:   1,
			HighMin:  5,
			Type:     "prep",
		},
		//cooking operations
		randOperation{
			Name:     "grill",
			Notes:    "use high heat",
			MinPerGm: 2,
			LowMin:   10,
			HighMin:  60,
			Type:     "cook",
		},
		randOperation{
			Name:     "bake",
			Notes:    "400 degrees",
			MinPerGm: 10,
			LowMin:   30,
			HighMin:  60,
			Type:     "cook",
		},
		randOperation{
			Name:     "fry",
			Notes:    "Medium pan on high heat",
			MinPerGm: 2,
			LowMin:   7,
			HighMin:  15,
			Type:     "cook",
		},
		randOperation{
			Name:     "roast",
			Notes:    "In oven",
			MinPerGm: 15,
			LowMin:   30,
			HighMin:  60,
			Type:     "cook",
		},
	}
	defaultParameters = RandomParameters{
		MaxDepth:              3,
		MaxRequirements:       3,
		ProbabilityIngrediant: 0.8,
		HowManyCooks:          2,
		TooManyCooks:          false,
		Servings:              4.5,
		Dollars:               50,
		Minutes:               120,
	}
}

func randomTitle(R *Recipe) string {
	return "untited"
}

//RandomParameters is just a container
//for all the options given to the random
//recipe generation
type RandomParameters struct {
	ProbabilityIngrediant float64 //default: 0.7
	MaxDepth              int     //default: 3
	MaxRequirements       int     //default: 3
	HowManyCooks          int     //default: 2
	TooManyCooks          bool    //default: false
	Servings              float64 //default: between 3 and 5
	Dollars               float64 //default: 50
	Minutes               float64 //default: 120
}

//RandomRecipe returns a recipe and error
//public version of random rcipe function below
func RandomRecipe(options RandomParameters) (*Recipe, error) {
	return randomRecipe(options)
}

func randomRecipe(options RandomParameters) (r *Recipe, err error) {
	p, err := randomPreRecipe(options)
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
type RandomParameters struct {
	ProbabilityIngrediant float64 //default: 0.7
	MaxDepth              int     //default: 3
	MaxRequirements       int     //default: 3
	HowManyCooks          int     //default: 2
	TooManyCooks          bool    //default: false
	Servings              float64 //default: between 3 and 5
	Dollars               float64 //default: 50
	Minutes               float64 //default: 120
}
*/
//generates a random prerecipe using random foods from
//common table and also random operations from some list
//somewhere
func randomPreRecipe(options RandomParameters) (pr *preRecipe, err error) {
	//parse options
	if options.MaxDepth == 0 {
		options.MaxDepth = defaultParameters.MaxDepth
	}
	if options.MaxRequirements == 0 {
		options.MaxRequirements = defaultParameters.MaxRequirements
	}
	if options.Dollars == 0 {
		options.Dollars = (rand.Float64()*0.5 + 0.75) * defaultParameters.Dollars
	}
	if options.Minutes == 0 {
		options.Minutes = (rand.Float64()*0.5 + 0.75) * defaultParameters.Minutes
	}
	if options.Minutes == 0 {
		options.Minutes = (rand.Float64()*0.5 + 0.75) * defaultParameters.Minutes
	}
	if options.Servings == 0 {
		options.Servings = (rand.Float64() + 0.5) * defaultParameters.Servings
	}
	//call initial function
	ptr, _, err := randomPreRecipeHelper(options, randPreHelperTotals{}, 0, randOperation{})
	if err != nil {
		return
	}
	pr = &ptr
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

type randPreHelperTotals struct {
	Servings float64
	Dollars  float64
	Minutes  float64
	Weight   float64
}

func randomPreRecipeHelper(options RandomParameters, totals randPreHelperTotals, depth int, previous randOperation) (pr preRecipe, postTotals randPreHelperTotals, err error) {
	//decide whether to build an ingrediant
	if (rand.Float64() < options.ProbabilityIngrediant || depth >= options.MaxDepth) && depth > 0 {
		//pick a random food using query function
		ndbNo, err2 := randomNdbNo()
		//choose some number of servings
		servings := units.Quantity{Amount: (rand.Float64()*1.25 + 0.75), Unit: "servings"}
		if err2 != nil {
			err = err2
			return
		}
		measurement, data, nutrition, err2 := searchForFood(ndbNo, servings)
		if err2 != nil {
			err = err2
			return
		}
		//generate raw ingrediant
		pr.Name = data["Com_Desc"]
		pr.Quantity = fmt.Sprintf("%f %s", measurement.Amount, measurement.Unit)
		//record information to send back up chain
		//store servings
		totals.Servings += servings.Amount
		//store price
		var price float64
		price, err = strconv.ParseFloat(data["price"], 64)
		if err != nil {
			price = 0
		}
		totals.Dollars += price
		//store weight
		weight, ok := nutrition["Gm_Wgt"]
		if ok {
			totals.Weight += weight.Amount
		}

		//return information
		postTotals = totals
	} else {
		//otherwise generate an operation and recurse
		op := operationList[rand.Intn(len(operationList))]
		for op.Type == previous.Type {
			op = operationList[rand.Intn(len(operationList))]
		}
		pr.Operation = op.Name
		//generate number of requirements
		n := rand.Intn(options.MaxRequirements) + 1
		pr.Ingrediants = make([]preRecipe, n)
		//build requirements recursively
		for i := 0; i < n; i++ {
			//make sure to track totals!
			subPr, totals, err2 := randomPreRecipeHelper(options, totals, depth+1, op)
			if err2 != nil {
				err = err2
				return
			}
			pr.Ingrediants[i] = subPr
			//update return information
			postTotals = totals
		}
		t := math.Max(math.Min(totals.Weight*op.MinPerGm, op.HighMin), op.LowMin)
		pr.Time = fmt.Sprintf("%f %s", t, "min")
	}
	return
}
