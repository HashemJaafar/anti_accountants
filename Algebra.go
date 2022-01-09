package anti_accountants

import (
	"fmt"
	"log"
	"math"
	"strconv"
)

func check_map_keys_for_equations(equations [][]string, m map[string]float64) {
	var elements []string
	for _, equation := range equations {
		elements = append(elements, equation[0], equation[1], equation[3])
	}
	elements, _ = RETURN_SET_AND_DUPLICATES_SLICES(elements)
	for keyb := range m {
		if !IS_IN(keyb, elements) {
			log.Panic(keyb, " is not in ", elements)
		}
	}
}

func EQUATIONS_SOLVER(print, check_if_keys_in_the_equations bool, m map[string]float64, equations [][]string) {
	if check_if_keys_in_the_equations {
		check_map_keys_for_equations(equations, m)
	}
	for a := 0; a <= len(equations); a++ {
		for _, equation := range equations {
			equation_solver(print, m, equation[0], equation[1], equation[2], equation[3])
		}
	}
}

func equation_solver(print bool, m map[string]float64, a, b, sign, c string) {
	la, erra := strconv.ParseFloat(a, 64)
	lb, errb := strconv.ParseFloat(b, 64)
	lc, errc := strconv.ParseFloat(c, 64)
	if erra == nil {
		m[a] = la
	}
	if errb == nil {
		m[b] = lb
	}
	if errc == nil {
		m[c] = lc
	}
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
		equations_generator(print, m, a, b, sign, c, POW(m[b], m[c]), ROOT(m[a], m[c]), LOG(m[a], m[b]))
	case "root":
		equations_generator(print, m, a, b, sign, c, ROOT(m[b], m[c]), POW(m[a], m[c]), LOG(m[b], m[a]))
	case "log":
		equations_generator(print, m, a, b, sign, c, LOG(m[b], m[c]), POW(m[c], m[a]), ROOT(m[b], m[a]))
	default:
		log.Panic(sign, " is not in [+-*/**root log]")
	}
}

func equations_generator(print bool, m map[string]float64, a, b, sign, c string, a_value, b_value, c_value float64) {
	la, oka := m[a]
	lb, okb := m[b]
	lc, okc := m[c]
	var inf bool
	for _, a := range []float64{m[a], m[b], m[c]} {
		if math.IsInf(a, 0) {
			inf = true
		}
	}
	if math.IsNaN(la) {
		m[a] = 0
	}
	if math.IsNaN(lb) {
		m[b] = 0
	}
	if math.IsNaN(lc) {
		m[c] = 0
	}
	if math.IsNaN(a_value) {
		a_value = 0
	}
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
	case oka && okb && okc && math.Round(la*1000)/1000 != math.Round(a_value*1000)/1000 && !inf:
		log.Fatal(a, " ", m[a], " != ", b, " ", m[b], " ", sign, " ", c, " ", m[c])
	}
}

func print_equation(print bool, m map[string]float64, a, b, sign, c string) {
	if print {
		fmt.Println(a, m[a], " = ", b, m[b], " ", sign, " ", c, m[c])
	}
}
