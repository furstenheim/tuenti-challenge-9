package main

import (
	"bufio"
	"os"
	"strings"
	"strconv"
	"log"
	"sort"
	"fmt"
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
	punches := c.punches
	folds := c.folds
	x := c.x
	y := c.y
	for len(folds) > 0 {
		var nextFold Fold
		folds, nextFold = folds[:len(folds) - 1], folds[len(folds) - 1]
		newPunches := make([]Punch, 0, len(punches) * 2)
		for _, p := range(punches) {
			if nextFold == Top {
				newPunches = append(
					newPunches,
					Punch{p[0], p[1] + y},
					Punch{p[0], y - p[1]- 1},
				)
				y = 2 * y
			} else if nextFold == Bottom {
				newPunches = append(
					newPunches,
					Punch{p[0], p[1]},
					Punch{p[0], 2 * y - p[1]},
				)
				y = 2 * y
			} else if nextFold == Right {
				newPunches = append(
					newPunches,
					Punch{p[0], p[1]},
					Punch{2 * x - p[0], p[1]},
				)
				x = 2 * x
			} else if nextFold == Left {
				newPunches = append(
					newPunches,
					Punch{p[0] + x, p[1]},
					Punch{x - p[0], p[1]},
				)
				x = 2 * x
			} else {
				log.Fatal("Unknown fold", nextFold)
			}
		}
		punches = newPunches
	}
	log.Println(punches)
	sort.Slice(punches, func (i, j int) bool {
		log.Println(i, j)
		return punches[i][0] < punches[j][0] ||
			(punches[i][0] == punches[j][0] && punches[i][1] < punches[j][1])
	})
	return punches
}

func printResult (i int, punches []Punch) {

	text1 := "Case #" + strconv.Itoa(i + 1) + ":"
	fmt.Println(text1)
	for _, p := range(punches) {
		text := strconv.Itoa(p[0]) + " " + strconv.Itoa(p[1])
		fmt.Println(text)
	}
}


func handleError (e error) {
	if e != nil {
		log.Fatal(e)
	}
}

