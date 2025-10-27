package tyozetsu

import (
	"testing"
)

func TestNumberConvert(t *testing.T) {
	fizzBuzz := &NumberConverter{}
	var tests = []struct {
		n int
		want string
	}{
		{1, "1"},
		{2, "2"},
		{3, "Fizz"},
		{4, "4"},
		{5, "Buzz"},
		{6, "Fizz"},
		{10, "Buzz"},
		{15, "FizzBuzz"},
		{30, "FizzBuzz"},
	}

	for _, tt := range tests {
		t.Run("convert", func(t *testing.T) {
			got := fizzBuzz.convert(tt.n)
			if got != tt.want {
				t.Errorf(got, tt.want)
			}
			
		})
	}
}