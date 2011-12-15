package util

func Gcd(n, m int) int {
	for m != 0 {
		n, m = m, n%m
	}
	return n
}

func Lcm(n, m int) int {
	return m / Gcd(n, m) * n
}
