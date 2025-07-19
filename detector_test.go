// Reference: https://github.com/tonaljs/tonal/tree/main/packages/chord-detect/test.ts
package detector

import (
	"testing"

	"github.com/go-music-theory/music-theory/note"
	"github.com/stretchr/testify/assert"
)

func createNotes(noteNames []string) []*note.Note {
	var notes []*note.Note
	for _, name := range noteNames {
		noteClass := note.ClassNamed(name)
		if noteClass != note.Nil {
			notes = append(notes, &note.Note{Class: noteClass})
		}
	}
	return notes
}

func TestDetect(t *testing.T) {
	notes := createNotes([]string{"D", "F#", "A", "C"})
	result := Detect(notes)
	assert.Contains(t, result, "D7", "Should detect D7 chord")

	notes = createNotes([]string{"F#", "A", "C", "D"})
	result = Detect(notes)
	assert.Contains(t, result, "D7/F#", "Should detect D7/F# inversion")

	notes = createNotes([]string{"A", "C", "D", "F#"})
	result = Detect(notes)
	assert.Contains(t, result, "D7/A", "Should detect D7/A inversion")

	notes = createNotes([]string{"E", "G#", "B", "C#"})
	result = Detect(notes)
	hasE6 := false
	hasCSharpM7 := false
	for _, chord := range result {
		if chord == "E6" {
			hasE6 = true
		}
		if chord == "C#m7/E" {
			hasCSharpM7 = true
		}
	}
	assert.True(t, hasE6 && hasCSharpM7, "Should detect E6 and C#m7/E, got: %v", result)
}

func TestAssumePerfectFifth(t *testing.T) {
	notes := createNotes([]string{"D", "F", "C"})
	result := DetectWithOptions(notes, DetectOptions{AssumePerfectFifth: true})
	assert.Contains(t, result, "Dm7", "Should detect Dm7 with assumePerfectFifth=true")

	result = DetectWithOptions(notes, DetectOptions{AssumePerfectFifth: false})
	assert.Empty(t, result, "Should not detect any chord with assumePerfectFifth=false")

	notes = createNotes([]string{"D", "F", "A", "C"})
	result = DetectWithOptions(notes, DetectOptions{AssumePerfectFifth: true})
	assert.Contains(t, result, "Dm7", "Should detect Dm7 with complete chord and assumePerfectFifth=true")

	result = DetectWithOptions(notes, DetectOptions{AssumePerfectFifth: false})
	assert.Contains(t, result, "Dm7", "Should detect Dm7 with complete chord and assumePerfectFifth=false")

	notes = createNotes([]string{"D", "F", "Ab", "C"})
	result = DetectWithOptions(notes, DetectOptions{AssumePerfectFifth: true})
	hasDm7b5 := false
	hasFm6 := false
	for _, chord := range result {
		if chord == "Dm7b5" {
			hasDm7b5 = true
		}
		if chord == "Fm6/D" {
			hasFm6 = true
		}
	}
	assert.True(t, hasDm7b5 && hasFm6, "Should detect Dm7b5 and Fm6/D, got: %v", result)
}

func TestDetectAug(t *testing.T) {
	notes := createNotes([]string{"C", "E", "G#"})
	result := Detect(notes)

	hasCaug := false
	hasEaug := false
	hasGSharpAug := false
	for _, chord := range result {
		if chord == "Caug" {
			hasCaug = true
		}
		if chord == "Eaug/C" {
			hasEaug = true
		}
		if chord == "G#aug/C" {
			hasGSharpAug = true
		}
	}
	assert.True(t, hasCaug && hasEaug && hasGSharpAug, "Should detect at least one augmented chord variation, got: %v", result)
}

func TestEdgeCases(t *testing.T) {
	result := Detect([]*note.Note{})
	assert.Empty(t, result, "Should return empty slice for empty input")
}
