package recipes

import (
	"../units"
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"strconv"
)

/*
see recipes/r.yml
*/

type Operation struct {
	Name        string         `json:"name"`
	Description string         `json:"desc"`
	Id          int            `json:"id"`
	Time        units.Quantity `json:"time"`
	Requires    []int          `json:inputs`
	Notes       string         `json:"notes"`
}

func (s *Operation) IsIngrediant() bool {
	return false
}
func (s *Operation) MarshalJSON() ([]byte, error) {
	return json.Marshal(s)
}

type Ingrediant struct {
	Name        string            `json:"name"`
	Id          int               `json:"id"`
	Data        map[string]string `json:data`
	Measurement units.Quantity    `json:"quant"`
	Notes       string            `json:"notes"`
}

func (s *Ingrediant) IsIngrediant() bool {
	return true
}
func (s *Ingrediant) MarshalJSON() ([]byte, error) {
	return json.Marshal(s)
}

type Recipe struct {
	Steps     []*Step                   `json:"steps"`
	Title     string                    `json:"title"`
	Nutrition map[string]units.Quantity `json:"nutr"`
	Price     float64                   `json:price`
}
type Step interface { //because I don't know how to "extend" objects
	IsIngrediant() bool
	MarshalJSON() ([]byte, error)
}

type preRecipe struct {
	Name        string      //name of food/recipe/step
	Operation   string      //name of operation to make this step, nil for ingrediants
	Notes       string      //random notes to keep track of
	Time        string      //length of step, nil for ingrediants
	Quantity    string      //how much of ingrediant, e.g. "1/2 cup" or "3 slices"
	Id          int         //for keeping track
	Ingrediants []preRecipe //if empty then this is raw ingrediant
}

func preRecipe2Recipe(PR *preRecipe) (R *Recipe, err error) {
	//go through recipe collect steps
	steps := []*Step{}
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
			pr2.Name = fmt.Sprintf("Step %d", pr2.Id)
		}

		//convert to step
		var s *Step
		if len(pr2.Ingrediants) == 0 {
			// ---- for ingrediants
			i := Ingrediant{}
			i.Id = pr2.Id
			i.Name = pr2.Name
			i.Notes = pr2.Notes
			i.Measurement, err = units.Parse(pr2.Quantity)
			s = &i
		} else {
			// ---- for operations
			o := Operation{}
			o.Id = pr2.Id
			o.Name = pr2.Name
			o.Description = pr2.Operation
			o.Notes = pr2.Notes
			o.Time, err = units.Parse(pr2.Time)
			for k := 0; k < len(R.Ingrediants); k++ {
				o.Requires = append(o.Requires, pr2.Ingrediants[k].Id)
			}
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
	check(&r)
	if err != nil {
		return
	}

	//fill in recipe automagically
	R, err = steps2recipe(steps)

	//finally return the new set of steps
	return
}

//fills in steps by randomizing ingrediants
//calculating nutritional data
//and then returning the completed recipe
//with a random name
func steps2recipe(steps []*Step) (R Recipe, err error) {

	//copy steps into recipe
	for i := 0; i < len(steps); i++ {
		R.Steps = append(R.Steps, s)
	}

	//initialize nutrition map
	R.Nutrition = make(map[string]units.Quantity)

	//keep track of names for title creation
	names := []string{}

	//find ingrediants and fill in
	for i := 0; i < len(R.Steps) && err == nil; i++ {
		s := R.Steps[i]
		if s.isIngrediant() {

			ing := s.(Ingrediant)

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
			names = append(names, i.Name)

			//set measurement
			ing.Measurement = measurement
			ing.Data = data

			//replace old step
			R[i] = &ing
		}
	}
	if err != nil {
		return
	}

	//build a title
	R.Title = randomTitle(names)

	return
}

//parses yaml into full recipe structure
//fills in as best as possible
func ParseYaml(input string) (R *Recipe, err error) {
	//parse yaml into pre-recipe structure
	var pr PreRecipe
	err = yaml.Unmarshal([]byte(input), &pr)
	if err != nil {
		return
	}
	R, err = preRecipe2Recipe(pr)
	return
}
