package units

import (
	"errors"
	"fmt"
)

type Quantity struct {
	Unit   string  `json:"unit"`
	Amount float64 `json:"amt"`
	Type   string  `json:"-"`
}

//converts this quantity to a basic unit
//use with care!!! does NOT throw error
func (this *Quantity) ToBasic() (q Quantity) {
	var x float64
	var u string
	switch unit_type[this.Unit] {
	case "time":
		x = this.Amount * toSeconds[this.Unit] //convert to seconds
		u = "second"
	case "mass":
		x = this.Amount * toGrams[this.Unit] //convert to grams
		u = "gram"
	case "volume":
		x = this.Amount * toMilliters[this.Unit] //convert to milliters
		u = "milliter"
	default:
		x = q.Amount
	}
	q.Type = unit_type[u]
	q.Amount = x
	q.Unit = u
	return
}

//converts to given unit type
func (this *Quantity) Convert(u string) (q Quantity, err error) {
	//translate to actual unit
	u, err = ParseUnit(u)
	if err != nil {
		err = errors.New(fmt.Sprintf("unable to parse [%s] --- conversion failed", this))
		return
	}
	//check type
	t := unit_type[u]
	if t != this.Type {
		err = errors.New(fmt.Sprintf("unable to convert [%s] ---> [%s]", this.Type, t))
		return
	}
	//convert by pushing both to basic
	q1 := this.ToBasic()
	q = Quantity{Unit: u, Amount: 1, Type: t}
	q2 := q.ToBasic()
	x := q2.Amount
	y := q1.Amount
	q.Amount = x * 1.0 / y
	return
}

var toMilliters map[string]float64 //unit ---> milliliters
var toGrams map[string]float64     //unit ---> grams
var toSeconds map[string]float64   //unit ---> seconds
var unit_type map[string]string    //just a nice map of standardized time units
var alternatives map[string]string //takes common unit names to specified ones

//initalize all static variables
func init() {
	unit_type = map[string]string{}

	//these hold only standardized unit names!
	toSeconds = map[string]float64{
		"day":    86400,
		"hour":   3600,
		"minute": 60,
		"second": 1}

	for k, _ := range toSeconds {
		unit_type[k] = "time"
	}

	toMilliters = map[string]float64{
		"cup":              236.588,
		"tablespoon":       14.7868,
		"teaspoon":         4.92892,
		"gallon":           3785.41,
		"quart":            946.353,
		"pint":             473.176,
		"liter":            1000,
		"cubic centimeter": 1,
		"milliter":         1,
		"fluid ounce":      29.5735}

	for k, _ := range toMilliters {
		unit_type[k] = "volume"
	}

	toGrams = map[string]float64{
		"ounce":      28.3495,
		"pound":      453.592,
		"stone":      6350.29,
		"grain":      0.06479891,
		"metric ton": 1000000,
		"gram":       1,
		"microgram":  1.0 / 1000000,
		"kilogram":   1000}

	for k, _ := range toGrams {
		unit_type[k] = "mass"
	}

	alternatives = map[string]string{
		//volume
		"tbs":              "tablespoon",
		"tbl":              "tablespoon",
		"tbls":             "tablespoon",
		"tblsp":            "tablespoon",
		"tblsps":           "tablespoon",
		"tablespoons":      "tablespoon",
		"tablespoon":       "tablespoon",
		"tsp":              "teaspoon",
		"tsps":             "teaspoon",
		"teaspoons":        "teaspoon",
		"teaspoon":         "teaspoon",
		"cups":             "cup",
		"cup":              "cup",
		"gal":              "gallon",
		"gals":             "gallon",
		"gallons":          "gallon",
		"gallon":           "gallon",
		"qt":               "quart",
		"quart":            "quart",
		"pt":               "pint",
		"pint":             "pint",
		"l":                "liter",
		"litre":            "liter",
		"liter":            "liter",
		"cc":               "cubic centimeter",
		"centimeter cubed": "cubic centimeter",
		"cm^3":             "cubic centimeter",
		"centimeter^3":     "cubic centimeter",
		"cubic centimeter": "cubic centimeter",
		"floz":             "fluid ounce",
		"fl oz":            "fluid ounce",
		"fl ounce":         "fluid ounce",
		"fluid oz":         "fluid ounce",
		"fluid_ounce":      "fluid ounce",
		"US_fluid_ounce ":  "fluid ounce",
		"US_liquid_ounce":  "fluid ounce",
		"fluid ounce":      "fluid ounce",
		//weight
		"oz":          "ounce",
		"ozs":         "ounce",
		"ounces":      "ounce",
		"ounce":       "ounce",
		"lb":          "pound",
		"lbs":         "pound",
		"pounds":      "pound",
		"pound":       "pound",
		"stones":      "stone",
		"stone":       "stone",
		"grains":      "grain",
		"grain":       "grain",
		"ton":         "metric ton",
		"tons":        "metric ton",
		"metric tons": "metric ton",
		"metric ton":  "metric ton",
		"g":           "gram",
		"gm":          "gram",
		"gms":         "gram",
		"grams":       "gram",
		"gram":        "gram",
		"um":          "microgram",
		"µm":          "microgram",
		"ums":         "microgram",
		"µms":         "microgram",
		"micrograms":  "microgram",
		"microgram":   "microgram",
		"kg":          "kilogram",
		"kgs":         "kilogram",
		"kilograms":   "kilogram",
		"kilogram":    "kilogram",
		//time
		"d":       "day",
		"day":     "day",
		"days":    "day",
		"h":       "hour",
		"hr":      "hour",
		"hrs":     "hour",
		"hour":    "hour",
		"hours":   "hour",
		"m":       "minute",
		"min":     "minute",
		"minute":  "minute",
		"minutes": "minute",
		"s":       "second",
		"sec":     "second",
		"seconds": "second",
		"second":  "second"}
}
