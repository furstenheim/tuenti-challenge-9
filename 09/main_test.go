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

func TestToConfusedNumber (t *testing.T) {
	testCases := []struct{
		original KanjiNumber
		expected ConfusedNumber
	}{
		{
			original: KanjiNumber{'二'},
			expected: ConfusedNumber{
				powersOfTen: []PowerOfTen{},
				digits: []int{2},
				allowFirstDigit: true,
			},

		},
		{
			original: KanjiNumber{'七', '十', '二', '千'},
			expected: ConfusedNumber{
				powersOfTen: []PowerOfTen{powersOfTenByRune['十'], powersOfTenByRune['千']},
				digits: []int{7, 2},
				allowFirstDigit: true,
			},

		},		{
			original: KanjiNumber{'七', '十', '二', '二', '千'},
			expected: ConfusedNumber{
				powersOfTen: []PowerOfTen{powersOfTenByRune['十'], powersOfTenByRune['千']},
				digits: []int{7, 2, 2},
				allowFirstDigit: true,
			},

		},
		{
			original: KanjiNumber{'十'},
			expected: ConfusedNumber{
				powersOfTen: []PowerOfTen{powersOfTenByRune['十']},
				digits: []int{},
				allowFirstDigit: false,
			},

		},
	}
	for _, tc := range(testCases) {
		expected := tc.original.toConfusedNumber()
		assert.Equal(t, tc.expected, expected)
	}
}

func TestFindAllPossibleCombinations (t *testing.T) {
	testCases := []struct{
		expected []int
		confused ConfusedNumber
		digits []int
	}{
		{
			confused: ConfusedNumber{
				powersOfTen: []PowerOfTen{},
				digits: []int{2},
				allowFirstDigit: true,
			},
			digits: []int{2},
			expected: []int{2},

		},
		{
			confused: ConfusedNumber{
				powersOfTen: []PowerOfTen{powersOfTenByRune['十'], powersOfTenByRune['千']},
				digits: []int{7, 2},
				allowFirstDigit: true,
			},
			digits: []int{7, 2},
			expected: []int{2017, 1027, 2070},

		},
		{
			confused: ConfusedNumber{
				powersOfTen: []PowerOfTen{powersOfTenByRune['十'], powersOfTenByRune['千']},
				digits: []int{7, 2, 3},
				allowFirstDigit: true,
			},
			digits: []int{7, 2, 3},
			expected: []int{3027},

		},
		{
			confused: ConfusedNumber{
				powersOfTen: []PowerOfTen{powersOfTenByRune['十']},
				digits: []int{},
				allowFirstDigit: false,
			},
			digits: []int{},
			expected: []int{10},

		},
		{
			confused: ConfusedNumber{
				powersOfTen: []PowerOfTen{powersOfTenByRune['十'], powersOfTenByRune['万']},
				digits: []int{1},
				allowFirstDigit: true,
			},
			digits: []int{1},
			expected: []int{10010},

		},
		{
			confused: ConfusedNumber{
				powersOfTen: []PowerOfTen{powersOfTenByRune['十'], powersOfTenByRune['万']},
				digits: []int{1, 2},
				allowFirstDigit: true,
			},
			digits: []int{1, 2},
			expected: []int{20011},

		},
		{
			confused: ConfusedNumber{
				powersOfTen: []PowerOfTen{powersOfTenByRune['十'], powersOfTenByRune['万']},
				digits: []int{2, 1},
				allowFirstDigit: true,
			},
			digits: []int{2, 1},
			expected: []int{10012, 10020},

		},
	}
	for _, tc := range(testCases) {
		expected := tc.confused.findAllPossibleNumbersForPermutation(tc.digits)
		assert.Equal(t, tc.expected, expected)
	}
}
func TestPermutations (t *testing.T) {
	testCases := []struct{
		original []int
		expected [][]int
	}{
		{
			original: []int{2},
			expected: [][]int{{2}},

		},
		{
			original: []int{5, 2},
			expected: [][]int{{5, 2}, {2, 5}},

		},
		{
			original: []int{2, 2},
			expected: [][]int{{2, 2}, {2, 2}},

		},
		{
			original: []int{5, 2, 1},
			expected: [][]int{{1, 2, 5}, {1, 5, 2}, {5, 2, 1}, {5, 1, 2}, {2, 1, 5}, {2, 5, 1}},

		},
	}
	for _, tc := range(testCases) {
		expected := permutations(tc.original)
		assert.Subset(t, tc.expected, expected)
		assert.Subset(t, expected, tc.expected)
	}
}
