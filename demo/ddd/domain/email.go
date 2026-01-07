package domain

import (
	"errors"
	"strings"
)

type Email struct {
	value string
}

func NewEmail(value string) (Email, error) {
	// 正規化: 前後の空白を除去し、小文字に変換
	normalized := strings.ToLower(strings.TrimSpace(value))

	// 形式バリデーション
	if normalized == "" {
		return Email{}, errors.New("email cannot be empty")
	}

	atIndex := strings.Index(normalized, "@")
	if atIndex == -1 {
		return Email{}, errors.New("email must contain @")
	}

	if atIndex == 0 {
		return Email{}, errors.New("email local part cannot be empty")
	}

	domain := normalized[atIndex+1:]
	if domain == "" {
		return Email{}, errors.New("email domain cannot be empty")
	}

	if !strings.Contains(domain, ".") {
		return Email{}, errors.New("email domain must contain a dot")
	}

	return Email{value: normalized}, nil
}

func (e Email) String() string {
	return e.value
}

func (e Email) Domain() string {
	atIndex := strings.Index(e.value, "@")
	return e.value[atIndex+1:]
}

func (e Email) Equals(other Email) bool {
	return e.value == other.value
}
