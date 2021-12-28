package anti_accountants

import (
	"fmt"
	"log"
	"math"
)

type Math struct {
	Points [][]float64
}

func Equations_solver(m map[string]float64, equations [][]string) {
	for a := 0; a <= len(equations); a++ {
		for _, equation := range equations {
			equation_solver(m, equation[0], equation[1], equation[2], equation[3])
		}
	}
}

func equation_solver(m map[string]float64, a, b, sign, c string) {
	switch sign {
	case "+":
		equations_generator(m, a, b, sign, c, m[b]+m[c], m[a]-m[c], m[a]-m[b])
	case "-":
		equations_generator(m, a, b, sign, c, m[b]-m[c], m[a]+m[c], m[b]-m[a])
	case "*":
		equations_generator(m, a, b, sign, c, m[b]*m[c], m[a]/m[c], m[a]/m[b])
	case "/":
		equations_generator(m, a, b, sign, c, m[b]/m[c], m[a]*m[c], m[b]/m[a])
	case "**":
		equations_generator(m, a, b, sign, c, math.Pow(m[b], m[c]), Root(m[a], m[c]), Log(m[a], m[b]))
	case "root":
		equations_generator(m, a, b, sign, c, Root(m[b], m[c]), math.Pow(m[a], m[c]), Log(m[b], m[a]))
	case "log":
		equations_generator(m, a, b, sign, c, Log(m[b], m[c]), math.Pow(m[c], m[a]), Root(m[b], m[a]))
	default:
		log.Panic(sign, " is not in [+-*/**root log]")
	}
}

func equations_generator(m map[string]float64, a, b, sign, c string, a_value, b_value, c_value float64) {
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
		print_equation(m, a, b, sign, c)
	case oka && !okb && okc:
		m[b] = b_value
		print_equation(m, a, b, sign, c)
	case oka && okb && !okc:
		m[c] = c_value
		print_equation(m, a, b, sign, c)
	case oka && okb && okc && math.Round(la*1000)/1000 != math.Round(a_value*1000)/1000 && !inf:
		log.Fatal(a, m[a], " != ", b, m[b], " ", sign, " ", c, m[c])
	}
}

func print_equation(m map[string]float64, a, b, sign, c string) {
	fmt.Println(a, m[a], " = ", b, m[b], " ", sign, " ", c, m[c])
}

func Log(a, b float64) float64  { return math.Log(a) / math.Log(b) }
func Root(a, b float64) float64 { return math.Pow(a, 1/b) }

func (s Math) High_low() float64 {
	var y2, y1, x2, x1 float64
	for _, i := range s.Points {
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

func (s Math) Least_squares_regression() (float64, float64) {
	var sum_x, sum_y, sum_x_quadratic, sum_xy float64
	for _, i := range s.Points {
		sum_x += i[0]
		sum_y += i[1]
		sum_x_quadratic += math.Pow(i[0], 2)
		sum_xy += i[0] * i[1]
	}
	n := float64(len(s.Points))
	m := (n*sum_xy - sum_x*sum_y) / ((n * sum_x_quadratic) - math.Pow(sum_x, 2))
	b := (sum_y - (m * sum_x)) / n
	return m, b
}
