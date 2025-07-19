package chordtype

import (
	"sort"
	"strings"
	"testing"

	"github.com/Golevka2001/go-chord-detector/pcset"
	"github.com/Golevka2001/go-chord-detector/pitchinterval"
	"github.com/stretchr/testify/assert"
)

func TestNames(t *testing.T) {
	names := Names()

	expectedFirst5 := []string{"fifth", "suspended fourth", "suspended fourth seventh", "augmented", "major seventh flat sixth"}
	assert.Equal(t, expectedFirst5, names[:5], "First 5 names should match expected order")
}

func TestSymbols(t *testing.T) {
	symbols := Symbols()

	expectedFirst3 := []string{"5", "M7#5sus4", "7#5sus4"}
	assert.Equal(t, expectedFirst3, symbols[:3], "First 3 symbols should match expected order")
}

func TestAll(t *testing.T) {
	allChords := All()
	assert.Len(t, allChords, 106, "Should have exactly 106 chord types")
}

func TestGet(t *testing.T) {
	major := Get("major")
	expected := ChordType{
		Pcset: pcset.Pcset{
			Empty:      false,
			SetNum:     2192,
			Chroma:     "100010010000",
			Normalized: "100001000100",
			Intervals:  []string{"1P", "3M", "5P"},
		},
		Name:      "major",
		Quality:   Major,
		Intervals: []string{"1P", "3M", "5P"},
		Aliases:   []string{"M", "^", "", "maj"},
	}
	assert.Equal(t, expected, major, "Should return correct major chord")
}

func TestAdd(t *testing.T) {
	// Add simple chord
	Add([]string{"1P", "5P"}, []string{"q"}, "")
	quinta := Get("q")
	assert.Equal(t, "100000010000", quinta.Chroma, "Should have correct chroma")

	// Add with name
	Add([]string{"1P", "5P"}, []string{"q"}, "quinta")
	quintaByName := Get("quinta")
	assert.Equal(t, Get("q"), quintaByName, "Should get same chord by name")
}

func TestRemoveAll(t *testing.T) {
	RemoveAll()
	assert.Empty(t, All(), "Should have no chords after RemoveAll")
	assert.Empty(t, Keys(), "Should have no keys after RemoveAll")
}

func TestChordData(t *testing.T) {
	intervals := make([]string, 0, len(chords))
	for _, data := range chords {
		intervals = append(intervals, data[0])
	}
	sort.Strings(intervals)

	t.Run("no repeated intervals", func(t *testing.T) {
		for i := 1; i < len(intervals); i++ {
			assert.NotEqual(t, intervals[i-1], intervals[i], "Intervals should not be repeated")
		}
	})

	t.Run("all chords must have abbreviations", func(t *testing.T) {
		for _, data := range chords {
			if len(data) >= 3 {
				abbreviations := strings.TrimSpace(data[2])
				assert.Greater(t, len(abbreviations), 0, "Chord should have abbreviations: %v", data)
			}
		}
	})

	t.Run("intervals should be in ascending order", func(t *testing.T) {
		for _, data := range chords {
			intervalList := data[0]
			intervals := strings.Split(intervalList, " ")

			var semitones []int
			for _, intervalStr := range intervals {
				interval := pitchinterval.Parse(intervalStr)
				if !interval.Empty {
					semitones = append(semitones, interval.Semitones)
				}
			}

			for i := 1; i < len(semitones); i++ {
				assert.Less(t, semitones[i-1], semitones[i],
					"Intervals should be in ascending order for chord: %v", data)
			}
		}
	})
}
