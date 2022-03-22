package anti_accountants

import (
	"fmt"
	"math"
	"strconv"
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

func ASSIGN_NUMBER_IF_NUMBER(m map[string]float64, str string) {
	number, err := strconv.ParseFloat(str, 64)
	if err == nil {
		m[str] = number
	}
}

func EQUATIONS_GENERATOR(print bool, m map[string]float64, a, b, sign, c string, a_value, b_value, c_value float64) {
	la, oka := m[a]
	lb, okb := m[b]
	lc, okc := m[c]
	CONVERT_NAN_TO_ZERO_FOR_MAP(m, a, la)
	CONVERT_NAN_TO_ZERO_FOR_MAP(m, b, lb)
	CONVERT_NAN_TO_ZERO_FOR_MAP(m, c, lc)
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
	case oka && okb && okc && math.Round(la*1000)/1000 != math.Round(a_value*1000)/1000 && !CHECK_IF_INF(m, a, b, c):
		// fmt.Errorf(m, a, b, sign, c)
	}
}

func CHECK_IF_INF(m map[string]float64, a string, b string, c string) bool {
	for _, a := range []float64{m[a], m[b], m[c]} {
		if math.IsInf(a, 0) {
			return true
		}
	}
	return false
}

func CONVERT_NAN_TO_ZERO(value float64) float64 {
	if math.IsNaN(value) {
		return 0
	}
	return value
}

func CONVERT_NAN_TO_ZERO_FOR_MAP(m map[string]float64, str string, value float64) {
	if math.IsNaN(value) {
		m[str] = 0
	}
}

func PRINT_EQUATION(print bool, m map[string]float64, a, b, sign, c string) {
	if print {
		fmt.Fprintln(PRINT_TABLE, a, "\t", m[a], "\t", " = ", "\t", b, "\t", m[b], "\t", sign, "\t", c, "\t", m[c])
	}
}
