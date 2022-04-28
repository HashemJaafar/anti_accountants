package main

import (
	"fmt"
	"math"
)

func CHECK_MAP_KEYS_FOR_EQUATIONS(equations [][]string, m map[string]float64) error {
	var elements []string
	for _, equation := range equations {
		elements = append(elements, equation[0], equation[1], equation[3])
	}
	elements, _ = RETURN_SET_AND_DUPLICATES_SLICES(elements)
	for keyb := range m {
		if !IS_IN(keyb, elements) {
			return ERROR_NOT_LISTED
		}
	}
	return nil
}

func EQUATIONS_GENERATOR(print bool, m map[string]float64, a, b, sign, c string, a_value, b_value, c_value float64) {
	la, oka := m[a]
	lb, okb := m[b]
	lc, okc := m[c]
	m[a] = CONVERT_NAN_TO_ZERO(la)
	m[b] = CONVERT_NAN_TO_ZERO(lb)
	m[c] = CONVERT_NAN_TO_ZERO(lc)
	a_value = CONVERT_NAN_TO_ZERO(a_value)
	b_value = CONVERT_NAN_TO_ZERO(b_value)
	c_value = CONVERT_NAN_TO_ZERO(c_value)
	switch {
	case !oka && okb && okc:
		m[a] = a_value
		PRINT_EQUATION(print, m, a, b, sign, c)
	case oka && !okb && okc:
		m[b] = b_value
		PRINT_EQUATION(print, m, a, b, sign, c)
	case oka && okb && !okc:
		m[c] = c_value
		PRINT_EQUATION(print, m, a, b, sign, c)
	case oka && okb && okc && math.Round(la*1000)/1000 != math.Round(a_value*1000)/1000 && !IS_INF_IN(la, lb, lc):
		// fmt.Errorf(m, a, b, sign, c)
	}
}

func EQUATIONS_SOLVER(print, check_if_keys_in_the_equations bool, m map[string]float64, equations [][]string) error {
	if check_if_keys_in_the_equations {
		err := CHECK_MAP_KEYS_FOR_EQUATIONS(equations, m)
		if err != nil {
			return err
		}
	}
	for a := 0; a <= len(equations)+1; a++ {
		for _, equation := range equations {
			EQUATION_SOLVER(print, m, equation[0], equation[1], equation[2], equation[3])
		}
	}
	PRINT_TABLE.Flush()
	return nil
}

func EQUATION_SOLVER(print bool, m map[string]float64, a, b, sign, c string) {
	ASSIGN_NUMBER_IF_NUMBER(m, a)
	ASSIGN_NUMBER_IF_NUMBER(m, b)
	ASSIGN_NUMBER_IF_NUMBER(m, c)
	switch sign {
	case "+":
		EQUATIONS_GENERATOR(print, m, a, b, sign, c, m[b]+m[c], m[a]-m[c], m[a]-m[b])
	case "-":
		EQUATIONS_GENERATOR(print, m, a, b, sign, c, m[b]-m[c], m[a]+m[c], m[b]-m[a])
	case "*":
		EQUATIONS_GENERATOR(print, m, a, b, sign, c, m[b]*m[c], m[a]/m[c], m[a]/m[b])
	case "/":
		EQUATIONS_GENERATOR(print, m, a, b, sign, c, m[b]/m[c], m[a]*m[c], m[b]/m[a])
	case "**":
		EQUATIONS_GENERATOR(print, m, a, b, sign, c, math.Pow(m[b], m[c]), ROOT(m[a], m[c]), LOGARITHM(m[a], m[b]))
	case "ROOT":
		EQUATIONS_GENERATOR(print, m, a, b, sign, c, ROOT(m[b], m[c]), math.Pow(m[a], m[c]), LOGARITHM(m[b], m[a]))
	default:
		EQUATIONS_GENERATOR(print, m, a, b, sign, c, LOGARITHM(m[b], m[c]), math.Pow(m[c], m[a]), ROOT(m[b], m[a]))
	}
}

func FIRST_POINT(points [][2]float64) (float64, float64) {
	return points[0][0], points[0][1]
}

func HIGH_LOW(points [][2]float64) float64 {
	x_max, y_max := MAX(points)
	x_low, y_low := MIN(points)
	return (y_max - y_low) / (x_max - x_low)
}

func LEAST_SQUARES_REGRESSION(points [][2]float64) (float64, float64) {
	var sum_x, sum_y, sum_x_quadratic, sum_xy float64
	for _, i := range points {
		sum_x += i[0]
		sum_y += i[1]
		sum_x_quadratic += math.Pow(i[0], 2)
		sum_xy += i[0] * i[1]
	}
	n := float64(len(points))
	m := (n*sum_xy - sum_x*sum_y) / ((n * sum_x_quadratic) - math.Pow(sum_x, 2))
	b := (sum_y - (m * sum_x)) / n
	return m, b
}

func LOGARITHM(a, b float64) float64 { return math.Log(a) / math.Log(b) }

func MAX(points [][2]float64) (float64, float64) {
	x, y := FIRST_POINT(points)
	for _, i := range points {
		if i[0] > x {
			x = i[0]
			y = i[1]
		}
	}
	return x, y
}

func MIN(points [][2]float64) (float64, float64) {
	x, y := FIRST_POINT(points)
	for _, i := range points {
		if i[0] < x {
			x = i[0]
			y = i[1]
		}
	}
	return x, y
}

func PRINT_EQUATION(print bool, m map[string]float64, a, b, sign, c string) {
	if print {
		fmt.Fprintln(PRINT_TABLE, a, "\t", m[a], "\t", " = ", "\t", b, "\t", m[b], "\t", sign, "\t", c, "\t", m[c])
	}
}

func ROOT(a, b float64) float64 { return math.Pow(a, 1/b) }

func X_UNDER_X(points [][2]float64, x_max float64) (float64, float64) {
	x, y := 0.0, 0.0
	for _, i := range points {
		if i[0] > x && i[0] <= x_max {
			x = i[0]
			y = i[1]
		}
	}
	return x, y
}
