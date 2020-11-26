package ipmi

import (
	"hcc/flute/lib/logger"
	"math"
)

func intToASCII(input int, isUpperCase bool) string {
	if input >= 0 && input <= 9 {
		input += 48
	} else if input >= 10 && input <= 15 {
		if isUpperCase {
			input += 87
		} else {
			input += 55
		}
	}

	return string(input)
}

func lastMacOffset(input string) string {
	if len(input) != 2 {
		logger.Logger.Fatal("wrong input of last mac")
		return ""
	}

	var sum = 0
	var isUpperCase = false
	for i, r := range input {
		if r >= 48 && r <= 57 {
			r -= 48
		} else if r >= 65 && r <= 70 {
			r -= 55
		} else if r >= 97 && r <= 102 {
			r -= 87
			isUpperCase = true
		}

		var pos = len(input) - i - 1
		if pos != 0 {
			r *= int32(math.Pow(16, float64(pos)))
		}

		sum += int(r)
	}

	sum -= 2

	part1 := sum / 16
	part2 := sum - part1*16

	strPart1 := intToASCII(part1, isUpperCase)
	strPart2 := intToASCII(part2, isUpperCase)

	lastMAC := strPart1 + strPart2

	return lastMAC
}
