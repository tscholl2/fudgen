package units

import (
	"testing"
)

func check(t *testing.T, value interface{}, expected interface{}) {
	if value != expected {
		t.Error("Expected ", expected, ", got ", value)
	}
}

func TestNumber(t *testing.T) {
	var s string
	var x float64
	var err error
	s = "7"
	x, err = ParseNumber(s)
	check(t, err, nil)
	check(t, x, float64(7))
	s = "134.99"
	x, err = ParseNumber(s)
	check(t, err, nil)
	check(t, x, float64(134.99))
	s = "3/2"
	x, err = ParseNumber(s)
	check(t, err, nil)
	check(t, x, float64(1.5))
}

func TestUnit(t *testing.T) {
	var s string
	var u string
	var err error
	s = "cup"
	u, err = ParseUnit(s)
	check(t, err, nil)
	check(t, u, "cup")
}

func TestParse(t *testing.T) {
	var s string
	var q Quantity
	var err error
	s = "  2.0  slices   "
	q, err = Parse(s)
	check(t, err, nil)
	check(t, q.Amount, float64(2))
	check(t, q.Type, "")
	check(t, q.Unit, "")
	s = "  2.0  pieces of something   "
	q, err = Parse(s)
	check(t, err, nil)
	check(t, q.Amount, float64(2))
	check(t, q.Type, "")
	check(t, q.Unit, "")
	s = "3/2 cups"
	q, err = Parse(s)
	check(t, err, nil)
	check(t, q.Amount, 1.5)
	check(t, q.Type, "volume")
	check(t, q.Unit, "cup")
	s = "7.8 hrs"
	q, err = Parse(s)
	check(t, err, nil)
	check(t, q.Amount, 7.8)
	check(t, q.Type, "time")
	check(t, q.Unit, "hour")
	q1 := q.ToBasic()
	check(t, err, nil)
	check(t, q1.Amount, 7.8*3600)
	check(t, q1.Type, "time")
	check(t, q1.Unit, "second")
}
