package rparser

import (
	"gopkg.in/yaml.v2"
)

/*
name: spaghetti
recipe:
  -
    name: mix
    attributes:
       - .5 min
    ingrediants:
      -
        name: pasta
        attributes:
          - .1 box
      -
        name: sauce
        attributes:
          - .15 oz
*/

type operation struct {
	Name        string
	Attributes  []string
	Ingrediants []operation
}

type recipe struct {
	Name   string
	Recipe []operation
}

func Parse(input string) (r recipe, err error) {
	err = yaml.Unmarshal([]byte(input), &r)
	return r, err
}
