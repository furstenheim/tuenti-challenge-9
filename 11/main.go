package main

import (
	"github.com/lukpank/go-glpk/glpk"
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

type Case struct {
	moons []Moon
	capacity float64
	shipRange float64
}

type Solution struct {
	found bool
	loads []int
}

func (c *Case) solve () Solution {
	lp := glpk.New()
	lp.SetProbName("sample")
	lp.SetObjName("Z")
	lp.SetObjDir(glpk.MAX)
	lp.AddCols(c.getNCols())
	lp.AddRows(c.getNRows())
	tripCheckIndexes, tripCheckValues := c.addTimeMoonConstraints(lp)
	c.setUpQuadraticTerms(lp, tripCheckIndexes, tripCheckValues)
	c.setStartAt0Condition(lp)

	iocp := glpk.NewIocp()
	iocp.SetPresolve(true)

	found := true

	if err := lp.Intopt(iocp); err != nil {
		log.Fatalf("Mip error: %v", err)
	}

	fmt.Printf("%s = %g", lp.ObjName(), lp.MipObjVal())
	loads := []int{}
	for time := 0; time < len(c.moons); time++ {
		for moon, _ := range(c.moons) {
			value := lp.MipColVal(c.getTimeMoonIndex(time, moon))
			if value == 1 {
				loads = append(loads, c.moons[moon].loadInt)
			}
		}
	}
	for i := 0; i < c.getNCols(); i++ {
		fmt.Printf("; %s = %g for variable %d", lp.ColName(i+1), lp.MipColVal(i+1), i + 1)
		fmt.Println()
	}
	for j := 0; j < c.getNRows(); j++ {
		a, b := lp.MatRow(j + 1)
		log.Println(lp.RowName(j +1 ), lp.RowLB(j + 1), a, b, lp.RowUB(j + 1), lp.RowType(j + 1), glpk.LO)
	}
	fmt.Println()
	lp.Delete()
	sort.Slice(loads, func (i, j int) bool {
		return loads[i] < loads[j]
	})
	return Solution{
		loads: loads,
		found: found,
	}
}

func (c * Case) addTimeMoonConstraints (lp * glpk.Prob) (tripCheckIndexes []int32, tripCheckValues[]float64){
	capacityIndexes := []int32{-1}
	capacityValues := []float64{-1}
	tripCheckIndexes = []int32{-1}
	tripCheckValues = []float64{-1}
	lp.SetObjName("Load")
	lp.SetObjDir(glpk.MAX)
	for time := 0; time < len(c.moons); time++ {
		for moon, moonObj := range(c.moons) {
			lp.SetColName(c.getTimeMoonIndex(time, moon), fmt.Sprintf("Spaceshift at time %d at moon %d", time, moon))
			lp.SetColKind(c.getTimeMoonIndex(time, moon), glpk.IV)
			lp.SetColBnds(c.getTimeMoonIndex(time, moon), glpk.DB, 0, 1)
			lp.SetObjCoef(c.getTimeMoonIndex(time, moon), moonObj.load)
			capacityIndexes = append(capacityIndexes, int32(c.getTimeMoonIndex(time, moon)))
			capacityValues = append(capacityValues, moonObj.load)
			tripCheckIndexes = append(tripCheckIndexes, int32(c.getTimeMoonIndex(time, moon)))
			tripCheckValues = append(tripCheckValues, 1)
		}
	}
	lp.SetRowName(c.getCapacityIndex(), "Capacity constraint")
	lp.SetRowBnds(c.getCapacityIndex(), glpk.UP, 0, c.capacity)
	lp.SetMatRow(c.getCapacityIndex(), capacityIndexes, capacityValues)
	for time := 0; time < len(c.moons); time++ {
		indexes := []int32{-1} // 0 index is ignored
		matrixValues := []float64{-1}
		for moon := 0; moon < len(c.moons); moon++ {
			indexes = append(indexes, int32(c.getTimeMoonIndex(time, moon)))
			matrixValues = append(matrixValues, 1)
		}
		lp.SetRowName(c.getSameTimeConditionIndex(time), fmt.Sprintf("Spaceshift at time %d in only one moon", time))
		lp.SetRowBnds(c.getSameTimeConditionIndex(time), glpk.UP, 0, 1)
		lp.SetMatRow(c.getSameTimeConditionIndex(time), indexes, matrixValues)
	}
	for moon := 0; moon < len(c.moons); moon++ {
		indexes := []int32{-1} // 0 index is ignored
		matrixValues := []float64{-1}
		for time := 0; time < len(c.moons); time++ {
			indexes = append(indexes, int32(c.getTimeMoonIndex(time, moon)))
			matrixValues = append(matrixValues, 1)
		}
		lp.SetRowName(c.getSameMoonConditionIndex(moon), fmt.Sprintf("Spaceshift at moon %d in only one time", moon))
		lp.SetRowBnds(c.getSameMoonConditionIndex(moon), glpk.UP, 0, 1)
		lp.SetMatRow(c.getSameMoonConditionIndex(moon), indexes, matrixValues)
	}
	return
}

func (c * Case) setUpQuadraticTerms (lp *glpk.Prob, tripCheckIndexes []int32, tripCheckValues[]float64) {
	rangeIndexes := []int32{-1}
	rangeValues := []float64{-1}
	for time := 0; time < len(c.moons) -1; time++ {
		for m1, moon1 := range(c.moons) {
			for m2, moon2 := range(c.moons) {
				if m1 == m2 {
					continue
				}
				quadraticIndex := c.getQuadraticTermIndex(time, m1, m2)
				m1Index := c.getTimeMoonIndex(time, m1)
				m2Index := c.getTimeMoonIndex(time + 1, m2)
				lp.SetColName(quadraticIndex, fmt.Sprintf("Term z%d%d%d", time, m1, m2))
				lp.SetColBnds(quadraticIndex, glpk.LO, 0, 0)
				baseIndex := c.getBaseQuadraticConditionIndex(time, m1, m2)

				lp.SetRowName(baseIndex, fmt.Sprintf("constraint 0 on z%d%d%d", time, m1, m2))
				lp.SetRowBnds(baseIndex, glpk.LO, 0, 0)
				lp.SetMatRow(baseIndex, []int32{-1, int32(m1Index), int32(quadraticIndex)}, []float64{0, 1, -1})

				lp.SetRowName(baseIndex + 1, fmt.Sprintf("constraint 1 on z%d%d%d", time, m1, m2))
				lp.SetRowBnds(baseIndex + 1, glpk.LO, 0, 0)
				lp.SetMatRow(baseIndex + 1, []int32{-1, int32(m2Index), int32(quadraticIndex)}, []float64{0, 1, -1})

				lp.SetRowName(baseIndex + 2, fmt.Sprintf("constraint 2 on z%d%d%d", time, m1, m2))
				lp.SetRowBnds(baseIndex + 2, glpk.LO, -1, -1)
				lp.SetMatRow(baseIndex + 2, []int32{-1, int32(m1Index), int32(m2Index), int32(quadraticIndex)}, []float64{0, -1, -1, 1})
				rangeIndexes = append(rangeIndexes, int32(quadraticIndex))
				rangeValues = append(rangeValues, moon1.distanceTo(moon2, float64(time + 1)) - moon1.radius)

				tripCheckIndexes = append(tripCheckIndexes, int32(quadraticIndex))
				tripCheckValues = append(tripCheckValues, -1)
			}
		}
	}
	for time := 0; time < len(c.moons); time++ {
		for moon, moonObj := range(c.moons) {
			tIndex := c.getTimeMoonIndex(time, moon)
			rangeIndexes = append(rangeIndexes, int32(tIndex))
			if time == 0 {
				rangeValues = append(rangeValues, 2 * moonObj.radius)
			} else {
				rangeValues = append(rangeValues, moonObj.radius)
			}
		}
	}

	lp.SetRowName(c.getRangeIndex(), "Range is not overdone")
	lp.SetRowBnds(c.getRangeIndex(), glpk.UP, 0, c.shipRange)
	lp.SetMatRow(c.getRangeIndex(), rangeIndexes, rangeValues)

	lp.SetRowName(c.getCheckIndex(), "Check that quadratic terms add up to the planets")
	lp.SetRowBnds(c.getCheckIndex(), glpk.FX, 1, 1)
	lp.SetMatRow(c.getCheckIndex(), tripCheckIndexes, tripCheckValues)
}


func (c *Case) setStartAt0Condition (lp * glpk.Prob) {
	startIndexes := []int32{-1}
	startValues := []float64{-1}
	for moon, _ := range(c.moons) {
		startIndexes = append(startIndexes, int32(c.getTimeMoonIndex(0, moon)))
		startValues = append(startValues, 1)
	}
	lp.SetRowName(c.getStartCondition(), "Check that we start when the weather conditions are correct")
	lp.SetRowBnds(c.getStartCondition(), glpk.FX, 1, 1)
	lp.SetMatRow(c.getStartCondition(), startIndexes, startValues)

}
func (c * Case) getTimeMoonIndex(time, moon int) int {
	return 1 + len(c.moons) * time + moon
}

func (c * Case) getQuadraticTermIndex (time, m1, m2 int) int {
	return 1 + len(c.moons) * len(c.moons) + time * len(c.moons) * len(c.moons) + m1 * len(c.moons) + m2
}

func (c * Case) getNCols () int {
	return 1 + len(c.moons) * len(c.moons) + len(c.moons) * len(c.moons) * len(c.moons)
}

func (c *Case) getSameMoonConditionIndex (moon int) int {
	return 4 + 1 + moon
}
func (c * Case) getSameTimeConditionIndex (time int) int {
	return 4+ 1 + len(c.moons) + time
}
func (c * Case) getBaseQuadraticConditionIndex (time, m1, m2 int) int {
	return 4 + 1 + len(c.moons) + len(c.moons) +  3 * time * len(c.moons) * len(c.moons) + 3 * m1 * len(c.moons) + 3 * m2
}

func (c * Case) getRangeIndex () int {
	return 1
}
func (c * Case) getCapacityIndex () int {
	return 2
}

func (c * Case) getCheckIndex () int {
	return 3
}

func (c * Case) getStartCondition () int {
	return 4
}

func (c * Case) getNRows () int {
	return 3 + 1 + len(c.moons) + len(c.moons) + 3 * len(c.moons) * len(c.moons) * len(c.moons)
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
