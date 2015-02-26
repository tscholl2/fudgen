package recipes

import (
	"../units"
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"strconv"
)

/*
see recipes/r.yml
*/

type Operation struct {
	Name     string         `json:"name"`
	Id       int            `json:"id"`
	Time     units.Quantity `json:"time"`
	Requires []int          `json:inputs`
	Notes    string         `json:"notes"`
}

type Ingrediant struct {
	Name        string            `json:"name"`
	Id          int               `json:"id"`
	Data        map[string]string `json:data`
	Measurement units.Quantity    `json:"quant"`
	Notes       string            `json:"notes"`
}

type Step struct { //because I don't know how to "extend" objects
	Ingrediant Ingrediant `json:"ingrediant"`
	Operation  Operation  `json:"op":`
}

type Recipe struct { //TODO fix the stupid steps thing
	Steps     []*Step                   `json:"steps"`
	Title     string                    `json:"title"`
	Nutrition map[string]units.Quantity `json:"nutr"`
	Price     float64                   `json:price`
}

func (s *Step) Name() string {
	if s.isIngrediant() {
		return s.Ingrediant.Name
	} else {
		return s.Operation.Name
	}
}
func (s *Step) isOperation() bool {
	return s.Operation.Name != ""
}
func (s *Step) isIngrediant() bool {
	return !s.isOperation()
}
func (s *Step) getId() int {
	if s.isIngrediant() {
		return s.Ingrediant.Id
	} else {
		return s.Operation.Id
	}
}
func (s *Step) getTimeInSeconds() float64 {
	if s.isOperation() {
		q := s.Operation.Time.ToBasic()
		return q.Amount
	} else {
		return 15
	}
}
func (s *Step) copy() (t Step) {
	if s.isIngrediant() {
		t.Ingrediant.Name = s.Ingrediant.Name
		t.Ingrediant.Notes = s.Ingrediant.Notes
		t.Ingrediant.Id = s.Ingrediant.Id
		t.Ingrediant.Measurement = s.Ingrediant.Measurement
		t.Ingrediant.Data = make(map[string]string)
		for k, v := range s.Ingrediant.Data {
			t.Ingrediant.Data[k] = v
		}
	} else {
		t.Operation.Id = s.Operation.Id
		t.Operation.Name = s.Operation.Name
		t.Operation.Notes = s.Operation.Notes
		t.Operation.Time = s.Operation.Time
		t.Operation.Requires = make([]int, len(s.Operation.Requires))
		for i := 0; i < len(s.Operation.Requires); i++ {
			t.Operation.Requires[i] = s.Operation.Requires[i]
		}
	}
	return
}

type PreRecipe struct {
	Name        string      //name of food/recipe/step
	Operation   string      //name of operation to make this step, nil for ingrediants
	Notes       string      //random notes to keep track of
	Time        string      //length of step, nil for ingrediants
	Quantity    string      //how much of ingrediant, e.g. "1/2 cup" or "3 slices"
	Id          int         //for keeping track
	Ingrediants []PreRecipe //if empty then this is raw ingrediant
}

//fills in steps by randomizing ingrediants
//calculating nutritional data
//and then returning the completed recipe
//with a random name
func steps2recipe(steps []*Step) (R Recipe, err error) {

	//copy steps into recipe
	for i := 0; i < len(steps); i++ {
		var s Step
		s = (steps[i]).copy()
		R.Steps = append(R.Steps, &s)
	}

	//initialize nutrition map
	R.Nutrition = make(map[string]units.Quantity)

	//keep track of names for title creation
	names := []string{}

	//find ingrediants and fill in
	for i := 0; i < len(R.Steps) && err == nil; i++ {
		s := R.Steps[i]
		if s.isIngrediant() {

			//look for closest/slightly random food
			measurement, data, nutrition, err := searchForFood(s.Ingrediant.Name, s.Ingrediant.Measurement)
			if err != nil {
				//panic(err) //TODO return here?
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
			names = append(names, s.Ingrediant.Name)
			//set measurement
			s.Ingrediant.Measurement = measurement
			s.Ingrediant.Data = data
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
func ParseYaml(input string) (R Recipe, err error) {
	//parse yaml into pre-recipe structure
	var r PreRecipe
	err = yaml.Unmarshal([]byte(input), &r)
	if err != nil {
		return
	}

	//first go through and set id's
	id := 0
	var setId func(*PreRecipe)
	setId = func(R *PreRecipe) {
		(*R).Id = id
		id += 1
		//recurse into dependencies
		for k := 0; k < len((*R).Ingrediants); k++ {
			setId(&((*R).Ingrediants[k]))
		}
	}
	setId(&r)

	//go through recipe collect steps
	steps := []*Step{}
	//and then convert to actual recipe structure
	var check func(*PreRecipe)
	check = func(R *PreRecipe) {
		if err != nil {
			return
		}
		//check for terrible things
		if len(R.Ingrediants) == 0 && R.Operation != "" {
			err = errors.New("can't have operation description on raw ingrediant")
			return
		}
		if R.Name == "" {
			R.Name = fmt.Sprintf("Step %d", R.Id)
		}

		//convert to step
		s := Step{}
		if len(R.Ingrediants) == 0 {
			// ---- for ingrediants
			s.Ingrediant.Id = R.Id
			s.Ingrediant.Name = R.Name
			s.Ingrediant.Notes = R.Notes
			s.Ingrediant.Measurement, err = units.Parse(R.Quantity)
		} else {
			// ---- for operations
			s.Operation.Id = R.Id
			s.Operation.Name = R.Name
			s.Operation.Notes = R.Notes
			s.Operation.Time, err = units.Parse(R.Time)
			for k := 0; k < len(R.Ingrediants); k++ {
				s.Operation.Requires = append(s.Operation.Requires, R.Ingrediants[k].Id)
			}
		}
		if err != nil {
			return
		}
		steps = append(steps, &s)
		//recurse into dependencies
		for k := 0; k < len(R.Ingrediants); k++ {
			check(&(R.Ingrediants[k]))
		}
	}
	check(&r)
	if err != nil {
		return
	}

	//fill in recipe automagically
	R, err = steps2recipe(steps)
	if err != nil {
		return
	}

	//finally return the new set of steps
	return
}
