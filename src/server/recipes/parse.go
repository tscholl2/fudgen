package recipes

import (
	"errors"
	"gopkg.in/yaml.v2"
	"regexp"
	"strconv"
	"strings"
)

/*
see recipes/r.yml
*/

type recipe struct {
	Name        string   //name of food/recipe/step
	Operation   string   //name of operation to make this step
	Notes       string   //random notes to keep track of
	Attributes  []string //list of things like quantity or time
	Ingrediants []recipe //if empty then this is raw ingrediant
}

type Step struct {
	Name        string            `json:"name"`
	Operation   string            `json:"op"`
	Attributes  []string          `json:"attr"`
	Data        map[string]string `json:"data"`
	Notes       string            `json:"note"`
	Identifier  int               `json:"id"`
	Depedencies []int             `json:"deps"`
}

func indexOfRecipe(arr *[]recipe, ptr *recipe) (i int) {
	for i := 0; i < len(*arr); i++ {
		if &((*arr)[i]) == ptr {
			return i
		}
	}
	return -1
}

//returns time in seconds!
func getTime(attr []string) int {
	times := map[string]int{
		"days": 3600 * 24,
		"d":    3600 * 24,
		"h":    3600,
		"hr":   3600,
		"hour": 3600,
		"m":    60,
		"min":  60,
		"s":    1,
		"sec":  1}
	re, err := regexp.Compile(`[\d]+`)
	if err != nil {
		panic(err) //will never happen
	}
	for _, s := range attr {
		for k, v := range times {
			if strings.Index(s, k) != -1 {
				n := re.FindString(s)
				if n != "" {
					t, err := strconv.Atoi(n)
					if err != nil {
						panic(err)
					}
					return t * v
				}
			}
		}
	}
	return 0
}

func getQuantitiy(step Step) (q float32) {
	measurements := map[string]string{
		"cup":        "cup",
		"can":        "can",
		"jar":        "jar",
		"package":    "package",
		"ounce":      "ounce",
		"oz":         "ounce",
		"pound":      "pound",
		"whole":      "whole",
		"tablespoon": "tbl",
		"teaspoon":   "tsp",
		"pinch":      "pinch",
		"bunch":      "bunch"}
	for _, s := range step.Attributes {
		for k, _ := range measurements {
			if n := strings.Index(s, k); n != -1 {
				x, err := strconv.ParseFloat(s[:n], 32)
				if err != nil {
					panic(err)
				}
				y, err := strconv.ParseFloat(step.Data["Amount"], 32)
				if err != nil {
					panic(err)
				}
				return float32(x) * float32(y)
			}
		}
	}
	return 1
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
