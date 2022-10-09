package handlers

import "testing"

func TestParsingId(t *testing.T) {
	type TestCase struct {
		name   string
		val    string
		expect int
	}
	testCases := []TestCase{{
		name:   "regular",
		val:    "/fact/1",
		expect: 1}, {
		name:   "zero",
		val:    "/fact/0",
		expect: 0}, {
		name:   "minus",
		val:    "/fact/-1",
		expect: 0}, {
		name:   "with additional letters",
		val:    "/fact/123sadasf",
		expect: 0}, {
		name:   "string instead of int",
		val:    "/fact/sadasf",
		expect: 0}, {
		name:   "longer URL",
		val:    "/fact/123/",
		expect: 0},
	}

	for _, v := range testCases {
		if returned, _ := parseID(v.val); returned != v.expect {
			t.Fatalf("test: %v, value: %v, expected %v got %v",
				v.name, v.val, v.expect, returned)
		}
	}

}
