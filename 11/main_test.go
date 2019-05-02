package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"math"
)

func TestAngleAt (t *testing.T) {
	testCases := []struct{
		moon Moon
		t float64
		expected float64
	}{
		{
			moon: Moon{
				radius: 2.0,
				angle0: 0,
				periodHours: 12.0,
			},
			t: 0,
			expected: 0,

		},
		{
			moon: Moon{
				radius: 2.0,
				angle0: 0,
				periodHours: 12.0,
			},
			t: 1,
			expected: math.Pi,

		},
		{
			moon: Moon{
				radius: 2.0,
				angle0: math.Pi / 2,
				periodHours: 8,
			},
			t: 1,
			expected: 2 * math.Pi,

		},
	}
	for _, tc := range(testCases) {
		expected := tc.moon.currentAngle(tc.t)
		assert.Equal(t, tc.expected, expected)
	}
}

