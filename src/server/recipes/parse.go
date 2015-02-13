package recipes

import (
	"errors"
	"gopkg.in/yaml.v2"
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
			s.Data, err = FindFood(s.Name)
			if err != nil {
				return
			}
		}
		steps = append(steps, s)
	}
	return
}
