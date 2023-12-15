package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

// var ErrInvalidString = errors.New("invalid string")
var ErrInvalidString_firstnumber = errors.New("The first letter should not be digit")
var ErrInvalidString_multiplire = errors.New("The multiplier should be one digit")

func Unpack(input string) (string, error) {
	var result strings.Builder
	runeSlice := []rune(input)
	length := len(runeSlice)

	if length == 0 {
		return "", nil
	} else if unicode.IsDigit(runeSlice[0]) {
		return "", ErrInvalidString_firstnumber
	}

	for i := 0; i < length; i++ {
		currentChar := runeSlice[i]

		if i+1 < length && unicode.IsDigit(runeSlice[i+1]) {
			if i+2 < length && unicode.IsDigit(runeSlice[i+2]) {
				return "", ErrInvalidString_multiplire
			}
			count, _ := strconv.Atoi(string(runeSlice[i+1]))
			result.WriteString(strings.Repeat(string(currentChar), count))
			i++
		} else {
			result.WriteRune(currentChar)
		}
	}

	return result.String(), nil
}
