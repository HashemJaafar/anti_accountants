package main

import (
	"fmt"
	"math"
)

func FRoot(a, b float64) float64      { return math.Pow(a, 1/b) }
func FLogarithm(a, b float64) float64 { return math.Log(a) / math.Log(b) }

func FCheckMapKeysForEquations(equations [][]string, m map[string]float64) error {
	var elements []string
	for _, v1 := range equations {
		elements = append(elements, v1[0], v1[1], v1[3])
	}
	elements, _ = FReturnSetAndDuplicatesSlices(elements)
	for k1 := range m {
		_, isIn := FFind(k1, elements)
		if !isIn {
			return VErrorNotListed
		}
	}
	return nil
}

func FEquationsGenerator(print bool, m map[string]float64, a, b, sign, c string, aValue, bValue, cValue float64) {
	la, oka := m[a]
	lb, okb := m[b]
	lc, okc := m[c]
	m[a] = FConvertNanToZero(la)
	m[b] = FConvertNanToZero(lb)
	m[c] = FConvertNanToZero(lc)
	aValue = FConvertNanToZero(aValue)
	bValue = FConvertNanToZero(bValue)
	cValue = FConvertNanToZero(cValue)
	switch {
	case !oka && okb && okc:
		m[a] = aValue
		FPrintEquation(print, m, a, b, sign, c)
	case oka && !okb && okc:
		m[b] = bValue
		FPrintEquation(print, m, a, b, sign, c)
	case oka && okb && !okc:
		m[c] = cValue
		FPrintEquation(print, m, a, b, sign, c)
	case oka && okb && okc && math.Round(la*1000)/1000 != math.Round(aValue*1000)/1000 && !FIsInfIn(la, lb, lc):
	}
}

func FEquationsSolver(print, checkIfKeysInTheEquations bool, m map[string]float64, equations [][]string) error {
	if checkIfKeysInTheEquations {
		err := FCheckMapKeysForEquations(equations, m)
		if err != nil {
			return err
		}
	}
	for k1 := 0; k1 <= len(equations)+1; k1++ {
		for _, v2 := range equations {
			FEquationSolver(print, m, v2[0], v2[1], v2[2], v2[3])
		}
	}
	VPrintTable.Flush()
	return nil
}

func FEquationSolver(print bool, m map[string]float64, a, b, sign, c string) {
	FAssignNumberIfNumber(m, a)
	FAssignNumberIfNumber(m, b)
	FAssignNumberIfNumber(m, c)
	switch sign {
	case "+":
		FEquationsGenerator(print, m, a, b, sign, c, m[b]+m[c], m[a]-m[c], m[a]-m[b])
	case "-":
		FEquationsGenerator(print, m, a, b, sign, c, m[b]-m[c], m[a]+m[c], m[b]-m[a])
	case "*":
		FEquationsGenerator(print, m, a, b, sign, c, m[b]*m[c], m[a]/m[c], m[a]/m[b])
	case "/":
		FEquationsGenerator(print, m, a, b, sign, c, m[b]/m[c], m[a]*m[c], m[b]/m[a])
	case "**":
		FEquationsGenerator(print, m, a, b, sign, c, math.Pow(m[b], m[c]), FRoot(m[a], m[c]), FLogarithm(m[a], m[b]))
	case "Root":
		FEquationsGenerator(print, m, a, b, sign, c, FRoot(m[b], m[c]), math.Pow(m[a], m[c]), FLogarithm(m[b], m[a]))
	default:
		FEquationsGenerator(print, m, a, b, sign, c, FLogarithm(m[b], m[c]), math.Pow(m[c], m[a]), FRoot(m[b], m[a]))
	}
}

func FFirstPoint(points [][2]float64) (float64, float64) {
	return points[0][0], points[0][1]
}

func FHighLow(points [][2]float64) float64 {
	xMax, yMax := FMax(points)
	xLow, yLow := FMin(points)
	return (yMax - yLow) / (xMax - xLow)
}

func FLeastSquaresRegression(points [][2]float64) (float64, float64) {
	var sumX, sumY, sumXQuadratic, sumXY float64
	for _, v1 := range points {
		sumX += v1[0]
		sumY += v1[1]
		sumXQuadratic += math.Pow(v1[0], 2)
		sumXY += v1[0] * v1[1]
	}
	n := float64(len(points))
	m := (n*sumXY - sumX*sumY) / ((n * sumXQuadratic) - math.Pow(sumX, 2))
	b := (sumY - (m * sumX)) / n
	return m, b
}

func FMax(points [][2]float64) (float64, float64) {
	x, y := FFirstPoint(points)
	for _, v1 := range points {
		if v1[0] > x {
			x = v1[0]
			y = v1[1]
		}
	}
	return x, y
}

func FMin(points [][2]float64) (float64, float64) {
	x, y := FFirstPoint(points)
	for _, v1 := range points {
		if v1[0] < x {
			x = v1[0]
			y = v1[1]
		}
	}
	return x, y
}

func FPrintEquation(print bool, m map[string]float64, a, b, sign, c string) {
	if print {
		fmt.Fprintln(VPrintTable, a, "\t", m[a], "\t", " = ", "\t", b, "\t", m[b], "\t", sign, "\t", c, "\t", m[c])
	}
}
