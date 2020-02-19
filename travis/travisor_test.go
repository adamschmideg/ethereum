package main

import (
	"reflect"
	"testing"
)

var log = `
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
	got := GetStats(&log)
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("GetStats expected: %#v\n\t got:%#v", expected, got)
	}
}

