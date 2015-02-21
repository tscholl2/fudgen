package recipes

import (
	"errors"
	"gopkg.in/yaml.v2"
	"strconv"
	"strings"
)

/*
see recipes/r.yml
*/

type Quantity struct {
	Unit   string  `json:"unit"`
	Amount float64 `json:"amt"`
}

type Operation struct {
	Name     string   `json:"name"`
	Id       int      `json:"id"`
	Time     Quantity `json:"time"`
	Requires []int    `json:inputs`
	Notes    string   `json:"notes"`
}

type Ingrediant struct {
	Name        string            `json:"name"`
	Id          int               `json:"id"`
	Data        map[string]string `json:data`
	Measurement Quantity          `json:"quant"`
	Notes       string            `json:"notes"`
}

type Step struct {
	Ingrediant Ingrediant
	Operation  Operation
}

type Recipe struct {
	Steps     []*Step             `json:"steps"`
	Title     string              `json:"title"`
	Nutrition map[string]Quantity `json:"nutr"`
}

func (s *Step) isOperation() bool {
	return (*s).Operation.Name != ""
}
func (s *Step) isIngrediant() bool {
	return !(*s).isOperation()
}
func (s *Step) getId() int {
	if (*s).isIngrediant() {
		return (*s).Ingrediant.Id
	} else {
		return (*s).Operation.Id
	}
}
func (s *Step) getTimeInSeconds() float64 {
	if (*s).isOperation() {
		return (*s).Operation.Time.Amount * float64(times[(*s).Operation.Time.Unit])
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
	Servings    float64     //number of servings of ingrediant
	Id          int         //for keeping track
	Ingrediants []PreRecipe //if empty then this is raw ingrediant
}

var times map[string]int

//given a string like "3 days"
//returns the number of seconds
func str2sec(s string) float64 {
	//remove extra space
	s = strings.TrimSpace(s)
	for k, v := range times {
		//look for units
		n := strings.Index(s, k)
		//if this is the right unit
		if n >= 0 && s[n:] == k {
			i, err := strconv.ParseFloat(strings.TrimSpace(s[:n]), 64)
			if err != nil {
				panic(err)
			}
			return i * float64(v)
		}
	}
	return -1
}

//var measurements map[string]string

func init() {
	//all acceptable times
	times = map[string]int{
		"d":       3600 * 24,
		"day":     3600 * 24,
		"days":    3600 * 24,
		"h":       3600,
		"hr":      3600,
		"hour":    3600,
		"hours":   3600,
		"m":       60,
		"min":     60,
		"minute":  60,
		"minutes": 60,
		"s":       1,
		"sec":     1,
		"seconds": 1}
}

//fills in steps by randomizing ingrediants
//calculating nutritional data
//and then returning the completed recipe
//with a random name
func steps2recipe(steps []*Step) (R Recipe, err error) {

	//copy steps into recipe
	for i := 0; i < len(steps); i++ {
		var s Step
		s = (*(steps[i])).copy()
		R.Steps = append(R.Steps, &s)
	}

	//initialize nutrition map
	R.Nutrition = make(map[string]Quantity)

	//keep track of names for title creation
	names := []string{}

	//find ingrediants and fill in
	for i := 0; i < len(R.Steps) && err == nil; i++ {
		s := R.Steps[i]
		if (*s).isIngrediant() {
			//convert measurement to # of servings
			servings := (*s).Ingrediant.Measurement.Amount

			measurement, data, nutrition, err := searchForFood((*s).Ingrediant.Name, servings)
			if err != nil {
				break
			}

			//add nutrition to totals
			for k, v := range nutrition {
				_, ok := R.Nutrition[k]
				if ok {
					q := R.Nutrition[k]
					R.Nutrition[k] = Quantity{Unit: v.Unit, Amount: q.Amount + v.Amount}
				} else {
					R.Nutrition[k] = v
				}

			}
			//add name to list
			names = append(names, (*s).Ingrediant.Name)
			//set measurement
			(*s).Ingrediant.Measurement = measurement
			(*s).Ingrediant.Data = data
		}
	}
	if err != nil {
		return
	}

	//build a title
	R.Title = randomTitle(names)

	return
}

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
		//check for terrible things
		if len((*R).Ingrediants) == 0 && (*R).Operation != "" {
			err = errors.New("can't have operation description on raw ingrediant")
			return
		}
		if (*R).Name == "" {
			err = errors.New("No name!")
			return
		}

		//convert to step
		s := Step{}
		if len((*R).Ingrediants) == 0 {
			// ---- for ingrediants
			s.Ingrediant.Id = (*R).Id
			s.Ingrediant.Name = (*R).Name
			s.Ingrediant.Notes = (*R).Notes
			s.Ingrediant.Measurement = Quantity{Amount: (*R).Servings, Unit: "servings"}
		} else {
			// ---- for operations
			s.Operation.Id = (*R).Id
			s.Operation.Name = (*R).Name
			s.Operation.Notes = (*R).Notes
			s.Operation.Time = Quantity{Amount: str2sec((*R).Time), Unit: "seconds"}
			for k := 0; k < len((*R).Ingrediants); k++ {
				s.Operation.Requires = append(s.Operation.Requires, (*R).Ingrediants[k].Id)
			}
		}
		steps = append(steps, &s)

		//recurse into dependencies
		for k := 0; k < len((*R).Ingrediants); k++ {
			check(&((*R).Ingrediants[k]))
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
