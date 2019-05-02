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
			n: KanjiNumber{'二', '二'},
			m: FrequencyMap{'二': 1},
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
func TestNumberToKanji (t *testing.T) {
	testCases := []struct{
		n int
		expected KanjiNumber
	}{
		{
			n: 2,
			expected: KanjiNumber{'二'},
		},
		{
			n: 1,
			expected: KanjiNumber{'一'},
		},
		{
			n: 10,
			expected: KanjiNumber{'十'},
		},
		{
			n: 12,
			expected: KanjiNumber{'十', '二'},
		},
		{
			n: 121,
			expected: KanjiNumber{'百', '二', '十', '一'},
		},
		{
			n: 1345,
			expected: KanjiNumber{'千', '三', '百', '四', '十', '五'},
		},
		{
			n: 12345,
			expected: KanjiNumber{'一', '万', '二', '千', '三', '百', '四', '十', '五'},
		},

	}
	for _, tc := range(testCases) {
		expected := intToKanji(tc.n)
		assert.Equal(t, tc.expected, expected)
	}
}
