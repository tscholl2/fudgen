package recipes

import (
	"math/rand"
	"sort"
	"strconv"

	"../units"
)

type randOperation struct {
	Name     string
	MinPerGm float64
	ETA      []float64 //time estimate [lowerbound,upperbound]
	MaxInput int       //max number of inputs
	Type     string
	Notes    string
}

var operationList []randOperation
var defaultParameters RandomParameters

func init() {
	operationList = []randOperation{
		//prepping operations
		randOperation{
			Name:     "mix",
			MinPerGm: 0.2,
			ETA:      []float64{1, 5},
			MaxInput: -1,
			Type:     "prep",
		},
		randOperation{
			Name:     "chop",
			Notes:    "very fine",
			MinPerGm: 0.6,
			ETA:      []float64{1, 5},
			MaxInput: 1,
			Type:     "prep",
		},
		randOperation{
			Name:     "smash",
			MinPerGm: 1,
			ETA:      []float64{1, 15},
			MaxInput: 1,
			Type:     "prep",
		},
		randOperation{
			Name:     "slice",
			Notes:    "thinly",
			MinPerGm: 0.3,
			ETA:      []float64{1, 5},
			MaxInput: 1,
			Type:     "prep",
		},
		randOperation{
			Name:     "layer",
			MinPerGm: 0.5,
			ETA:      []float64{1, 5},
			MaxInput: -1,
			Type:     "prep",
		},
		randOperation{
			Name:     "toss",
			MinPerGm: 0.3,
			ETA:      []float64{1, 5},
			MaxInput: -1,
			Type:     "prep",
		},
		//cooking operations
		randOperation{
			Name:     "grill",
			Notes:    "use high heat",
			MinPerGm: 2,
			ETA:      []float64{10, 60},
			MaxInput: -1,
			Type:     "cook",
		},
		randOperation{
			Name:     "bake",
			Notes:    "400 degrees",
			MinPerGm: 10,
			ETA:      []float64{30, 60},
			MaxInput: 1,
			Type:     "cook",
		},
		randOperation{
			Name:     "fry",
			Notes:    "Medium pan on high heat",
			MinPerGm: 2,
			ETA:      []float64{7, 15},
			MaxInput: -1,
			Type:     "cook",
		},
		randOperation{
			Name:     "roast",
			Notes:    "In oven",
			MinPerGm: 15,
			ETA:      []float64{30, 60},
			MaxInput: 1,
			Type:     "cook",
		},
	}
	defaultParameters = RandomParameters{
		HowManyCooks: 2,
		TooManyCooks: false,
		Servings:     4.5,
		Dollars:      50,
		Minutes:      120,
	}
}

func randomTitle(R *Recipe) string {
	return "untited"
}

//RandomParameters is just a container
//for all the options given to the random
//recipe generation
type RandomParameters struct {
	HowManyCooks int     //default: 2
	TooManyCooks bool    //default: false
	Servings     float64 //default: between 3 and 5
	Dollars      float64 //default: 50
	Minutes      float64 //default: 120
}

//RandomRecipe generates a random prerecipe using random foods from
//common table and also random operations from some list somewhere
func RandomRecipe(options RandomParameters) (*Recipe, error) {
	return randomRecipe(options)
}

