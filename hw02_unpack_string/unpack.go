package hw02unpackstring

import (
	"errors"
	"strconv"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	if len(str) == 0 {
		return "", nil
	}

	runeStr := []rune(str)
	if unicode.IsDigit(runeStr[0]) {
		return "", ErrInvalidString
	}

	sb := newUnpackStringBuilder()
	var prev rune
	var err error
	const slash rune = '\\'
	for _, char := range runeStr {
		if (char == slash || unicode.IsDigit(char)) && prev == slash {
			err = sb.setSymbol(char)
			if err != nil {
				return "", err
			}

			prev = 0
			continue
		}

		switch {
		case unicode.IsDigit(char):
			if unicode.IsDigit(prev) {
				return "", ErrInvalidString
			}

			sb.repeat, err = strconv.Atoi(string(char))
			fallthrough
		case char == slash:
			prev = char
		default:
			if prev == slash {
				return "", ErrInvalidString
			}

			err = sb.setSymbol(char)
			prev = 0
		}

		if err != nil {
			return "", err
		}
	}

	if prev == slash {
		return "", ErrInvalidString
	}

	return sb.writeSymbol()
}
