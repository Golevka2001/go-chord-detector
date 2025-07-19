// Reference: https://github.com/tonaljs/tonal/tree/main/packages/pitch-interval/test.ts
package pitchinterval

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenizeInterval(t *testing.T) {
	t.Run("tokenize intervals", func(t *testing.T) {
		result := TokenizeInterval("-2M")
		assert.Equal(t, [2]string{"-2", "M"}, result)

		result = TokenizeInterval("M-3")
		assert.Equal(t, [2]string{"-3", "M"}, result)

		result = TokenizeInterval("4d")
		assert.Equal(t, [2]string{"4", "d"}, result)

		result = TokenizeInterval("P5")
		assert.Equal(t, [2]string{"5", "P"}, result)

		result = TokenizeInterval("1P")
		assert.Equal(t, [2]string{"1", "P"}, result)

		result = TokenizeInterval("invalid")
		assert.Equal(t, [2]string{"", ""}, result)
	})
}

func TestIntervalFromString(t *testing.T) {
	t.Run("has all properties", func(t *testing.T) {
		interval := Parse("4d")
		expected := Interval{
			Empty:     false,
			Name:      "4d",
			Num:       4,
			Q:         d,
			T:         perfectable,
			Alt:       -1,
			Chroma:    4,
			Simple:    4,
			Step:      3,
			Semitones: 4,
			Oct:       0,
		}
		assert.Equal(t, expected, interval)
	})

	t.Run("accepts interval as parameter", func(t *testing.T) {
		original := Parse("5P")
		duplicate := Parse("5P")
		assert.Equal(t, original, duplicate)
	})

	t.Run("interval names", func(t *testing.T) {
		testCases := []struct {
			input    string
			expected string
		}{
			{"1P", "1P"},
			{"2M", "2M"},
			{"3M", "3M"},
			{"4P", "4P"},
			{"5P", "5P"},
			{"6M", "6M"},
			{"7M", "7M"},
			{"P1", "1P"},
			{"M2", "2M"},
			{"M3", "3M"},
			{"P4", "4P"},
			{"P5", "5P"},
			{"M6", "6M"},
			{"M7", "7M"},
			{"-1P", "-1P"},
			{"-2M", "-2M"},
			{"-3M", "-3M"},
			{"-4P", "-4P"},
			{"-5P", "-5P"},
			{"-6M", "-6M"},
			{"-7M", "-7M"},
			{"P-1", "-1P"},
			{"M-2", "-2M"},
			{"M-3", "-3M"},
			{"P-4", "-4P"},
			{"P-5", "-5P"},
			{"M-6", "-6M"},
			{"M-7", "-7M"},
		}

		for _, tc := range testCases {
			interval := Parse(tc.input)
			assert.Equal(t, tc.expected, interval.Name, "Input: %s", tc.input)
		}

		assert.True(t, Parse("not-an-interval").Empty)
		assert.True(t, Parse("2P").Empty)
	})

	t.Run("quality", func(t *testing.T) {
		testCases := []struct {
			input    string
			expected Quaility
		}{
			{"1dd", dd},
			{"1d", d},
			{"1P", P},
			{"1A", A},
			{"1AA", AA},
			{"2dd", dd},
			{"2d", d},
			{"2m", m},
			{"2M", M},
			{"2A", A},
			{"2AA", AA},
		}

		for _, tc := range testCases {
			interval := Parse(tc.input)
			assert.Equal(t, tc.expected, interval.Q, "Input: %s", tc.input)
		}
	})

	t.Run("alteration", func(t *testing.T) {
		testCases := []struct {
			input    string
			expected int
		}{
			{"1dd", -2},
			{"2dd", -3},
			{"3dd", -3},
			{"4dd", -2},
		}

		for _, tc := range testCases {
			interval := Parse(tc.input)
			assert.Equal(t, tc.expected, interval.Alt, "Input: %s", tc.input)
		}
	})

	t.Run("simple interval", func(t *testing.T) {
		testCases := []struct {
			input    string
			expected int
		}{
			{"1P", 1},
			{"2M", 2},
			{"3M", 3},
			{"4P", 4},
			{"8P", 8},
			{"9M", 2},
			{"10M", 3},
			{"11P", 4},
			{"-8P", -8},
			{"-9M", -2},
			{"-10M", -3},
			{"-11P", -4},
		}

		for _, tc := range testCases {
			interval := Parse(tc.input)
			assert.Equal(t, tc.expected, interval.Simple, "Input: %s", tc.input)
		}
	})
}

func TestQToAlt(t *testing.T) {
	t.Run("perfect intervals", func(t *testing.T) {
		assert.Equal(t, 0, qToAlt(perfectable, P))
		assert.Equal(t, -1, qToAlt(perfectable, d))
		assert.Equal(t, -2, qToAlt(perfectable, dd))
		assert.Equal(t, 1, qToAlt(perfectable, A))
		assert.Equal(t, 2, qToAlt(perfectable, AA))
	})

	t.Run("major intervals", func(t *testing.T) {
		assert.Equal(t, 0, qToAlt(majorable, M))
		assert.Equal(t, -1, qToAlt(majorable, m))
		assert.Equal(t, -2, qToAlt(majorable, d))
		assert.Equal(t, -3, qToAlt(majorable, dd))
		assert.Equal(t, 1, qToAlt(majorable, A))
		assert.Equal(t, 2, qToAlt(majorable, AA))
	})
}

func TestNoInterval(t *testing.T) {
	t.Run("empty interval properties", func(t *testing.T) {
		assert.True(t, Nointerval.Empty)
		assert.Equal(t, "", Nointerval.Name)
		assert.Equal(t, 0, Nointerval.Num)
		assert.Equal(t, Quaility(""), Nointerval.Q)
		assert.Equal(t, Type(""), Nointerval.T)
		assert.Equal(t, 0, Nointerval.Step)
		assert.Equal(t, 0, Nointerval.Alt)
		assert.Equal(t, 0, Nointerval.Simple)
		assert.Equal(t, 0, Nointerval.Semitones)
		assert.Equal(t, 0, Nointerval.Chroma)
		assert.Equal(t, 0, Nointerval.Oct)
	})
}
