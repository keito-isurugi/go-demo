package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFizzBuzz(t *testing.T) {
	tests := []struct {
		input int
		expected string
		expectError bool
	}{
		{1, "1", false},
		{3, "Fizz", false},
		{6, "Fizz", false},
		{99, "Fizz", false},
		{5, "Buzz", false},
		{10, "Buzz", false},
		{25, "Buzz", false},
		{30, "FizzBuzz", false},
		{45, "FizzBuzz", false},
		{60, "FizzBuzz", false},
		{75, "FizzBuzz", false},
		{90, "FizzBuzz", false},
		{101, "", true},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%dを入力", tt.input), func(t *testing.T) {
			res, err := FizzBuzz(tt.input)
			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, err.Error(), fmt.Sprintf("input is not between %d and %d", min, max))
			} else {
				assert.NoError(t, err)
				assert.Equal(t, res, tt.expected)
			}
		})
	}
}