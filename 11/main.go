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

	lp.AddRows(10)
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


	fmt.Printf("col1: %v\n", lp.ColKind(1) == glpk.CV)

	// index 0 and value 0 are ignored
	ind := []int32{-1, 1, 2, 3, 4}
	mat := [][]float64{
		{0, 1, 0, 1.0, 0},
		{0, 0, 1, 0, 1},
		{0, 1, 1.0, 0.0, 0},
		{0, 0, 0, 1, 1},
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

func (m1 Moon) distanceTo (m2 Moon, t float64) float64 {
	angleDiff := m1.currentAngle(t) - m2.currentAngle(t)
	distance := m1.radius * m1.radius + m2.radius * m2.radius - 2 * m1.radius * m2.radius * math.Cos(angleDiff)
	return math.Sqrt(distance)
}

func (m1 Moon) currentAngle (t float64) float64 {
	return m1.angle0 + 2 * math.Pi * t *  6 / m1.periodHours
}

