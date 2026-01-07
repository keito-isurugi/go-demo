package domain

import "errors"

type Currency string

const (
	JPY Currency = "JPY"
	USD Currency = "USD"
)

type Money struct {
	amount   int
	currency Currency
}

func NewMoney(amount int, currency Currency) (Money, error) {
	if amount < 0 {
		return Money{}, errors.New("amount must be non-negative")
	}

	return Money{
		amount:   amount,
		currency: currency,
	}, nil
}

func (m Money) Add(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, errors.New("currency mismatch")
	}

	return NewMoney(m.amount+other.amount, m.currency)
}

func (m Money) Subtract(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, errors.New("currency mismatch")
	}

	newMoney := m.amount - other.amount
	if newMoney < 0 {
		return Money{}, errors.New("result cannot be negative")
	}

	return NewMoney(newMoney, m.currency)
}

func (m Money) Equals(other Money) bool {
	return m.amount == other.amount && m.currency == other.currency
}

func (m Money) IsGreaterThan(other Money) (bool, error) {
	if m.currency != other.currency {
		return false, errors.New("currency mismatch")
	}

	return m.amount > other.amount, nil
}

func (m Money) IsLessThan(other Money) (bool, error) {
	if m.currency != other.currency {
		return false, errors.New("currency mismatch")
	}

	return m.amount < other.amount, nil
}

func (m Money) Amount() int {
	return m.amount
}

func (m Money) Currency() Currency {
	return m.currency
}
