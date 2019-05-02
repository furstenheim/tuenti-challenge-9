package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestCaseIsKanjiPossible (t *testing.T) {
	testCases := []struct{
		n KanjiNumber
		m FrequencyMap
		probeRune rune
		expected bool
	}{
		{
			n: KanjiNumber{'二'},
			m: FrequencyMap{'二': 1},
			probeRune: '二',
			expected: true,
		},
		{
			n: KanjiNumber{'二'},
			m: FrequencyMap{'二': 2},
			probeRune: '二',
			expected: false,
		},
		{
			n: KanjiNumber{'十', '二'},
			m: FrequencyMap{'二': 2},
			probeRune: '二',
			expected: false,
		},
		{
			n: KanjiNumber{'二', '十', '二'},
			m: FrequencyMap{'二': 2, '十': 1},
			probeRune: '二',
			expected: true,
		},
	}
	for _, tc := range(testCases) {
		isPossible := tc.m.isKanjiPossible(tc.n)
		assert.Equal(t, tc.expected, isPossible)
		assert.GreaterOrEqual(t, tc.m[tc.probeRune], 1) // check we cloned correctly
	}
}
