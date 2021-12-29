package anti_accountants

import (
	"fmt"
	"log"
	"math"
	"strconv"
)

type MATH struct {
	POINTS [][]float64
}

func EQUATIONS_SOLVER(print bool, m map[string]float64, equations [][]string) {
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
		equations_generator(print, m, a, b, sign, c, math.Pow(m[b], m[c]), ROOT(m[a], m[c]), LOG(m[a], m[b]))
	case "root":
		equations_generator(print, m, a, b, sign, c, ROOT(m[b], m[c]), math.Pow(m[a], m[c]), LOG(m[b], m[a]))
	case "log":
		equations_generator(print, m, a, b, sign, c, LOG(m[b], m[c]), math.Pow(m[c], m[a]), ROOT(m[b], m[a]))
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
		log.Fatal(a, m[a], " != ", b, m[b], " ", sign, " ", c, m[c])
	}
}

func print_equation(print bool, m map[string]float64, a, b, sign, c string) {
	if print {
		fmt.Println(a, m[a], " = ", b, m[b], " ", sign, " ", c, m[c])
	}
}

func LOG(a, b float64) float64  { return math.Log(a) / math.Log(b) }
func ROOT(a, b float64) float64 { return math.Pow(a, 1/b) }

func (s MATH) HIGH_LOW() float64 {
	var y2, y1, x2, x1 float64
	for _, i := range s.POINTS {
		if i[0] >= x2 {
			x2 = i[0]
			y2 = i[1]
		} else if i[0] < x1 {
			x1 = i[0]
			y1 = i[1]
		}
	}
	return (y2 - y1) / (x2 - x1)
}

func (s MATH) LEAST_SQUARES_REGRESSION() (float64, float64) {
	var sum_x, sum_y, sum_x_quadratic, sum_xy float64
	for _, i := range s.POINTS {
		sum_x += i[0]
		sum_y += i[1]
		sum_x_quadratic += math.Pow(i[0], 2)
		sum_xy += i[0] * i[1]
	}
	n := float64(len(s.POINTS))
	m := (n*sum_xy - sum_x*sum_y) / ((n * sum_x_quadratic) - math.Pow(sum_x, 2))
	b := (sum_y - (m * sum_x)) / n
	return m, b
}
