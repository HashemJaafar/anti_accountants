package main

import (
	"fmt"
	"math"
)

func Root(a, b float64) float64      { return math.Pow(a, 1/b) }
func Logarithm(a, b float64) float64 { return math.Log(a) / math.Log(b) }

func CheckMapKeysForEquations(equations [][]string, m map[string]float64) error {
	var elements []string
	for _, v1 := range equations {
		elements = append(elements, v1[0], v1[1], v1[3])
	}
	elements, _ = ReturnSetAndDuplicatesSlices(elements)
	for k1 := range m {
		if !IsIn(k1, elements) {
			return ErrorNotListed
		}
	}
	return nil
}

func EquationsGenerator(print bool, m map[string]float64, a, b, sign, c string, aValue, bValue, cValue float64) {
	la, oka := m[a]
	lb, okb := m[b]
	lc, okc := m[c]
	m[a] = ConvertNanToZero(la)
	m[b] = ConvertNanToZero(lb)
	m[c] = ConvertNanToZero(lc)
	aValue = ConvertNanToZero(aValue)
	bValue = ConvertNanToZero(bValue)
	cValue = ConvertNanToZero(cValue)
	switch {
	case !oka && okb && okc:
		m[a] = aValue
		PrintEquation(print, m, a, b, sign, c)
	case oka && !okb && okc:
		m[b] = bValue
		PrintEquation(print, m, a, b, sign, c)
	case oka && okb && !okc:
		m[c] = cValue
		PrintEquation(print, m, a, b, sign, c)
	case oka && okb && okc && math.Round(la*1000)/1000 != math.Round(aValue*1000)/1000 && !IsInfIn(la, lb, lc):
		// fmt.Errorf(m, a, b, sign, c)
	}
}

func EquationsSolver(print, checkIfKeysInTheEquations bool, m map[string]float64, equations [][]string) error {
	if checkIfKeysInTheEquations {
		err := CheckMapKeysForEquations(equations, m)
		if err != nil {
			return err
		}
	}
	for k1 := 0; k1 <= len(equations)+1; k1++ {
		for _, v2 := range equations {
			EquationSolver(print, m, v2[0], v2[1], v2[2], v2[3])
		}
	}
	PrintTable.Flush()
	return nil
}

func EquationSolver(print bool, m map[string]float64, a, b, sign, c string) {
	AssignNumberIfNumber(m, a)
	AssignNumberIfNumber(m, b)
	AssignNumberIfNumber(m, c)
	switch sign {
	case "+":
		EquationsGenerator(print, m, a, b, sign, c, m[b]+m[c], m[a]-m[c], m[a]-m[b])
	case "-":
		EquationsGenerator(print, m, a, b, sign, c, m[b]-m[c], m[a]+m[c], m[b]-m[a])
	case "*":
		EquationsGenerator(print, m, a, b, sign, c, m[b]*m[c], m[a]/m[c], m[a]/m[b])
	case "/":
		EquationsGenerator(print, m, a, b, sign, c, m[b]/m[c], m[a]*m[c], m[b]/m[a])
	case "**":
		EquationsGenerator(print, m, a, b, sign, c, math.Pow(m[b], m[c]), Root(m[a], m[c]), Logarithm(m[a], m[b]))
	case "Root":
		EquationsGenerator(print, m, a, b, sign, c, Root(m[b], m[c]), math.Pow(m[a], m[c]), Logarithm(m[b], m[a]))
	default:
		EquationsGenerator(print, m, a, b, sign, c, Logarithm(m[b], m[c]), math.Pow(m[c], m[a]), Root(m[b], m[a]))
	}
}

func FirstPoint(points [][2]float64) (float64, float64) {
	return points[0][0], points[0][1]
}

func HighLow(points [][2]float64) float64 {
	xMax, yMax := Max(points)
	xLow, yLow := Min(points)
	return (yMax - yLow) / (xMax - xLow)
}

func LeastSquaresRegression(points [][2]float64) (float64, float64) {
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

func Max(points [][2]float64) (float64, float64) {
	x, y := FirstPoint(points)
	for _, v1 := range points {
		if v1[0] > x {
			x = v1[0]
			y = v1[1]
		}
	}
	return x, y
}

func Min(points [][2]float64) (float64, float64) {
	x, y := FirstPoint(points)
	for _, v1 := range points {
		if v1[0] < x {
			x = v1[0]
			y = v1[1]
		}
	}
	return x, y
}

func PrintEquation(print bool, m map[string]float64, a, b, sign, c string) {
	if print {
		fmt.Fprintln(PrintTable, a, "\t", m[a], "\t", " = ", "\t", b, "\t", m[b], "\t", sign, "\t", c, "\t", m[c])
	}
}
