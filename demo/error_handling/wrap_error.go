package main

import (
	"errors"
	"fmt"
)

// ラップしたエラーを返す
func wrapError(originalError error) error {
	return fmt.Errorf("wraped error: %w", originalError)
}

// 独自のエラー型
type MyError struct {
	Message string
}
func (e *MyError) Error() string {
	return e.Message
}

func execWrapError() {
	// 独自のエラー型をswitch文で分岐する
	// err := &MyError{Message: "my error"}
	// switch err.(type) {
	// case *MyError:
	// 	fmt.Println("my error")
	// default:
	// 	fmt.Println("unknown error")
	// }

	// errors.Asを使ってエラー型を判定する
	err := &MyError{Message: "my error"}
	var myErr *MyError
	if errors.As(err, &myErr) {
		fmt.Println("errors.As my error")
	}

	// errors.Isを使ってエラー型を判定する
	if errors.Is(err, &MyError{}) {
		fmt.Println("errors.Is my error")
	}
}
