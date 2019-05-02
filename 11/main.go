package main

import (
	"fmt"
	"log"
	"math"
	"bufio"
	"os"
	"strings"
	"strconv"
	"sort"
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

func parseCase (reader *bufio.Reader) Case {
	line, e := reader.ReadString('\n')
	handleError(e)
	nMoons, e := strconv.Atoi(strings.Fields(line)[0])
	handleError(e)
	moons := make([]Moon, nMoons)
	radiusLine, e := reader.ReadString('\n')
	handleError(e)
	for i, s := range(strings.Fields(radiusLine)) {
		moons[i].radius, e  = strconv.ParseFloat(s, 64)
		handleError(e)
	}
	angleLine, e := reader.ReadString('\n')
	handleError(e)
	for i, s := range(strings.Fields(angleLine)) {
		moons[i].angle0, e = strconv.ParseFloat(s, 64)
		handleError(e)
	}

	periodLine, e := reader.ReadString('\n')
	handleError(e)
	for i, s := range(strings.Fields(periodLine)) {
		moons[i].periodHours, e = strconv.ParseFloat(s, 64)
		handleError(e)
	}

	loadLine, e := reader.ReadString('\n')
	handleError(e)
	for i, s := range(strings.Fields(loadLine)) {
		moons[i].loadInt, e = strconv.Atoi(s)
		handleError(e)
		moons[i].load = float64(moons[i].loadInt)
	}

	capacityLine, e := reader.ReadString('\n')
	handleError(e)
	capacity, e := strconv.Atoi(strings.Fields(capacityLine)[0])
	handleError(e)

	shipRangeLine, e := reader.ReadString('\n')
	handleError(e)
	shipRange, e := strconv.ParseFloat(strings.Fields(shipRangeLine)[0], 64)
	handleError(e)

	return Case{
		moons: moons,
		capacity: float64(capacity),
		visitedCache: map[Position]CacheItem{},
		shipRange: shipRange,
	}
}

func (s *Solution) printResult (i int) {
	var text string
	if s.found {
		values := []string{}
		for _, v := range(s.loads) {
			values = append(values, fmt.Sprintf("%d", v))
		}
		text1 := strings.Join(values, " ")
		text = fmt.Sprintf("Case #%d: %s", i + 1, text1)
	} else {
		text = fmt.Sprintf("Case #%d: None", i + 1)
	}
	fmt.Println(text)
}
type Moon struct {
	load        float64
	radius      float64
	angle0      float64
	periodHours float64
	loadInt int
}


type Position struct {
	currentMoonIndex int
	visitedMoons     [15]bool
}
type CacheItem struct {
	valid bool
	currentRemainingRange float64
	totalCapacityObtained float64
}
type VisitState struct {
	time float64
	remainingRange, remainingCapacity float64
}
type Case struct {
	moons []Moon
	capacity float64
	shipRange float64
	visitedCache map[Position]CacheItem
}

type Solution struct {
	found bool
	loads []int
}

func (c *Case) solve () Solution {
	found := false
	var bestVisits [15]bool
	capacity := -1.

	for i, m := range (c.moons) {
		moons := [15]bool{}
		moons[i] = true
		branchFound, branchCapacity, branchBestVisits := c.bestBranch(Position{
			currentMoonIndex: i,
			visitedMoons: moons,
		}, VisitState{
			time: 0,
			remainingRange: c.shipRange - m.radius,
			remainingCapacity: c.capacity,
		})
		if branchFound && branchCapacity > capacity {
			capacity = branchCapacity
			bestVisits = branchBestVisits
			found = true
		}
	}
	loads := make([]int, 0)
	for i, b := range(bestVisits) {
		if b {
			loads = append(loads, c.moons[i].loadInt)
		}
	}
	sort.Slice(loads, func (i, j int) bool {
		return loads[i] < loads[j]
	})
	return Solution{
		found: found,
		loads: loads,
	}
}

func (c *Case) bestBranch (position Position, vs VisitState) (found bool, capacity float64,  bestVisits [15]bool) {
	if p, ok := c.visitedCache[position]; ok && p.currentRemainingRange > vs.remainingRange {
		return false, -1, [15]bool{}
	}
	c.visitedCache[position] = CacheItem{
		currentRemainingRange: vs.remainingRange,
	}
	currentMoon := c.moons[position.currentMoonIndex]
	if vs.remainingCapacity < currentMoon.load {
		return false, -1,  [15]bool{}
	}
	if vs.remainingRange < currentMoon.radius {
		return false, -1, [15]bool{}
	}
	currentCapacity := currentMoon.load
	capacity = 0
	bestVisits = position.visitedMoons
	for i, m := range(c.moons) {
		if !position.visitedMoons[i] {
			copyMoons := copyMoonIndexes(position.visitedMoons)
			copyMoons[i] = true
			nextTime := vs.time + 1
			vs2 := VisitState{
				remainingRange: vs.remainingRange - currentMoon.distanceTo(m, nextTime),
				remainingCapacity: vs.remainingCapacity - m.load,
			}
			position2 := Position{
				currentMoonIndex: i,
				visitedMoons: copyMoons,
			}
			branchFound, branchCapacity , branchBestVisits := c.bestBranch(position2, vs2) // TODO define other terms
			if branchFound && branchCapacity > capacity {
				bestVisits = branchBestVisits
				capacity = branchCapacity
				found = branchFound
			}
		}
	}
	return true, capacity + currentCapacity, bestVisits
}

func copyMoonIndexes (indexes [15]bool) [15]bool {
	return indexes
}


func (m1 Moon) distanceTo (m2 Moon, t float64) float64 {
	angleDiff := m1.currentAngle(t) - m2.currentAngle(t)
	distance := m1.radius * m1.radius + m2.radius * m2.radius - 2 * m1.radius * m2.radius * math.Cos(angleDiff)
	return math.Sqrt(distance)
}

func (m1 Moon) currentAngle (t float64) float64 {
	return m1.angle0 + 2 * math.Pi * t *  6 / m1.periodHours
}

func handleError (e error) {
	if e != nil {
		log.Fatal(e)
	}
}
