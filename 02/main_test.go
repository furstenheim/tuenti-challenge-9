package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestCaseSolve (t *testing.T) {
	testCases := []struct{
		c Case
		expected int
	}{
		{
			c: Case{
				nPlanets: 1,
				matrix: map[Planet][]Planet{
					"Galactica": []Planet{"New Earth"},
				},
			},
			expected: 1,
		},
		{
			c: Case{
				nPlanets: 2,
				matrix: map[Planet][]Planet{
					"Galactica": []Planet{"New Earth", "A"},
					"A": []Planet{"New Earth"},
				},
			},
			expected: 2,
		},
	}
	for _, tc := range(testCases) {
		nRoutes := tc.c.Solve()
		assert.Equal(t, nRoutes, tc.expected)
	}
}