package anti_accountants

import "math"

func pow(a ...float64) float64 {
	result := a[0]
	for i := 1; i <= len(a)-1; i++ {
		result = math.Pow(result, a[i])
	}
	return result
}

func logarithm(a ...float64) float64 {
	result := a[0]
	for i := 1; i <= len(a)-1; i++ {
		result = math.Log(result) / math.Log(a[i])
	}
	return result
}

func root(a ...float64) float64 {
	result := a[0]
	for i := 1; i <= len(a)-1; i++ {
		result = math.Pow(result, 1/a[i])
	}
	return result
}

func max(points [][2]float64) (float64, float64) {
	x, y := first_point(points)
	for _, i := range points {
		if i[0] > x {
			x = i[0]
			y = i[1]
		}
	}
	return x, y
}

func min(points [][2]float64) (float64, float64) {
	x, y := first_point(points)
	for _, i := range points {
		if i[0] < x {
			x = i[0]
			y = i[1]
		}
	}
	return x, y
}

func x_under_x(points [][2]float64, x_max float64) (float64, float64) {
	x, y := 0.0, 0.0
	for _, i := range points {
		if i[0] > x && i[0] <= x_max {
			x = i[0]
			y = i[1]
		}
	}
	return x, y
}

func first_point(points [][2]float64) (float64, float64) { return points[0][0], points[0][1] }

func HIGH_LOW(points [][2]float64) float64 {
	x_max, y_max := max(points)
	x_low, y_low := min(points)
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
