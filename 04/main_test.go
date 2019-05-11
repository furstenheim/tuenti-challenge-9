package main

import (
"testing"
"github.com/stretchr/testify/assert"
	"math/big"
)

func TestCaseSolve (t *testing.T) {
	testCases := []struct{
		c Case
		expectedN, expectedD *big.Int
	}{
	}
	for _, tc := range(testCases) {
		n, d := tc.c.Solve()
		assert.Equal(t, tc.expectedD, d)
		assert.Equal(t, tc.expectedN, n)
	}
}