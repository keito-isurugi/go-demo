package demo

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAddMonthsPreservingEndOfMonth(t *testing.T) {
	testCases := []struct {
		name     string
		input    time.Time
		months   int
		expected time.Time
	}{
		{
			name:     "月を加算/通常のケース",
			input:    time.Date(2023, 7, 15, 0, 0, 0, 0, time.UTC),
			months:   0,
			expected: time.Date(2023, 7, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "月を加算/月末付近のケース/月末30日の月 -> 月末30日の月",
			input:    time.Date(2023, 4, 30, 0, 0, 0, 0, time.UTC),
			months:   2,
			expected: time.Date(2023, 6, 30, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "月を加算/月末付近のケース/月末30日の月 -> 月末31日の月",
			input:    time.Date(2023, 6, 30, 0, 0, 0, 0, time.UTC),
			months:   1,
			expected: time.Date(2023, 7, 31, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "月を加算/月末付近のケース/月末31日の月 -> 月末30日の月",
			input:    time.Date(2023, 5, 31, 0, 0, 0, 0, time.UTC),
			months:   1,
			expected: time.Date(2023, 6, 30, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "月を加算/月末付近のケース/月末31日の月 -> 月末31日の月",
			input:    time.Date(2023, 7, 31, 0, 0, 0, 0, time.UTC),
			months:   1,
			expected: time.Date(2023, 8, 31, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "月を加算/月末付近のケース/月末30日の月 -> 2月28日",
			input:    time.Date(2022, 11, 30, 0, 0, 0, 0, time.UTC),
			months:   3,
			expected: time.Date(2023, 2, 28, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "月を加算/月末付近のケース/月末30日の月 -> 2月29日(うるう年)",
			input:    time.Date(2019, 11, 30, 0, 0, 0, 0, time.UTC),
			months:   3,
			expected: time.Date(2020, 2, 29, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "月を加算/月末付近のケース/月末31日の月 -> 2月28日",
			input:    time.Date(2023, 1, 31, 0, 0, 0, 0, time.UTC),
			months:   1,
			expected: time.Date(2023, 2, 28, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "月を加算/月末付近のケース/月末31日の月 -> 2月29日(うるう年)",
			input:    time.Date(2020, 1, 31, 0, 0, 0, 0, time.UTC),
			months:   1,
			expected: time.Date(2020, 2, 29, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "月を加算/月末付近のケース/2月28日 -> 2月29日(うるう年)",
			input:    time.Date(2019, 2, 28, 0, 0, 0, 0, time.UTC),
			months:   12,
			expected: time.Date(2020, 2, 29, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "月を加算/月末付近のケース/2月29日(うるう年) -> 2月28日",
			input:    time.Date(2020, 2, 29, 0, 0, 0, 0, time.UTC),
			months:   12,
			expected: time.Date(2021, 2, 28, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "月を加算/年をまたぐケース",
			input:    time.Date(2023, 11, 30, 0, 0, 0, 0, time.UTC),
			months:   2,
			expected: time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "月を減算/月末付近のケース/月末30日の月 -> 月末30日の月",
			input:    time.Date(2023, 6, 30, 0, 0, 0, 0, time.UTC),
			months:   -2,
			expected: time.Date(2023, 4, 30, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "月を減算/月末付近のケース/月末30日の月 -> 月末31日の月",
			input:    time.Date(2023, 6, 30, 0, 0, 0, 0, time.UTC),
			months:   -1,
			expected: time.Date(2023, 5, 31, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "月を減算/月末付近のケース/月末31日の月 -> 月末30日の月",
			input:    time.Date(2023, 5, 31, 0, 0, 0, 0, time.UTC),
			months:   -1,
			expected: time.Date(2023, 4, 30, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "月を減算/月末付近のケース/月末31日の月 -> 月末31日の月",
			input:    time.Date(2023, 8, 31, 0, 0, 0, 0, time.UTC),
			months:   -1,
			expected: time.Date(2023, 7, 31, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "月を減算/月末付近のケース/月末30日の月 -> 2月28日",
			input:    time.Date(2023, 4, 30, 0, 0, 0, 0, time.UTC),
			months:   -2,
			expected: time.Date(2023, 2, 28, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "月を減算/月末付近のケース/月末30日の月 -> 2月29日(うるう年)",
			input:    time.Date(2020, 4, 30, 0, 0, 0, 0, time.UTC),
			months:   -2,
			expected: time.Date(2020, 2, 29, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "月を減算/月末付近のケース/月末31日の月 -> 2月28日",
			input:    time.Date(2023, 3, 31, 0, 0, 0, 0, time.UTC),
			months:   -1,
			expected: time.Date(2023, 2, 28, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "月を減算/月末付近のケース/月末31日の月 -> 2月29日(うるう年)",
			input:    time.Date(2020, 3, 31, 0, 0, 0, 0, time.UTC),
			months:   -1,
			expected: time.Date(2020, 2, 29, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "月を減算/月末付近のケース/2月28日 -> 2月29日(うるう年)",
			input:    time.Date(2021, 2, 28, 0, 0, 0, 0, time.UTC),
			months:   -12,
			expected: time.Date(2020, 2, 29, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "月を減算/月末付近のケース/2月29日(うるう年) -> 2月28日",
			input:    time.Date(2020, 2, 29, 0, 0, 0, 0, time.UTC),
			months:   -12,
			expected: time.Date(2019, 2, 28, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "月を減算/年をまたぐケース",
			input:    time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			months:   -1,
			expected: time.Date(2023, 12, 15, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := AddMonthsPreservingEndOfMonth(tc.input, tc.months)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
