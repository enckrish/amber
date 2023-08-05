package main

// Max returns the larger of x or y.
func Max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

// Min returns the smaller of x or y.
func Min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

func ClIncr(v int, s int) int {
	return (v + 1) % s
}

func ClDecr(v int, s int) int {
	r := (v - 1) % s
	if r < 0 {
		r += s
	}
	return r
}
