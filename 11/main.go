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

	lp.AddRows(3)
	lp.SetRowName(1, "c1")
	lp.SetRowBnds(1, glpk.UP, 0.0, 20.0)
	lp.SetRowName(2, "c2")
	lp.SetRowBnds(2, glpk.UP, 0.0, 30.0)
	lp.SetRowName(3, "c3")
	lp.SetRowBnds(3, glpk.FX, 0.0, 0)

	lp.AddCols(4)
	lp.SetColName(1, "x1")
	lp.SetColBnds(1, glpk.DB, 0.0, 40.0)
	lp.SetObjCoef(1, 1.0)
	lp.SetColName(2, "x2")
	lp.SetColBnds(2, glpk.LO, 0.0, 0.0)
	lp.SetObjCoef(2, 2.0)
	lp.SetColName(3, "x3")
	lp.SetColBnds(3, glpk.LO, 0.0, 0.0)
	lp.SetObjCoef(3, 3.0)
	lp.SetColName(4, "x4")
	lp.SetColBnds(4, glpk.DB, 2.0, 3.0)
	lp.SetObjCoef(4, 1.0)
	lp.SetColKind(4, glpk.IV)

	fmt.Printf("col1: %v\n", lp.ColKind(1) == glpk.CV)

	ind := []int32{0, 1, 2, 3, 4}
	mat := [][]float64{
		{0, -1, 1.0, 1.0, 10},
		{0, 1.0, -3.0, 1.0, 0.0},
		{0, 0.0, 1.0, 0.0, -3.5}}
	for i := 0; i < 3; i++ {
		lp.SetMatRow(i+1, ind, mat[i])
	}

	iocp := glpk.NewIocp()
	iocp.SetPresolve(true)

	if err := lp.Intopt(iocp); err != nil {
		log.Fatalf("Mip error: %v", err)
	}

	fmt.Printf("%s = %g", lp.ObjName(), lp.MipObjVal())
	for i := 0; i < 4; i++ {
		fmt.Printf("; %s = %g", lp.ColName(i+1), lp.MipColVal(i+1))
	}
	fmt.Println()

	lp.Delete()
}

type Moon struct {
	load        int
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

