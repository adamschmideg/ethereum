package main

import (
	"reflect"
	"testing"
)

var logcontent = `
bla bla
--- FAIL: TestFoo (0.22s)
other stuff
FAIL
ok      good/pkg
FAIL	bad/guy	11.00s
FAIL
`

func TestGetStats(t *testing.T) {
	//expected := stats{"good/pkg": map[string]int{"TestFoo": 1}}
	expected := Stats{}
	expected["bad/guy"] = map[string]int{"TestFoo": 1}
	got := GetStats(&logcontent)
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("GetStats expected: %#v\n\t got:%#v", expected, got)
	}
}

func TestCombineStats(t *testing.T) {
	var all = make(Stats)
	all["private"] = map[string]int{"TestMine": 1}
	all["foo"] = map[string]int{"TestThis": 2, "TestThat": 1}
	all["bar"] = map[string]int{"TestBeer": 1}
	var one = make(Stats)
	one["foo"] = map[string]int{"TestThis": 1}
	one["bar"] = map[string]int{"TestWine": 1}
	one["qux"] = map[string]int{"TestQux": 1}
	var expected = make(Stats)
	expected["private"] = map[string]int{"TestMine": 1}
	expected["foo"] = map[string]int{"TestThis": 3, "TestThat": 1}
	expected["bar"] = map[string]int{"TestBeer": 1, "TestWine": 1}
	expected["qux"] = map[string]int{"TestQux": 1}
	combineStats(&all, &one)
	for k, v := range all {
		expectedSubmap, ok := expected[k]
		if !ok {
			t.Error("Missing key", k)
		}
		if !reflect.DeepEqual(expectedSubmap, v) {
			t.Errorf("Mismatch in submaps for %v: %#v vs %#v", k, expectedSubmap, v)
		}
	}
}
