package main


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

/*
func toRawNumber (n int) KanjiNumber {

}
*/

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





