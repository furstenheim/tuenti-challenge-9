package main

import (
	"github.com/lukpank/go-glpk/glpk"
	"fmt"
	"log"
	"math"
)

func main() {
	lp := glpk.New()
	lp.SetProbName("sample")
	lp.SetObjName("Z")
	lp.SetObjDir(glpk.MAX)

	shipRange := 6.
	capacity := 20.
	m0 := Moon{
		load        : 4,
		radius      : 2.0,
		angle0      : 0,
		periodHours : 12,
	}
	m1 := Moon{
		load        : 5,
		radius      : 2.5,
		angle0      : 3.14,
		periodHours : 100.0,
	}
	nCols := 6
	lp.AddCols(nCols)
	lp.SetColName(1, "c00")
	lp.SetColBnds(1, glpk.DB, 0.0, 1.0)
	lp.SetObjCoef(1, m0.load)
	lp.SetColKind(1, glpk.IV)
	lp.SetColName(2, "c10")
	lp.SetColBnds(2, glpk.DB, 0.0, 1.0)
	lp.SetObjCoef(2, m0.load)
	lp.SetColKind(2, glpk.IV)
	lp.SetColName(3, "c01")
	lp.SetColBnds(3, glpk.DB, 0.0, 1.0)
	lp.SetObjCoef(3, m1.load)
	lp.SetColKind(3, glpk.IV)
	lp.SetColName(4, "c11")
	lp.SetColBnds(4, glpk.DB, 0.0, 1.0)
	lp.SetObjCoef(4, m1.load)
	lp.SetColKind(4, glpk.IV)

	lp.SetColName(5, "z001")
	lp.SetColBnds(5, glpk.LO, 0, 0)
	lp.SetColName(6, "z010")
	lp.SetColBnds(6, glpk.LO, 0, 0)

	lp.AddRows(12)
	lp.SetRowName(1, "one planet at t 0")
	// TODO check this is not glpk.DB
	lp.SetRowBnds(1, glpk.UP, 0.0, 1)
	lp.SetRowName(2, "one planet at t 1")
	lp.SetRowBnds(2, glpk.UP, 0.0, 1)
	lp.SetRowName(3, "visited 0 only once")
	lp.SetRowBnds(3, glpk.UP, 0.0, 1)
	lp.SetRowName(4, "visited 1 only once")
	lp.SetRowBnds(4, glpk.UP, 0.0, 1)

	lp.SetRowName(5, "constraint 0 on z001")
	lp.SetRowBnds(5, glpk.LO, 0, 0)
	lp.SetRowName(6, "constraint 1 on z001")
	lp.SetRowBnds(6, glpk.LO, 0, 0)
	lp.SetRowName(7, "constraint 2 on z001")
	lp.SetRowBnds(7, glpk.LO, -1, 0)

	lp.SetRowName(8, "constraint 0 on z010")
	lp.SetRowBnds(8, glpk.LO, 0, 0)
	lp.SetRowName(9, "constraint 1 on z010")
	lp.SetRowBnds(9, glpk.LO, 0, 0)
	lp.SetRowName(10, "constraint 2 on z010")
	lp.SetRowBnds(10, glpk.LO, -1, 0)

	lp.SetRowName(11, "constraint on range")
	lp.SetRowBnds(11, glpk.UP, 0, shipRange)

	lp.SetRowName(12, "Constraints on capacity")
	lp.SetRowBnds(12, glpk.UP, 0, capacity)


	// TODO set capacity cap
	fmt.Printf("col1: %v\n", lp.ColKind(1) == glpk.CV)

	// index 0 and value 0 are ignored
	ind := []int32{-1, 1, 2, 3, 4, 5, 6}
	mat := [][]float64{
		{0, 1, 0, 1.0, 0, 0, 0},
		{0, 0, 1, 0, 1, 0, 0},
		{0, 1, 1.0, 0.0, 0, 0, 0},
		{0, 0, 0, 1, 1, 0, 0},

		{0, 1, 0, 0, 0, -1, 0},
		{0, 0, 0, 0, 1, -1, 0},
		{0, -1, 0, 0, -1, 1, 0},

		{0, 0, 1, 0, 0, 0, -1},
		{0, 0, 0, 1, 0, 0, -1},
		{0, 0, -1, -1, 0, 0, 1},

		{0, m0.radius, m0.radius, m1.radius, m1.radius, m0.distanceTo(m1, 1), m1.distanceTo(m0, 1)},
		{0, m0.load, m0.load, m1.load, m1.load, 0, 0},



	}
	for i := 0; i < len(mat); i++ {
		// main constraints
		lp.SetMatRow(i+1, ind, mat[i])
	}

	iocp := glpk.NewIocp()
	iocp.SetPresolve(true)

	if err := lp.Intopt(iocp); err != nil {
		log.Fatalf("Mip error: %v", err)
	}

	fmt.Printf("%s = %g", lp.ObjName(), lp.MipObjVal())
	for i := 0; i < nCols; i++ {
		fmt.Printf("; %s = %g", lp.ColName(i+1), lp.MipColVal(i+1))
	}
	fmt.Println()

	lp.Delete()
}

type Moon struct {
	load        float64
	radius      float64
	angle0      float64
	periodHours float64
}

type Case struct {
	moons []Moon
	capacity float64
	shipRange float64
}

func (c *Case) solve () {
	lp := glpk.New()
	lp.SetProbName("sample")
	lp.SetObjName("Z")
	lp.SetObjDir(glpk.MAX)
	c.addTimeMoonConstraints(lp)
	c.setUpQuadraticTerms(lp)

}

