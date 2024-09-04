package utils_test

import (
	"testing"

	"github.com/nanoteck137/dwebble/tools/utils"
)

func TestExtract(t *testing.T) {
	type test struct {
		s        string
		expected int
	}

	tests := []test{
		{
			s:        "01 - Testing.flac",
			expected: 1,
		},
		{
			s:        "10 - Testing.opus",
			expected: 10,
		},
		{
			s:        "123 - Testing.mp3",
			expected: 123,
		},
		{
			s:        "03 01 02 - Testing.mp4",
			expected: 3,
		},
		{
			s:        "100.wav",
			expected: 100,
		},
		{
			s:        "Hello World.wav",
			expected: 0,
		},
		{
			s:        "Hello 04 World 02 .wav",
			expected: 0,
		},
	}

	for i, test := range tests {
		num := utils.ExtractNumber(test.s)
		if num != test.expected {
			t.Errorf("Test %d Failed: (\"%s\") Expected %d got %d", i, test.s, test.expected, num)
		}
	}
}
