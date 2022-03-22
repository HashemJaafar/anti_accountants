package anti_accountants

import "math"

func LOGARITHM(a, b float64) float64 {
	return math.Log(a) / math.Log(b)
}

func ROOT(a, b float64) float64 {
	return math.Pow(a, 1/b)
}

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
