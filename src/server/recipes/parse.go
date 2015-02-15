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
	Name   string   `json:"name"`
	Time   Quantity `json:"time"`
	Output string   `json:"output"`
	Inputs []int    `json:inputs`
	Notes  string   `json:"notes"`
}

type Ingrediant struct {
	Name        string              `json:"name"`
	Nutrition   map[string]Quantity `json:"nutr"`
	Measurement Quantity            `json:"quant"`
	Notes       string              `json:"notes"`
}

type recipe struct {
	Ingrediant Ingrediant
	Operation  Operation
	Requires   []recipe //if empty then this is raw ingrediant
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
	var r recipe
	err = yaml.Unmarshal([]byte(input), &r)

	//go through recipe collect steps
	//and then
	rs := []*recipe{}
	var check func(*recipe)
	check = func(R *recipe) {
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
		var rec recipe
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
