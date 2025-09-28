package utils

import (
	"errors"
	"fmt"
	"strconv"
)

const (
	min int = 1
	max int = 100
)

func FizzBuzz(input int) (string, error) {
	if input < min || input > max {
		return "", errors.New(fmt.Sprintf("input is not between %d and %d", min, max))
	}
	
	if input % 15 == 0 {
		return "FizzBuzz", nil
	}

	if input % 3 == 0 {
		return "Fizz", nil
	}

	if input % 5 == 0 {
		return "Buzz", nil
	}
	return strconv.Itoa(input), nil
}