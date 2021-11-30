package hw02unpackstring

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnpack(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "a4bc2d5e", expected: "aaaabccddddde"},
		{input: "abccd", expected: "abccd"},
		{input: "", expected: ""},
		{input: "aaa0b", expected: "aab"},
		{input: "a0", expected: ""},
		{input: "a1", expected: "a"},
		{input: "d\n5abc", expected: "d\n\n\n\n\nabc"},
		{input: "a3a0b0", expected: "aaa"},
		{input: "*&^$@", expected: "*&^$@"},
		{input: "*2&4", expected: "**&&&&"},
		{input: "дваб2айта3", expected: "дваббайтааа"},
		{input: "д0вабайта0", expected: "вабайт"},
		{input: "при\n2вет", expected: "при\n\nвет"},
		{input: "異体字3", expected: "異体字字字"},

		{input: `qwe\4\5`, expected: `qwe45`},
		{input: `qwe\45`, expected: `qwe44444`},
		{input: `qwe\\5`, expected: `qwe\\\\\`},
		{input: `qwe\\\3`, expected: `qwe\3`},
		{input: `check\\\\`, expected: `check\\`},
		{input: `двабайта\\`, expected: `двабайта\`},
		{input: `число3\22`, expected: `числооо22`},
		{input: `число0\02`, expected: `числ00`},
		{input: `один\1абв`, expected: `один1абв`},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			result, err := Unpack(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestUnpackInvalidString(t *testing.T) {
	invalidStrings := []string{"3abc", "45", "aaa10b", "aaa01bb", "a123b", `check\*`}
	for _, tc := range invalidStrings {
		tc := tc
		t.Run(tc, func(t *testing.T) {
			_, err := Unpack(tc)
			require.Truef(t, errors.Is(err, ErrInvalidString), "actual error %q", err)
		})
	}
}
