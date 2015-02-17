package recipes

import (
	"bufio"
	"errors"
	"gopkg.in/yaml.v2"
	"os"
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
	Name        string              `json:"name"`
	Id          int                 `json:"id"`
	Nutrition   map[string]Quantity `json:"nutr"`
	Measurement Quantity            `json:"quant"`
	Notes       string              `json:"notes"`
}

type Step struct {
	Ingrediant Ingrediant
	Operation  Operation
}

type PreRecipe struct {
	Name        string      //name of food/recipe/step
	Operation   string      //name of operation to make this step
	Notes       string      //random notes to keep track of
	Attributes  []string    //list of things like quantity or time
	Ingrediants []PreRecipe //if empty then this is raw ingrediant
}

func (r *PreRecipe) getTime() Quantity {
	if len((*r).Ingrediants) > 0 { //if raw foods
		return Quantity{Amount: 0, Unit: "s"}
	} else {
		return Quantity{
			Amount: (*s).Operation.Time.Amount * times[(*s).Operation.Time.Unit], //will crash if time unit not in the table!
			Unit:   "s"}
	}
}
func (r *PreRecipe) getQuantity() Quantity {
	//TODO
}

var times map[string]int
var measurements map[string]string

func init() {
	//all acceptable times
	times = map[string]int{
		"days": 3600 * 24,
		"d":    3600 * 24,
		"h":    3600,
		"hr":   3600,
		"hour": 3600,
		"m":    60,
		"min":  60,
		"s":    1,
		"sec":  1}
	//load units from file
	var measurements []string
	file, err := os.Open("../../data/units.txt")
	if err != nil {
		panic(err)
	}
	reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		measurements = append(measurements, scanner.Text())
	}
}

func indexOfRecipe(arr *[]recipe, ptr *recipe) (i int) {
	for i := 0; i < len(*arr); i++ {
		if &((*arr)[i]) == ptr {
			return i
		}
	}
	return -1
}

func ParseYaml(input string) (steps []Step, err error) {
	var r PreRecipe
	err = yaml.Unmarshal([]byte(input), &r)

	//go through recipe collect steps
	//and then
	rs := []*PreRecipe{}
	var check func(*PreRecipe)
	check = func(R *PreRecipe) {
		rs = append(rs, R)
		if len((*R).Ingrediants) == 0 && (*R).Operation != "" {
			err = errors.New("can't have operation description on raw ingrediant")
			return
		}
		if len((*R).Ingrediants) > 0 {
			for k := 0; k < len((*R).Ingrediants); k++ {
				check(&((*R).Ingrediants[k]))
			}
		}
	}
	check(&r)
	if err != nil {
		return
	}

	for i := 0; i < len(rs); i++ {
		var rec PreRecipe
		rec = *(rs[i])
		s := Step{}
		s.Name = rec.Name
		s.Operation = rec.Operation
		s.Attributes = rec.Attributes
		s.Notes = rec.Notes
		s.Identifier = i
		s.Depedencies = []int{}
		for j := 0; j < len(rec.Ingrediants); j++ {
			s.Depedencies = append(s.Depedencies, indexOfRecipe(&(rec.Ingrediants), &(rec.Ingrediants[j])))
		}
		//if raw ingrediant, look up info
		if len(s.Depedencies) == 0 {
			s.Data, err = findFood(s.Name)
			if err != nil {
				return
			}
		}
		steps = append(steps, s)
	}
	return
}
