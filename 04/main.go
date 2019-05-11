package main

import (
	"bufio"
	"os"
	"strings"
	"strconv"
	"log"
	"math/big"
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
		n, d := c.Solve()
		printResult(i, n, d)
	}
}

type Case struct {
	samples map[int]int
}

func parseCase (reader * bufio.Reader) Case {
	line, e := reader.ReadString('\n')
	handleError(e)
	nSample, e := strconv.Atoi(strings.Fields(line)[0])
	handleError(e)
	samplesLine, e := reader.ReadString('\n')
	samples := strings.Fields(samplesLine)
	if len(samples) != nSample {
		log.Fatal("Unexpected length for line", samples, nSample)
	}
	sampleMap := make(map[int]int)
	for _, s := range(samples) {
		v, e := strconv.Atoi(s)
		handleError(e)
		sampleMap[v]++
	}
	return Case{samples: sampleMap}
}

func (c * Case) Solve () (n, d *big.Int) {
	lcm := c.getLCMOfSamples()
	totalCount := big.NewInt(0)
	totalSum := big.NewInt(0)
	for i, count := range (c. samples) {
		countBig := big.NewInt(int64(count))
		iBig := big.NewInt(int64(i))
		countBig.Mul(lcm, countBig)
		totalSum.Add(countBig, totalSum)
		countBig.Div(countBig, iBig)
		totalCount.Add(totalCount, countBig)
	}
	gcd := big.NewInt(0)
	gcd.GCD(nil, nil, totalCount, totalSum)
	totalCount.Div(totalCount, gcd)
	totalSum.Div(totalSum, gcd)
	return totalSum, totalCount
}

func (c *Case) getLCMOfSamples () *big.Int {
	lcm := big.NewInt(1)
	for i, count := range(c.samples) {
		g := gcd(i, count)
		factor := big.NewInt(int64(i / g))
		g2 := new(big.Int)
		g2.GCD(nil, nil, factor, lcm)
		factor.Div(factor, g2)
		lcm.Mul(factor, lcm)
	}
	return lcm
}

func printResult (i int, n, d *big.Int) {
	text1 := "Case #" + strconv.Itoa(i + 1) + ": " + n.String() + "/" + d.String()
	fmt.Println(text1)
}


func handleError (e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func gcd(x, y int) int {
	return int(new(big.Int).GCD(nil, nil, big.NewInt(int64(x)), big.NewInt(int64(y))).Int64())
}
