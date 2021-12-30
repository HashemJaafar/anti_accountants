package anti_accountants

import "math"

func LOG(a, b float64) float64  { return math.Log(a) / math.Log(b) }
func ROOT(a, b float64) float64 { return math.Pow(a, 1/b) }
