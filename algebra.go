package anti_accountants

import (
	"fmt"
	"math"
	"strconv"
)

func check_map_keys_for_equations(equations [][]string, m map[string]float64) error {
	var elements []string
	for _, equation := range equations {
		elements = append(elements, equation[0], equation[1], equation[3])
	}
	elements, _ = return_set_and_duplicates_slices(elements)
	for keyb := range m {
		if !is_in(keyb, elements) {
			return error_not_listed
		}
	}
	return nil
}

func EQUATIONS_SOLVER(print, check_if_keys_in_the_equations bool, m map[string]float64, equations [][]string) {
	if check_if_keys_in_the_equations {
		check_map_keys_for_equations(equations, m)
	}
	for a := 0; a <= len(equations)+1; a++ {
		for _, equation := range equations {
			equation_solver(print, m, equation[0], equation[1], equation[2], equation[3])
		}
	}
	print_table.Flush()
}

func equation_solver(print bool, m map[string]float64, a, b, sign, c string) {
	assign_number_if_number(m, a)
	assign_number_if_number(m, b)
	assign_number_if_number(m, c)
	switch sign {
	case "+":
		equations_generator(print, m, a, b, sign, c, m[b]+m[c], m[a]-m[c], m[a]-m[b])
	case "-":
		equations_generator(print, m, a, b, sign, c, m[b]-m[c], m[a]+m[c], m[b]-m[a])
	case "*":
		equations_generator(print, m, a, b, sign, c, m[b]*m[c], m[a]/m[c], m[a]/m[b])
	case "/":
		equations_generator(print, m, a, b, sign, c, m[b]/m[c], m[a]*m[c], m[b]/m[a])
	case "**":
		equations_generator(print, m, a, b, sign, c, pow(m[b], m[c]), root(m[a], m[c]), logarithm(m[a], m[b]))
	case "root":
		equations_generator(print, m, a, b, sign, c, root(m[b], m[c]), pow(m[a], m[c]), logarithm(m[b], m[a]))
	default:
		equations_generator(print, m, a, b, sign, c, logarithm(m[b], m[c]), pow(m[c], m[a]), root(m[b], m[a]))
	}
}

func assign_number_if_number(m map[string]float64, str string) {
	number, err := strconv.ParseFloat(str, 64)
	if err == nil {
		m[str] = number
	}
}

func equations_generator(print bool, m map[string]float64, a, b, sign, c string, a_value, b_value, c_value float64) {
	la, oka := m[a]
	lb, okb := m[b]
	lc, okc := m[c]
	convert_nan_to_zero_for_map(m, a, la)
	convert_nan_to_zero_for_map(m, b, lb)
	convert_nan_to_zero_for_map(m, c, lc)
	a_value = convert_nan_to_zero(a_value)
	b_value = convert_nan_to_zero(b_value)
	c_value = convert_nan_to_zero(c_value)
	switch {
	case !oka && okb && okc:
		m[a] = a_value
		print_equation(print, m, a, b, sign, c)
	case oka && !okb && okc:
		m[b] = b_value
		print_equation(print, m, a, b, sign, c)
	case oka && okb && !okc:
		m[c] = c_value
		print_equation(print, m, a, b, sign, c)
	case oka && okb && okc && math.Round(la*1000)/1000 != math.Round(a_value*1000)/1000 && !check_if_inf(m, a, b, c):
		// fmt.Errorf(m, a, b, sign, c)
	}
}

func check_if_inf(m map[string]float64, a string, b string, c string) bool {
	for _, a := range []float64{m[a], m[b], m[c]} {
		if math.IsInf(a, 0) {
			return true
		}
	}
	return false
}

func convert_nan_to_zero(value float64) float64 {
	if math.IsNaN(value) {
		value = 0
	}
	return value
}

func convert_nan_to_zero_for_map(m map[string]float64, str string, value float64) {
	if math.IsNaN(value) {
		m[str] = 0
	}
}

func print_equation(print bool, m map[string]float64, a, b, sign, c string) {
	if print {
		fmt.Fprintln(print_table, a, "\t", m[a], "\t", " = ", "\t", b, "\t", m[b], "\t", sign, "\t", c, "\t", m[c])
	}
}
