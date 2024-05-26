package tyozetsu

import "strconv"

type NumberConverter struct {}
func (nc *NumberConverter) convert(n int) string{
	if n % 3 == 0 && n % 5 == 0 {
		return "FizzBuzz"
	}
	if n % 3 == 0 {
		return "Fizz"
	}
	if n % 5 == 0 {
		return "Buzz"
	}
	return strconv.Itoa(n)
}