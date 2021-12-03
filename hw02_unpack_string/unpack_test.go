package hw02unpackstring

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnpack(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "symbolsRepeat", input: "a4bc2d5e", expected: "aaaabccddddde"},
		{name: "noSymbolRepeat", input: "abccd", expected: "abccd"},
		{name: "emptyString", input: "", expected: ""},
		{name: "zeroDigit", input: "aaa0b", expected: "aab"},
		{name: "zeroDigits", input: "a0b0c0", expected: ""},
		{name: "negativeDigit", input: "aaa-1b", expected: "aaa-b"},
		{name: "negativeDigitRepeat", input: "aaa-3b", expected: "aaa---b"},
		{name: "spaceRepeat", input: "   3   ", expected: "        "},
		{name: "specSymbolsRepeat", input: "!1@2#3$4%5^6&7*8(9)0", expected: "!@@###$$$$%%%%%^^^^^^&&&&&&&********((((((((("},
		{name: "emojiRepeat", input: "🤯3", expected: "🤯🤯🤯"},
		{name: "unicodeRepeat", input: "\u00083", expected: "\u0008\u0008\u0008"},
		{name: "nullRepeat", input: "\u00003", expected: ""},
		// uncomment if task with asterisk completed
		{name: "escapeDigits", input: `qwe\4\5`, expected: `qwe45`},
		{name: "escapeDigitRepeat", input: `qwe\45`, expected: `qwe44444`},
		{name: "escapeSlashRepeat", input: `qwe\\5`, expected: `qwe\\\\\`},
		{name: "escapeSlashAndDigit", input: `qwe\\\3`, expected: `qwe\3`},
		{name: "escapeBegin", input: `\4\5abc`, expected: `45abc`},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name+"("+tc.input+")", func(t *testing.T) {
			result, err := Unpack(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestUnpackInvalidString(t *testing.T) {
	invalidStrings := []string{
		"3abc", "45", "aaa10b", `\abc`, `abc\`, `\-1abc`,
		`\🤯abc`, `\🤯abc`, "\\\u0008abc", "\\\u0000abc",
	}
	for _, tc := range invalidStrings {
		tc := tc
		t.Run(tc, func(t *testing.T) {
			_, err := Unpack(tc)
			require.Truef(t, errors.Is(err, ErrInvalidString), "actual error %q", err)
		})
	}
}
