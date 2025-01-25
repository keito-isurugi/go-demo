package main

import (
	"crypto/rand"
	"fmt"
)

func main() {
	o := NewOauth(rand.Read)
	h := NewHandler(o)
	res, err := h.GenerateCodeChallenges()
	if err != nil {
		panic(err)
	}

	fmt.Printf("State: %+v\n", res.State)
	fmt.Printf("CodeVerifier: %+v\n", res.CodeVerifier)
	fmt.Printf("CodeChallenge: %+v\n", res.CodeChallenge)
}
