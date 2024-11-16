package main

import "fmt"

// Ruleインターフェース
type Rule interface {
    Match(num int) bool
    String() string
}

// FizzRule構造体
type FizzRule struct{}

func (r *FizzRule) Match(num int) bool {
    return num%3 == 0
}

func (r *FizzRule) String() string {
    return "Fizz"
}

// BuzzRule構造体
type BuzzRule struct{}

func (r *BuzzRule) Match(num int) bool {
    return num%5 == 0
}

func (r *BuzzRule) String() string {
    return "Buzz"
}

// FizzBuzz関数
func FizzBuzz(num int, rules []Rule) string {
    for _, rule := range rules {
        if rule.Match(num) {
            return rule.String()
        }
    }
    return fmt.Sprintf("%d", num)
}

func main() {
    rules := []Rule{&FizzRule{}, &BuzzRule{}}

    for i := 1; i <= 100; i++ {
        fmt.Println(FizzBuzz(i, rules))
    }
}