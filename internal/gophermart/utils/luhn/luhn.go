package luhn

import "strconv"

// luhnAlgorithm - function for luhn algorithm
func luhnAlgorithm(n int) int {

	sum := 0
	for i := 0; n > 0; i++ {
		c := n % 10
		if i%2 == 0 {
			c *= 2
			if c > 9 {
				c = c%10 + c/10
			}
		}
		sum += c
		n /= 10
	}
	return sum % 10
}

// Validate - validate order
func Validate(order string) bool {
	check, err := strconv.Atoi(order)
	if err != nil {
		return false
	}
	return (check%10+luhnAlgorithm(check/10))%10 == 0
}
