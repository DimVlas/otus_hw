package hw02unpackstring

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

var (
	ErrInvalidString = errors.New("invalid string")

	nums = []rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}

	numZero = '0'

	bslash = '\\'
)

// True - если руна является цифрой.
func IsDigit(r rune) bool {
	return slices.Contains(nums, r)
}

func Unpack(text string) (string, error) {
	if text == "" { // с пустой строкой ничего не делаем
		return "", nil
	}

	runes := []rune(text)

	if IsDigit(runes[0]) { // если первая руна цифра
		return "", ErrInvalidString
	}

	var res strings.Builder

	lenRunes := len(runes)
	for i := 0; i < lenRunes; i++ {
		if IsDigit(runes[i]) {
			return "", ErrInvalidString
		}

		if runes[i] == bslash { // текущий символ слэш
			if i == lenRunes-1 { // это последний символ
				return "", ErrInvalidString
			}

			if !IsDigit(runes[i+1]) && runes[i+1] != bslash { // следующий символ не цифра ине слэш
				return "", ErrInvalidString
			}

			i++ // нужно обработать следующий символ как обычный
		}

		if i == lenRunes-1 { // это последний символ
			res.WriteRune(runes[i])
			break
		}
		if !IsDigit(runes[i+1]) { // следующий символ не цифра
			res.WriteRune(runes[i])
			continue
		}
		if IsDigit(runes[i+1]) { // следующий символ цифра
			if runes[i+1] == numZero { // следующий символ '0'
				i++
				continue
			}
			n, err := strconv.Atoi(string(runes[i+1]))
			if err != nil {
				return "", fmt.Errorf("error converting rune '%q' to number: %w", runes[i+1], err)
			}

			res.WriteString(strings.Repeat(string(runes[i]), n))
			i++
			continue
		}
	}

	return res.String(), nil
}
