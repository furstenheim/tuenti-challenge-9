package main

import (
	"bufio"
	"os"
	"strings"
	"strconv"
	"log"
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

type Case struct {
	clues     []Letter
	indegree  map[int]map[Letter]bool
	outdegree map[Letter]int
	outmatrix map[Letter][]Letter
}
type Solution struct {

}

func parseCase (reader * bufio.Reader) Case {
	line, e := reader.ReadString('\n')
	handleError(e)
	nClues, e := strconv.Atoi(strings.Fields(line)[0])
	handleError(e)
	clues := make([]Letter, nClues)
	for i, _ := range(clues) {
		clueLine, e := reader.ReadString('\n')
		handleError(e)
		clue := strings.Fields(clueLine)[0]
		if len(clue) + 1 != len(clueLine) {
			log.Fatal("Not matching length", clue, clueLine)
		}
		clues[i] = clue
	}
	return Case{clues: clues}
}

func (c * Case) Solve () Solution {
	return Solution{}
}

func (s * Solution) printResult (i int) {

}

func handleError (e error) {
	if e != nil {
		log.Fatal(e)
	}
}
