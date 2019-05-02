package main

import (
	"log"
	"sort"
	"bufio"
	"os"
	"strings"
	"strconv"
	"fmt"
)

func main () {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	firstLineFields := strings.Fields(line)
	numberOfCases, err := strconv.Atoi(firstLineFields[0])
	if (err != nil) {
		log.Fatal(err)
	}
	for i := 0; i < numberOfCases; i ++ {
		log.Println(i)
		c := parseCase(reader)
		s := c.solve()
		s.printResult(i)
	}
}

var delimiterRunes = map[rune]bool{
	'O': true,
	'P': true,
	'E': true,
	'R': true,
	'A': true,
	'T': true,
	' ': true,
	'=': true,
	'\n': true,
}
func parseCase (reader *bufio.Reader) Case {
	line, e := reader.ReadString('\n')
	handleError(e)
	fields := strings.FieldsFunc(line, func (r rune) bool {
		return delimiterRunes[r]
	})
	return Case{
		op1: KanjiNumber(fields[0]),
		op2: KanjiNumber(fields[1]),
		result: KanjiNumber(fields[2]),
	}
}

func (s Solution) printResult (i int) {
	var text1 string
	text1 = fmt.Sprintf("Case #%d: %d %s %d = %d", i + 1, s.n1, s.operator, s.n2, s.result)
	fmt.Println(text1)
}
type Case struct {
	op1, op2, result KanjiNumber
}

type Solution struct {
	operator string
	n1, n2, result int
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

type Operation struct {
	symbol string
	operation func (i, j int) int
}

var operations = []Operation{
	{
		symbol: "+",
		operation: func (i, j int) int {
			return i + j
		},
	},{
		symbol: "-",
		operation: func (i, j int) int {
			return i - j
		},
	},{
		symbol: "*",
		operation: func (i, j int) int {
			return i * j
		},
	},
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



func (c Case) solve () Solution {
	c1 := c.op1.toConfusedNumber()
	c2 := c.op2.toConfusedNumber()
	fm := c.result.toFrequencymap()
	possibleNumbers1 := c1.findAllPossibleNumbers()
	possibleNumbers2 := c2.findAllPossibleNumbers()
	for _, p1 := range(possibleNumbers1) {
		for _, p2 := range(possibleNumbers2) {
			for _, op := range(operations) {
				possibleResult := op.operation(p1, p2)
				if possibleResult <= MAX_NUMBER && possibleResult >= MIN_NUMBER {
					possibleKanji := intToKanji(possibleResult)
					isPossible := fm.isKanjiPossible(possibleKanji)
					if isPossible {
						return Solution{
							operator: op.symbol,
							n1: p1,
							n2: p2,
							result: possibleResult,
						}
					}
				}
			}

		}
	}
	log.Fatal("Solution not possible", string(c.op1), ' ', string(c.op2), ' ', string(c.result))
	return Solution{}
}

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

func (n KanjiNumber) toFrequencymap () FrequencyMap {
	fm := map[rune]int{}
	for _, k := range(n) {
		fm[k]++
	}
	return fm
}


func (c ConfusedNumber) findAllPossibleNumbers () []int {
	digitsPermutations := permutations(c.digits)
	result := []int{}
	for _, digits := range(digitsPermutations) {
		result = append(result, c.findAllPossibleNumbersForPermutation(digits)...)
	}
	return result
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
func handleError (e error) {
	if e != nil {
		log.Fatal(e)
	}
}
