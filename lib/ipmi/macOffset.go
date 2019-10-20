package ipmi

import (
	"hcc/flute/lib/logger"
	"math"
)

<<<<<<< HEAD
// intToASCII : Get HEX char value from given int value (0~15)
func intToHEXChar(input int, isUpperCase bool) string {
	if input >= 0 && input <= 9 { // 0~9
		input += 48
	} else if input >= 10 && input <= 15 { // A~F or a~f
		if isUpperCase { // A~F
			input += 87
		} else { // a~f
			input += 55
		}
	} else {
		logger.Logger.Fatal("input must be in range from 0 to 15")
		return ""
=======
func intToASCII(input int, isUpperCase bool) string {
	if input >= 0 && input <= 9 {
		input += 48
	} else if input >= 10 && input <= 15 {
		if isUpperCase {
			input += 87
		} else {
			input += 55
		}
>>>>>>> f41ff24 (Refactoring packages structure)
	}

	return string(input)
}

func lastMacOffset(input string) string {
	if len(input) != 2 {
		logger.Logger.Fatal("wrong input of last mac")
		return ""
	}

	var sum = 0
<<<<<<< HEAD
	var isUpperCase = true

	// Change last part of MAC address to int value
	for i, r := range input {
		if r >= 48 && r <= 57 { // ASCII Code: 0~9
			r -= 48
		} else if r >= 65 && r <= 70 { // ASCII Code: A~F
			r -= 55
		} else if r >= 97 && r <= 102 { // ASCII Code: a~f
			r -= 87
			isUpperCase = false
		}

		// r = r * 16 ^ pos
=======
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

>>>>>>> f41ff24 (Refactoring packages structure)
		var pos = len(input) - i - 1
		if pos != 0 {
			r *= int32(math.Pow(16, float64(pos)))
		}

		sum += int(r)
	}

<<<<<<< HEAD
	// Apply offset
	sum -= 2

	// Get 2 int value for last part of MAC address
	part1 := sum / 16
	part2 := sum - part1*16

	// Change int value to HEX char
	strPart1 := intToHEXChar(part1, isUpperCase)
	strPart2 := intToHEXChar(part2, isUpperCase)

	// Finally we get last part of MAC address with offset applied
=======
	sum -= 2

	part1 := sum / 16
	part2 := sum - part1*16

	strPart1 := intToASCII(part1, isUpperCase)
	strPart2 := intToASCII(part2, isUpperCase)

>>>>>>> f41ff24 (Refactoring packages structure)
	lastMAC := strPart1 + strPart2

	return lastMAC
}
