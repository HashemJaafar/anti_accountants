package anti_accountants

import (
	"math"
)

func POW(a ...float64) float64 {
	result := a[0]
	for i := 1; i <= len(a)-1; i++ {
		result = math.Pow(result, a[i])
	}
	return result
}

func LOG(a ...float64) float64 {
	result := a[0]
	for i := 1; i <= len(a)-1; i++ {
		result = math.Log(result) / math.Log(a[i])
	}
	return result
}

func ROOT(a ...float64) float64 {
	result := a[0]
	for i := 1; i <= len(a)-1; i++ {
		result = math.Pow(result, 1/a[i])
	}
	return result
}

func HIGH_LOW(points [][]float64) float64 {
	var y2, y1, x2, x1 float64
	for _, i := range points {
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

func LEAST_SQUARES_REGRESSION(points [][]float64) (float64, float64) {
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
