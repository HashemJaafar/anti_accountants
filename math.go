package anti_accountants

import (
	"fmt"
	"math"
	"strconv"
)

func FRoot(a, b float64) float64      { return math.Pow(a, 1/b) }
func FLogarithm(a, b float64) float64 { return math.Log(a) / math.Log(b) }

func FEquationsSolver(print, checkIfKeysInTheEquations bool, m map[string]float64, equations [][]string) error {
	if checkIfKeysInTheEquations {
	bigLoop:
		for k1 := range m {
			for _, v1 := range equations {
				if k1 == v1[0] || k1 == v1[1] || k1 == v1[3] {
					continue bigLoop
				}
			}
			return fmt.Errorf("key (%s) is not in the equations", k1)
		}
	}
	for range equations {
		for _, v2 := range equations {
			FEquationSolver(print, m, v2[0], v2[1], v2[2], v2[3])
		}
	}
	VPrintTable.Flush()
	return nil
}

func FEquationSolver(print bool, m map[string]float64, a, b, sign, c string) {
	equationsGenerator := func(aValue, bValue, cValue float64) {
		printEquation := func() {
			if print {
				fmt.Fprintln(VPrintTable, a, "\t", m[a], "\t", " = ", "\t", b, "\t", m[b], "\t", sign, "\t", c, "\t", m[c])
			}
		}

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
			printEquation()
		case oka && !okb && okc:
			m[b] = bValue
			printEquation()
		case oka && okb && !okc:
			m[c] = cValue
			printEquation()
		case oka && okb && okc && math.Round(la*1000)/1000 != math.Round(aValue*1000)/1000 && !FIsInfIn(la, lb, lc):
		}
	}

	assignNumberIfNumber := func(str string) {
		number, err := strconv.ParseFloat(str, 64)
		if err == nil {
			m[str] = number
		}
	}

	assignNumberIfNumber(a)
	assignNumberIfNumber(b)
	assignNumberIfNumber(c)

	switch sign {
	case "+":
		equationsGenerator(m[b]+m[c], m[a]-m[c], m[a]-m[b])
	case "-":
		equationsGenerator(m[b]-m[c], m[a]+m[c], m[b]-m[a])
	case "*":
		equationsGenerator(m[b]*m[c], m[a]/m[c], m[a]/m[b])
	case "/":
		equationsGenerator(m[b]/m[c], m[a]*m[c], m[b]/m[a])
	case "**":
		equationsGenerator(math.Pow(m[b], m[c]), FRoot(m[a], m[c]), FLogarithm(m[a], m[b]))
	case "Root":
		equationsGenerator(FRoot(m[b], m[c]), math.Pow(m[a], m[c]), FLogarithm(m[b], m[a]))
	default:
		equationsGenerator(FLogarithm(m[b], m[c]), math.Pow(m[c], m[a]), FRoot(m[b], m[a]))
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

func FMaxMin[t INumber](slice []t) (t, t) {
	var max, min t
	for _, v1 := range slice {
		if v1 > max {
			max = v1
		}
		if v1 < min {
			min = v1
		}
	}
	return max, min
}

func FConvertNanToZero(value float64) float64 {
	if math.IsNaN(value) {
		return 0
	}
	return value
}
