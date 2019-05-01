package main

import (
	"bufio"
	"os"
	"strings"
	"strconv"
	"log"
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
		s := c.Solve()
		s.printResult(i)
	}
}

type Letter string
type Hint []Letter
type Case struct {
	hints            []Hint
	dictionary       map[Letter]bool
	indegree         map[int]map[Letter]bool
	indegreeByLetter map[Letter]int
	outmatrix        map[Letter]map[Letter]bool
}
type Solution struct {
	correct bool
	order []Letter
}

func parseCase (reader * bufio.Reader) Case {
	line, e := reader.ReadString('\n')
	handleError(e)
	nHints, e := strconv.Atoi(strings.Fields(line)[0])
	handleError(e)
	hints := make([]Hint, nHints)
	for i, _ := range(hints) {
		hintLine, e := reader.ReadString('\n')
		handleError(e)
		hint := strings.Fields(hintLine)[0]
		if len(hint) + 1 != len(hintLine) {
			log.Fatal("Not matching length", hint, hintLine)
		}
		if len(hint) == 0 {
			log.Fatal("Empty clue")
		}
		hArray := make([]Letter, len(hint))
		for i, c := range(hint) {
			hArray[i] = Letter(c)
		}
		hints[i] = Hint(hArray)
	}
	return Case{hints: hints}
}

func (c * Case) Solve () Solution {
	c.computeDictionary()
	c.computeMatrices()
	return c.computeSolution()
}

func (c *Case) computeSolution () Solution {
	order := []Letter{}
	fail := false
	for len(c.dictionary) > 0 {
		if len(c.indegree[0]) != 1 {
			fail = true
			log.Println("indegree for 0 at fail", order, c.indegree[0], c.dictionary, c.indegree)
			break
		}
		var nextLetter Letter
		for l, _ := range(c.indegree[0]) {
			nextLetter = l // only one interation here
		}
		order = append(order, nextLetter)
		delete(c.indegree[0], nextLetter)
		delete(c.dictionary, nextLetter)
		for l, _ := range(c.outmatrix[nextLetter]) {
			currentIndegree := c.indegreeByLetter[l]
			delete(c.indegree[currentIndegree], l)
			indegree := currentIndegree - 1
			_, ok := c.indegree[indegree]
			if ok {
				c.indegree[indegree][l] = true
			} else {
				c.indegree[indegree] = map[Letter]bool{l: true}
			}
			c.indegreeByLetter[l] = indegree
		}
	}
	if fail {
		return Solution{correct: false}
	}
	return Solution{correct: true, order: order}
}

func (c * Case) computeMatrices () {
	for i, h1 := range (c.hints) {
		if i == len(c.hints) - 1 {
			break
		}
		h2 := c.hints[i + 1]
		l1, l2, found := c.compareClues(h1, h2)
		if found {
			alreadyFound := c.outmatrix[l1][l2]
			if alreadyFound {
				continue
			}
			c.outmatrix[l1][l2] = true
			previousIndegree := c.indegreeByLetter[l2]
			delete(c.indegree[previousIndegree], l2)
			indegree := previousIndegree + 1
			_, ok := c.indegree[indegree]
			if ok {
				c.indegree[indegree][l2] = true
			} else {
				c.indegree[indegree] = map[Letter]bool{l2: true}
			}
			c.indegreeByLetter[l2] = indegree
		}
	}
}

func (c *Case) computeDictionary () {
	dictionary := map[Letter]bool{}
	indegree := make(map[int]map[Letter]bool, 1)
	indegree[0] = map[Letter]bool{}
	indegreeByLetter := make(map[Letter]int, 0)
	outmatrix := make(map[Letter]map[Letter]bool, 0)
	for _, h := range(c.hints) {
		for _, l := range(h) {
			dictionary[l] = true
			indegree[0][l] = true
			indegreeByLetter[l] = 0
			outmatrix[l] = make(map[Letter]bool, 0)
		}
	}
	c.dictionary = dictionary
	c.indegree = indegree
	c.indegreeByLetter = indegreeByLetter
	c.outmatrix = outmatrix
}

func (c * Case) compareClues (h1, h2 Hint) (out, in Letter, found bool) {
	for i, l1 := range(h1) {
		if i >= len(h2) {
			// Not really a clue. This is for example gg -> ggg
			log.Println("Found useless clues", h1, h2)
			return "0", "0", false
		}
		l2 := h2[i]
		if l2 != l1 {
			return l1, l2, true
		}
	}
	log.Println("Found useless clue, 2", h1, h2)
	return "0", "0", false
}

func (s * Solution) printResult (i int) {
	var text1 string
	if !s.correct {
		text1 = fmt.Sprintf("Case #%d: AMBIGUOUS", i + 1)
	} else {
		orderText := []string{}
		for _, l := range(s.order) {
			orderText = append(orderText, string(l))
		}
		text1 = fmt.Sprintf("Case #%d: %s", i + 1, strings.Join(orderText, " "))
	}
	fmt.Println(text1)
}

func handleError (e error) {
	if e != nil {
		log.Fatal(e)
	}
}
