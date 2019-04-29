package main

import (
	"os"
	"strings"
	"strconv"
	"log"
	"bufio"
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
		c := parseCase(reader, i)
		nRoutes := c.Solve()
		printResult(i, nRoutes)
		log.Println(c)
	}
}

type Case struct {
	iCase int
	nPlanets int
	matrix map[Planet][]Planet
}


type Planet string
func parseCase (reader * bufio.Reader, i int) Case {
	planetsLine, e := reader.ReadString('\n')
	handleError(e)
	planetsNumberField := strings.Fields(planetsLine)
	nPlanets, e := strconv.Atoi(planetsNumberField[0])
	handleError(e)
	c := Case{iCase: i + 1, nPlanets: nPlanets, matrix: make(map[Planet][]Planet, nPlanets)}
	for i := 0; i < nPlanets; i++ {
		planetLine, e := reader.ReadString('\n')
		handleError(e)
		planetLineFields := strings.FieldsFunc(planetLine, func (c rune) bool {
			return c == ':' || c == '\n'
		})
		currentPlanet := planetLineFields[0]
		destinationPlanets := strings.FieldsFunc(planetLineFields[1], func (c rune) bool {
			return c == ','
		})
		destinations := make([]Planet, len(destinationPlanets))
		for i, s := range(destinationPlanets) {
			destinations[i] = Planet(s)
		}
		c.matrix[Planet(currentPlanet)] = destinations
	}
	return c
}

func (c * Case) Solve() int {
	routes := make([]Planet, 0)
	nRoutes := 0
	for _, p := range(c.matrix["Galactica"]) {
		routes = append(routes, p)
	}

	for len(routes) > 0 {
		var lastPlanet Planet
		routes, lastPlanet = routes[:len(routes) - 1], routes[len(routes) - 1]
		if lastPlanet == "New Earth" {
			nRoutes++
			continue
		}
		if (len(c.matrix[lastPlanet]) == 0) {
			log.Fatal("Missing length")
		}
		for _, p2 := range(c.matrix[lastPlanet]) {
			routes = append(routes, p2)
		}
	}
	return nRoutes
}

func printResult (i, nRoutes int) {
	text := "Case #" + strconv.Itoa(i + 1) + ": " + strconv.Itoa(nRoutes)
	fmt.Println(text)
}

func handleError (e error) {
	if e != nil {
		log.Fatal(e)
	}
}

