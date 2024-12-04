package main

import (
	"testing"
)

func TestCaesarCipher(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		shift  int
		output string
	}{
		{"LowercaseWithSmallShift", "abc", 1, "bcd"},
		{"UppercaseWithSmallShift", "ABC", 1, "BCD"},
		{"MixCaseAndWrapAround", "aBcXyZ", 3, "dEfAbC"},
		{"NegativeShift", "abc", -1, "zab"},
		{"NoShift", "abc", 0, "abc"},
		{"FullAlphabetShift", "z", 26, "z"},
		{"SpecialCharacters", "abc!@#", 1, "bcd!@#"},
		{"SpacesAndNumbers", "Hello World 123", 5, "Mjqqt Btwqi 123"},
		{"NegativeShiftWrap", "cde", -3, "zab"},
		{"EmptyString", "", 3, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := caesarCipher(tt.input, tt.shift)
			if result != tt.output {
				t.Errorf("caesarCipher(%q, %d) = %q, want %q", tt.input, tt.shift, result, tt.output)
			}
		})
	}
}
