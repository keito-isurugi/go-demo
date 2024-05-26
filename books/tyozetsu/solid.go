package tyozetsu

// import (
// 	"fmt"
// 	"strconv"
// )

// func FizzBuzz(i int) {
// 	if i % 15 == 0 {
// 		fmt.Println("FizzBuzz")
// 	} else if i % 3 == 0 {
// 		fmt.Println("Fizz")
// 	} else if i % 5 == 0 {
// 		fmt.Println("Buzz")
// 	} else {
// 		fmt.Println(i)
// 	}
// }

// type ReplaceRule interface {
// 	match(carry string, n int) bool
// 	apply(carry string, n int) string
// }

// type NumberConverter struct {
// 	Rules []ReplaceRule
// }
// func NewNumberConverter(
// 	replaceRules []ReplaceRule,
// ) *NumberConverter {
// 	return &NumberConverter{Rules: replaceRules}
// }
// func (nc *NumberConverter) convert(n int) string {
// 	carry := ""
// 	for _, r := range nc.Rules {
// 		if r.match(carry, n) {
// 			carry = r.apply(carry, n)
// 		}
// 	}
// 	return carry
// }

// type CyclicNumberRule struct {
// 	Base int
// 	Replacement string
// }
// func (cnr *CyclicNumberRule) match(carry string, n int) bool {
// 	return n % cnr.Base == 0
// }
// func (cnr *CyclicNumberRule) apply(carry string, n int) string {
// 	return carry + cnr.Replacement
// }

// type PassThroughRule struct {}
// func (cnr *PassThroughRule) match(carry string, n int) bool {
// 	return carry == ""
// }
// func (cnr *PassThroughRule) apply(carry string, n int) string {
// 	return strconv.Itoa(n)
// }

// func FizzBuzz2() {
// 	fizzbuzz := NewNumberConverter([]ReplaceRule{
// 		&CyclicNumberRule{Base: 3, Replacement: "Buzz"},
// 		&CyclicNumberRule{Base: 5, Replacement: "Fizz"},
// 		&PassThroughRule{},
// 	})

// 	fmt.Println(fizzbuzz.convert(1))
// 	fmt.Println(fizzbuzz.convert(3))
// 	fmt.Println(fizzbuzz.convert(5))
// 	fmt.Println(fizzbuzz.convert(15))
// }