/*
type Recipe struct {
	Steps     []Step                    `json:"steps"`
	Title     string                    `json:"title"`
	Nutrition map[string]units.Quantity `json:"nutr"`
	Price     float64                   `json:"price"`
}
type RandomParameters struct {
	HowManyCooks          int     //default: 2
	TooManyCooks          bool    //default: false
	Servings              float64 //default: between 3 and 5
	Dollars               float64 //default: 50
	Minutes               float64 //default: 120
}
*/
func randomRecipe(options RandomParameters) (r *Recipe, err error) {
	//initailize output
	r = &Recipe{}
	r.Steps = []Step{}
	r.Nutrition = make(map[string]units.Quantity)

	//parse options
	if options.Dollars == 0 {
		options.Dollars = (rand.Float64()*0.5 + 0.75) * defaultParameters.Dollars
	}
	if options.Minutes == 0 {
		options.Minutes = (rand.Float64()*0.5 + 0.75) * defaultParameters.Minutes
	}
	if options.HowManyCooks == 0 {
		options.HowManyCooks = defaultParameters.HowManyCooks
	}
	if options.Servings == 0 {
		options.Servings = (rand.Float64() + 0.5) * defaultParameters.Servings
	}

	//
	//generate all ingrediants first
	//
	//temp variables. this makes things easier
	var ndbNo string
	var measurement units.Quantity
	var data map[string]string
	var nutrition map[string]units.Quantity
	var totalServings float64
	var servings units.Quantity
	var price float64
	for (totalServings < options.Servings && r.Price < options.Dollars) || len(r.Steps) == 0 {
		//
		//generate random food
		//
		ndbNo, err = randomNdbNo()
		if err != nil {
			return
		}
		//generate random quantity
		servings = units.Quantity{Amount: (rand.Float64()*1.25 + 0.75), Unit: "servings"}
		//gather all necessary info
		measurement, data, nutrition, err = searchForFood(ndbNo, servings)
		if err != nil {
			return
		}
		//build actual ingrediant
		ingrediant := Ingrediant{}
		//fill in details
		ingrediant.ID = len(r.Steps)
		ingrediant.Name = data["Com_Desc"]
		ingrediant.Data = data
		ingrediant.Measurement = measurement
		ingrediant.Notes = "~random~"
		//
		//update reipce
		//
		//update servings
		totalServings += nutrition["servings"].Amount
		price, err = strconv.ParseFloat(data["price"], 64)
		if err != nil {
			price = 0
		}
		r.Price += price
		//update nutrition
		for k, v := range nutrition {
			_, ok := r.Nutrition[k]
			if ok {
				q := r.Nutrition[k]
				r.Nutrition[k] = units.Quantity{Unit: v.Unit, Amount: q.Amount + v.Amount, Type: v.Type}
			} else {
				r.Nutrition[k] = v
			}
		}
		//add into recipe
		r.Steps = append(r.Steps, &ingrediant)
	}

	//
	//build graph structure to keep track of to-do list
	//
	stuffToDo := make([]Step, len(r.Steps))
	for i, s := range r.Steps {
		stuffToDo[i] = s
	}

	//
	//combine recipes as necessary
	//
	for len(stuffToDo) > 1 {
		//randomize steps
		stuffToDo = shuffleSteps(stuffToDo)

		//TODO FIX THIS!!!

		//randomly select and build an operation
		preOp := operationList[rand.Intn(len(operationList))]
		op := Operation{}
		op.Name = "~Name~"
		op.ID = len(r.Steps)
		op.Operation = preOp.Name
		op.Notes = preOp.Notes
		op.Time = units.Quantity{Amount: rand.Float64()*(preOp.ETA[1]-preOp.ETA[0]) + preOp.ETA[0], Unit: "minute", Type: "time"}

		//look for possible prereqs
		//TODO: FINISH THIS!
		//
		//op.Requires????
		//

		//if can't find prereqs, just mix
		preOp = operationList[0]
		op = Operation{}
		op.Name = "~Name~"
		op.ID = len(r.Steps)
		op.Operation = preOp.Name
		op.Notes = preOp.Notes
		op.Time = units.Quantity{Amount: rand.Float64()*(preOp.ETA[1]-preOp.ETA[0]) + preOp.ETA[0], Unit: "minute", Type: "time"}
		//grab a few prereqs
		n := rand.Intn(len(stuffToDo)-1) + 2 //at least 2 things
		op.Requires = make([]int, n)
		for i, v := range stuffToDo[:n] {
			op.Requires[i] = v.getID()
		}
		sort.Ints(op.Requires)
		r.Steps = append(r.Steps, &op)
		stuffToDo = append(stuffToDo[n:], &op)
	}
	return
}

func shuffleSteps(src []Step) []Step {
	dest := make([]Step, len(src))
	perm := rand.Perm(len(src))
	for i, v := range perm {
		dest[v] = src[i]
	}
	return dest
}
