package main

import (
	"log"
	"sort"
)

type Case struct {
	op1, op2, result KanjiNumber
}

type Solution struct {

}

type KanjiNumber []rune
type FrequencyMap map[rune]int
type ProcessedNumber struct {
	raw            KanjiNumber
	frequencyCount FrequencyMap
}
type PowerOfTen struct {
	hasRune, listsUnit bool
	rune               rune
	value int
}

type ConfusedNumber struct {
	allowFirstDigit bool
	powersOfTen []PowerOfTen
	digits 	    []int
}
const MAX_NUMBER = 99999
const MIN_NUMBER = 1
var powersOfTen = []PowerOfTen{
	{
		hasRune: false,
		listsUnit: true,
		rune: 'x',
		value: 1,
	},
	{
		hasRune: true,
		listsUnit: false,
		rune: '十',
		value: 10,
	},
	{
		hasRune: true,
		listsUnit: false,
		rune: '百',
		value: 100,
	},
	{
		hasRune: true,
		listsUnit: false,
		rune: '千',
		value: 1000,
	},
	{
		hasRune: true,
		listsUnit: true,
		rune: '万',
		value: 10000,
	},
}

var digitToKanji = map[int]rune {
	1: '一',
	2: '二',
	3: '三',
	4: '四',
	5: '五',
	6: '六',
	7: '七',
	8: '八',
	9: '九',
}

var runeToDigit = map[rune]int {
	'一': 1,
	'二': 2,
	'三': 3,
	'四': 4,
	'五': 5,
	'六': 6,
	'七': 7,
	'八': 8,
	'九': 9,
}
var powersOfTenByRune = map[rune]PowerOfTen {

}
func init () {
	for _, v := range(powersOfTen) {
		if v.hasRune {
			powersOfTenByRune[v.rune] = v
		}
	}
}

/*
func toRawNumber (n int) KanjiNumber {

}
*/


func (n KanjiNumber) toConfusedNumber () ConfusedNumber {
	allowFirstDigit := false
	confusedPowers := []PowerOfTen{}
	digits := []int{}
	for _, k := range(n) {
		if p, ok := powersOfTenByRune[k]; ok {
			confusedPowers = append(confusedPowers, p)
		} else {
			allowFirstDigit = true
			digits = append(digits, runeToDigit[k])
		}
	}
	sort.Slice(confusedPowers, func (i, j int) bool {
		return confusedPowers[i].value < confusedPowers[j].value
	})
	return ConfusedNumber{
		allowFirstDigit: allowFirstDigit,
		powersOfTen: confusedPowers,
		digits: digits,
	}
}

func (c ConfusedNumber) findAllPossibleNumbersForPermutation (digits []int) []int {
	var helper func (current []int, remainingPowers []PowerOfTen, remainingDigits []int, allowUnits bool, firstGo bool) []int
	helper = func (current []int, remainingPowers []PowerOfTen, remainingDigits []int, allowUnits bool, firstGo bool) []int {
		if len(remainingPowers) == 0 && len(remainingDigits) == 0 {
			return current
		}

		if len(remainingDigits) > 0 && !allowUnits && len(remainingPowers) == 0 {
			return []int{}
		}
		if len(remainingDigits) == 0 && remainingPowers[0].listsUnit {
			return []int{}
		}
		nextCurrent := []int{}
		if allowUnits {
			current2 := append(current, remainingDigits[0])
			nextCurrent = append(nextCurrent, helper(current2, remainingPowers, remainingDigits[1:], false, false)...)
		}
		if len(remainingPowers) > 0 {
			nextPower := remainingPowers[0]
			if !nextPower.listsUnit {
				current2 := make([]int, len(current))
				if firstGo {
					current2 = append(current2, nextPower.value)
				}
				for i, v := range(current) {
					current2[i] = v + nextPower.value
				}
				nextCurrent = append(nextCurrent,
					helper(current2, remainingPowers[1:], remainingDigits, false, false)...
				)
			}
			if len(remainingDigits) > 0 {
				if nextDigit := remainingDigits[0]; nextDigit != 1 || nextPower.listsUnit {
					current2 := make([]int, len(current))
					if firstGo {
						current2 = append(current2, nextPower.value * nextDigit)
					}

					for i, v := range(current) {
						current2[i] = v + nextPower.value * nextDigit
					}

					nextCurrent = append(nextCurrent,
						helper(current2, remainingPowers[1:], remainingDigits[1:], false, false)...
					)
				}
			}

		}
		return nextCurrent
	}

	return helper([]int{}, c.powersOfTen, digits, c.allowFirstDigit, true)
}


func isUnitRune (r rune) bool {
	return r == '一'
}

func intToKanji (n int) KanjiNumber {
	if n > MAX_NUMBER {
		log.Fatal("Number too big", n)
	}
	if n < MIN_NUMBER {
		log.Fatal("Negative number", n)
	}
	kanjiNumber := make(KanjiNumber, 0)
	for _, p := range(powersOfTen) {
		v := n % 10
		n = n / 10
		if v == 0 {
			continue
		}
		k, ok := digitToKanji[v]
		if !ok {
			log.Fatal("Unknown digit", k)
		}
		if !p.hasRune {
			kanjiNumber = append(kanjiNumber, k)
		} else if (!p.listsUnit && v == 1) {
			kanjiNumber = append(kanjiNumber, p.rune)
		} else {
			kanjiNumber = append(kanjiNumber, p.rune, k)
		}

	}
	kanjiNumber.reverse()
	return kanjiNumber

}


func (m FrequencyMap) isKanjiPossible (n KanjiNumber) bool {
	mapCopy := m.clone()
	for _, v := range (n) {
		freq := mapCopy[v]
		if freq == 0 {
			return false
		}
		mapCopy[v]--
	}
	mapIsEmpty := true
	for _, v := range (mapCopy) {
		if v > 0 {
			mapIsEmpty = false
			break
		}
	}
	return mapIsEmpty
}

func (m FrequencyMap) clone () FrequencyMap {
	cp := make(FrequencyMap, len(m))
	for k, v := range (m) {
		cp[k] = v
	}
	return cp
}




func (k KanjiNumber) reverse () {
	for i := len(k)/2-1; i >= 0; i-- {
		opp := len(k)-1-i
		k[i], k[opp] = k[opp], k[i]
	}
}

func permutations(arr []int)[][]int{
	var helper func([]int, int)
	res := [][]int{}

	helper = func(arr []int, n int){
		if n == 1{
			tmp := make([]int, len(arr))
			copy(tmp, arr)
			res = append(res, tmp)
		} else {
			for i := 0; i < n; i++{
				helper(arr, n - 1)
				if n % 2 == 1{
					tmp := arr[i]
					arr[i] = arr[n - 1]
					arr[n - 1] = tmp
				} else {
					tmp := arr[0]
					arr[0] = arr[n - 1]
					arr[n - 1] = tmp
				}
			}
		}
	}
	helper(arr, len(arr))
	return res
}
