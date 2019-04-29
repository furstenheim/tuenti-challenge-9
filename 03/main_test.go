package main

import (
"testing"
"github.com/stretchr/testify/assert"
)

func TestCaseSolve (t *testing.T) {
	testCases := []struct{
		c Case
		expected []Punch
	}{
		{
			c: Case{
				folds: []Fold{Top},
				x: 4,
				y: 2,
				punches: []Punch{Punch{0, 1}},
			},
			expected: []Punch{{0, 0}, {0, 3}},
		},
		{
			c: Case{
				folds: []Fold{Top, Top},
				x: 4,
				y: 2,
				punches: []Punch{Punch{0, 1}},
			},
			expected: []Punch{{0, 0}, {0, 3}, {0, 4}, {0, 7}},
		},
	}
	for _, tc := range(testCases) {
		punches := tc.c.Solve()
		assert.Equal(t, tc.expected, punches)
	}
}