func (c * Case) addTimeMoonConstraints (lp * glpk.Prob) {
	capacityIndexes := []int32{-1}
	capacityValues := []float64{-1}
	lp.SetObjName("Load")
	lp.SetObjDir(glpk.MAX)
	for time := 0; time < len(c.moons); time++ {
		for moon, moonObj := range(c.moons) {
			lp.SetColName(c.getTimeMoonIndex(time, moon), fmt.Sprintf("Spaceshift at time %d at moon %d", time, moon))
			lp.SetColKind(c.getTimeMoonIndex(time, moon), glpk.IV)
			lp.SetColBnds(c.getTimeMoonIndex(time, moon), glpk.DB, 0, 1)
			lp.SetObjCoef(c.getTimeMoonIndex(time, moon), moonObj.load)
			capacityIndexes = append(capacityIndexes, c.getTimeMoonIndex(time, moon))
			capacityValues = append(capacityValues, moonObj.load)
		}
	}
	lp.SetMatRow(c.getCapacityIndex(), capacityIndexes, capacityValues)
	for time := 0; time < len(c.moons); time++ {
		indexes := []int{-1} // 0 index is ignored
		matrixValues := []int{-1}
		for moon := 0; moon < len(c.moons); moon++ {
			indexes = append(indexes, c.getTimeMoonIndex(time, moon))
			matrixValues = append(matrixValues, 1)
		}
		lp.SetRowName(c.getSameTimeConditionIndex(time), fmt.Sprintf("Spaceshift at time %d in only one moon", time))
		lp.SetRowBnds(c.getSameTimeConditionIndex(time), glpk.UP, 0, 1)
		lp.SetMatRow(c.getSameTimeConditionIndex(time), indexes, matrixValues)
	}
	for moon := 0; moon < len(c.moons); moon++ {
		indexes := []int{-1} // 0 index is ignored
		matrixValues := []int{-1}
		for time := 0; time < len(c.moons); time++ {
			indexes = append(indexes, c.getTimeMoonIndex(time, moon))
			matrixValues = append(matrixValues, 1)
		}
		lp.SetRowName(c.getSameMoonConditionIndex(moon), fmt.Sprintf("Spaceshift at moon %d in only one time", moon))
		lp.SetRowBnds(c.getSameMoonConditionIndex(moon), glpk.UP, 0, 1)
		lp.SetMatRow(c.getSameMoonConditionIndex(moon), indexes, matrixValues)
	}
}

func (c * Case) setUpQuadraticTerms (lp *glpk.Prob) {
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
				m2Index := c.getTimeMoonIndex(time, m2)
				lp.SetColName(quadraticIndex, fmt.Sprintf("Term z%d%d%d", time, m1, m2))
				lp.SetColBnds(quadraticIndex, glpk.LO, 0, 0)
				baseIndex := c.getBaseQuadraticConditionIndex(time, m1, m2)

				lp.SetRowName(baseIndex, fmt.Sprintf("constraint 0 on z%d%d%d", time, m1, m2))
				lp.SetRowBnds(baseIndex, glpk.LO, 0, 0)
				lp.SetMatRow(baseIndex, []int32{-1, m1Index, quadraticIndex}, []float64{0, 1, -1})

				lp.SetRowName(baseIndex, fmt.Sprintf("constraint 1 on z%d%d%d", time, m1, m2))
				lp.SetRowBnds(baseIndex, glpk.LO, 0, 0)
				lp.SetMatRow(baseIndex, []int32{-1, m2Index, quadraticIndex}, []float64{0, 1, -1})

				lp.SetRowName(baseIndex, fmt.Sprintf("constraint 2 on z%d%d%d", time, m1, m2))
				lp.SetRowBnds(baseIndex, glpk.LO, -1, 0)
				lp.SetMatRow(baseIndex, []int32{-1, m1Index, m2Index, quadraticIndex}, []float64{0, -1, -1, 1})
				rangeIndexes = append(rangeIndexes, quadraticIndex)
				rangeValues = append(rangeValues, moon1.distanceTo(moon2, float64(time + 1)))
			}
		}
	}
	lp.SetRowName(c.getRangeIndex(), "Range is not overdone")
	lp.SetRowBnds(c.getRangeIndex(), glpk.UP, 0, c.shipRange)
	lp.SetMatRow(c.getRangeIndex(), rangeIndexes, rangeValues)
}

func (c * Case) getTimeMoonIndex(time, moon int) int {
	return 1 + len(c.moons) * time + moon
}

func (c * Case) getQuadraticTermIndex (time, m1, m2 int) int {
	return 1 + len(c.moons) * len(c.moons) + time * len(c.moons) * len(c.moons) + m1 * len(c.moons) + m2
}

func (c *Case) getSameMoonConditionIndex (moon int) int {
	return 2 + 1 + moon
}
func (c * Case) getSameTimeConditionIndex (time int) int {
	return 2+ 1 + len(c.moons) + time
}
func (c * Case) getBaseQuadraticConditionIndex (time, m1, m2 int) int {
	return 2 + 1 + len(c.moons) + len(c.moons) +  3 * time * len(c.moons) * len(c.moons) + m1 * len(c.moons) + m2
}

func (c * Case) getRangeIndex () int {
	return 1
}
func (c * Case) getCapacityIndex () int {
	return 2
}

func (m1 Moon) distanceTo (m2 Moon, t float64) float64 {
	angleDiff := m1.currentAngle(t) - m2.currentAngle(t)
	distance := m1.radius * m1.radius + m2.radius * m2.radius - 2 * m1.radius * m2.radius * math.Cos(angleDiff)
	return math.Sqrt(distance)
}

func (m1 Moon) currentAngle (t float64) float64 {
	return m1.angle0 + 2 * math.Pi * t *  6 / m1.periodHours
}

