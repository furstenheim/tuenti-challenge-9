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
	digits 	    []rune
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
	digits := []rune{}
	for _, k := range(n) {
		if p, ok := powersOfTenByRune[k]; ok {
			confusedPowers = append(confusedPowers, p)
		} else {
			allowFirstDigit = true
			digits = append(digits, k)
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