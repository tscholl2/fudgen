package recipes

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"../units"
	"gopkg.in/yaml.v2"
)

// Operation data structure
type Operation struct {
	Name        string         `json:"name"`
	Description string         `json:"desc"`
	ID          int            `json:"id"`
	Time        units.Quantity `json:"time"`
	Requires    []int          `json:"reqs"`
	Notes       string         `json:"notes"`
}

//IsIngrediant returns false for operations
func (s *Operation) IsIngrediant() bool {
	return false
}

//MarshalJSON returns this operation as json
func (s *Operation) MarshalJSON() ([]byte, error) {
	return json.Marshal(s)
}

//Ingrediant data structure
type Ingrediant struct {
	Name        string            `json:"name"`
	ID          int               `json:"id"`
	Data        map[string]string `json:"data"`
	Measurement units.Quantity    `json:"quant"`
	Notes       string            `json:"notes"`
}

//IsIngrediant returns true for ingrediants
func (s *Ingrediant) IsIngrediant() bool {
	return true
}

//MarshalJSON returns ingrediant as json
func (s *Ingrediant) MarshalJSON() ([]byte, error) {
	return json.Marshal(s)
}

//Recipe data struct, contains list of steps
//which are either operations or ingrediants
type Recipe struct {
	Steps     []Step                    `json:"steps"`
	Title     string                    `json:"title"`
	Nutrition map[string]units.Quantity `json:"nutr"`
	Price     float64                   `json:"price"`
}

//Step interface. just says what it can be cast to
type Step interface {
	IsIngrediant() bool
	MarshalJSON() ([]byte, error)
}

type preRecipe struct {
	Name        string      //name of food/recipe/step
	Operation   string      //name of operation to make this step, nil for ingrediants
	Notes       string      //random notes to keep track of
	Time        string      //length of step, nil for ingrediants
	Quantity    string      //how much of ingrediant, e.g. "1/2 cup" or "3 slices"
	ID          int         //for keeping track
	Ingrediants []preRecipe //if empty then this is raw ingrediant
}

func preRecipe2Steps(pr *preRecipe) (steps []Step, err error) {

	//and then convert to actual recipe structure
	var check func(*preRecipe)
	check = func(pr2 *preRecipe) {
		if err != nil {
			return
		}

		//check for terrible things
		if len(pr2.Ingrediants) == 0 && pr2.Operation != "" {
			err = errors.New("can't have operation description on raw ingrediant")
			return
		}
		if pr2.Name == "" {
			pr2.Name = fmt.Sprintf("Step %d", pr2.ID)
		}

		//convert to step
		var s Step

		if len(pr2.Ingrediants) == 0 {
			// ---- for ingrediants
			i := Ingrediant{}
			i.ID = pr2.ID
			i.Name = pr2.Name
			i.Notes = pr2.Notes
			i.Measurement, err = units.Parse(pr2.Quantity)
			s = &i
		} else {
			// ---- for operations
			o := Operation{}
			o.ID = pr2.ID
			o.Name = pr2.Name
			o.Description = pr2.Operation
			o.Notes = pr2.Notes
			o.Time, err = units.Parse(pr2.Time)
			for k := 0; k < len(pr2.Ingrediants); k++ {
				o.Requires = append(o.Requires, pr2.Ingrediants[k].ID)
			}
			s = &o
		}
		if err != nil {
			return
		}
		steps = append(steps, s)

		//recurse into dependencies
		for k := 0; k < len(pr2.Ingrediants); k++ {
			check(&(pr2.Ingrediants[k]))
		}
	}
	check(pr)

	//finally return the new set of steps
	return
}

//fills in steps by randomizing ingrediants
//calculating nutritional data
//and then returning the completed recipe
//with a random name
func steps2recipe(steps []Step) (R *Recipe, err error) {
	//initialize output
	r := Recipe{}
	R = &r

	//copy steps into recipe
	R.Steps = make([]Step, len(steps))
	copy(R.Steps, steps)

	//initialize nutrition map
	R.Nutrition = make(map[string]units.Quantity)

	//keep track of names for title creation
	names := []string{}

	//find ingrediants and fill in
	for i := 0; i < len(R.Steps) && err == nil; i++ {
		s := R.Steps[i]
		if s.IsIngrediant() {

			ing := s.(*Ingrediant)

			//look for closest/slightly random food
			measurement, data, nutrition, err := searchForFood(ing.Name, ing.Measurement)
			if err != nil {
				break
			}

			//add nutrition to totals
			for k, v := range nutrition {
				_, ok := R.Nutrition[k]
				if ok {
					q := R.Nutrition[k]
					R.Nutrition[k] = units.Quantity{Unit: v.Unit, Amount: q.Amount + v.Amount, Type: v.Type}
				} else {
					R.Nutrition[k] = v
				}
			}

			//add price to total if possible
			_, ok := data["price"]
			if ok {
				x, _ := strconv.ParseFloat(data["price"], 64)
				R.Price += x
			}

			//add name to list
			for _, n := range strings.Fields(ing.Name) {
				names = append(names, n)
			}

			//set measurement
			ing.Measurement = measurement
			ing.Data = data

			//replace old step
			R.Steps[i] = ing
		} else {
			op := s.(*Operation)

			//add name to list
			for _, n := range strings.Fields(op.Name) {
				names = append(names, n)
			}
		}
	}
	if err != nil {
		return
	}

	//build a title
	R.Title = randomTitle(R)

	//return recipe
	return
}

// ParseYaml takes yaml and returns
// full recipe structure
// fills in as best as possible, automatically
// randomizes ingrediants
func ParseYaml(input string) (R *Recipe, err error) {
	//parse yaml into pre-recipe structure
	var pr preRecipe
	err = yaml.Unmarshal([]byte(input), &pr)
	if err != nil {
		return
	}
	//set ids
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
	setID(&pr)
	steps, err := preRecipe2Steps(&pr)
	if err != nil {
		return
	}
	R, err = steps2recipe(steps)
	return
}
