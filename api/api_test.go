package api

import (
	"testing"
)

var parseAsIntAndUintTests = map[string]struct {
	numResults string
	min int
	max	int
	def	int
	expected int
} {
	"parse error test case": {"not a parsable integer", 101, 1000, 5000, 5000},
	"less than min test case": {"100", 101, 1000, 5000, 101},
	"greater than max test case": {"1001", 101, 1000, 5000, 1000},
	"within limits test case": {"100", 0, 100, 5000, 100},
}

func TestParseAsIntAndUint(t *testing.T) {
	for name, td := range parseAsIntAndUintTests {
		t.Run(name, func(t *testing.T) {
			if got, expected := 
				parseAsUintValue(td.numResults, uint(td.min), uint(td.max), uint(td.def)), uint(td.expected); got != expected {
					t.Fatalf("parseAsUintValue - %s: returned %d; expected %d", name, got, expected)
				}

			if got, expected := 
				parseAsIntValue(td.numResults, td.min, td.max, td.def), td.expected; got != expected {
					t.Fatalf("parseAsIntValue - %s: returned %d; expected %d", name, got, expected)
				}
		})	
	}
}