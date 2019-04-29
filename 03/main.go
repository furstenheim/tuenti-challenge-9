package main

import (
	"bufio"
	"os"
	"strings"
	"strconv"
	"log"
)

const (
	Top Fold = iota + 1
	Left
	Right
	Bottom
)

var letterToFold = map[byte]Fold{'T': Top, 'L': Left, 'R': Right, 'B': Bottom}

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
		log.Println(c)
		punches := c.Solve()
		printResult(i, punches)
	}
}

type Fold uint
type Punch [2]int

type Case struct {
	x, y int
	nFolds int
	folds []Fold // Folds saved in reversed order (so we pop instead of unshifting)
	punches []Punch
}

func parseCase (reader * bufio.Reader) Case {
	line, e := reader.ReadString('\n')
	handleError(e)
	fields := strings.Fields(line)
	if (len(fields) != 4) {
		log.Fatal("Expected size 4 at header", fields)
	}
	x, e := strconv.Atoi(fields[0])
	handleError(e)
	y, e := strconv.Atoi(fields[1])
	handleError(e)
	nFolds, e := strconv.Atoi(fields[2])
	handleError(e)
	nPunches, e := strconv.Atoi(fields[3])
	folds := parseFolds(nFolds, reader)
	punches := parsePunches(nPunches, reader)
	return Case{
		nFolds: nFolds,
		punches: punches,
		folds: folds,
		x: x,
		y: y,
	}
}

func parseFolds (nFolds int, reader *bufio.Reader) []Fold {
	folds := make([]Fold, nFolds)
	for i := 0; i < nFolds; i++ {
		foldLine, e := reader.ReadString('\n')
		handleError(e)
		f, ok := letterToFold[foldLine[0]]
		if !ok {
			log.Fatal("Unknown fold", foldLine[0])
		}
		folds[len(folds) - 1 - i] = f
	}
	return folds
}

func parsePunches (nPunch int, reader *bufio.Reader) []Punch {
	punches := make([]Punch, nPunch)
	for i := 0; i < nPunch; i++ {
		punchLine, e := reader.ReadString('\n')
		handleError(e)
		fields := strings.Fields(punchLine)
		x, e := strconv.Atoi(fields[0])
		handleError(e)
		y, e := strconv.Atoi(fields[1])
		punches[i] = Punch{x, y}
	}
	return punches
}

func (c * Case) Solve () []Punch {
	return nil
}

func printResult (i int, punches []Punch) {

}


func handleError (e error) {
	if e != nil {
		log.Fatal(e)
	}
}

