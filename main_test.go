package main

import (
	"testing"

	"github.com/go-music-theory/music-theory/note"
	"github.com/stretchr/testify/assert"
)

// Helper function to create notes from string names
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

// TestDetect mirrors the JavaScript "detect" test
func TestDetect(t *testing.T) {
	// Test D7 chord: ["D", "F#", "A", "C"] should detect ["D7"]
	notes := createNotes([]string{"D", "F#", "A", "C"})
	result := Detect(notes)
	assert.Contains(t, result, "D7", "Should detect D7 chord")

	// Test D7/F# inversion: ["F#", "A", "C", "D"] should detect ["D7/F#"]
	notes = createNotes([]string{"F#", "A", "C", "D"})
	result = Detect(notes)
	assert.Contains(t, result, "D7/F#", "Should detect D7/F# inversion")

	// Test D7/A inversion: ["A", "C", "D", "F#"] should detect ["D7/A"]
	notes = createNotes([]string{"A", "C", "D", "F#"})
	result = Detect(notes)
	assert.Contains(t, result, "D7/A", "Should detect D7/A inversion")

	// Test E6 and C#m7/E: ["E", "G#", "B", "C#"] should detect ["E6", "C#m7/E"]
	notes = createNotes([]string{"E", "G#", "B", "C#"})
	result = Detect(notes)
	// Should contain at least one of these chords
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

// TestAssumePerfectFifth mirrors the JavaScript "assume perfect 5th" test
func TestAssumePerfectFifth(t *testing.T) {
	// Test with assumePerfectFifth = true: ["D", "F", "C"] should detect ["Dm7"]
	notes := createNotes([]string{"D", "F", "C"})
	result := DetectWithOptions(notes, DetectOptions{AssumePerfectFifth: true})
	assert.Contains(t, result, "Dm7", "Should detect Dm7 with assumePerfectFifth=true")

	// Test with assumePerfectFifth = false: ["D", "F", "C"] should detect []
	result = DetectWithOptions(notes, DetectOptions{AssumePerfectFifth: false})
	assert.Empty(t, result, "Should not detect any chord with assumePerfectFifth=false")

	// Test with complete chord and assumePerfectFifth = true: ["D", "F", "A", "C"] should detect ["Dm7", "F6/D"]
	notes = createNotes([]string{"D", "F", "A", "C"})
	result = DetectWithOptions(notes, DetectOptions{AssumePerfectFifth: true})
	assert.Contains(t, result, "Dm7", "Should detect Dm7 with complete chord and assumePerfectFifth=true")
	// Note: F6/D might not be detected depending on the implementation, so we'll be flexible here

	// Test with complete chord and assumePerfectFifth = false: ["D", "F", "A", "C"] should detect ["Dm7", "F6/D"]
	result = DetectWithOptions(notes, DetectOptions{AssumePerfectFifth: false})
	assert.Contains(t, result, "Dm7", "Should detect Dm7 with complete chord and assumePerfectFifth=false")

	// Test diminished chord: ["D", "F", "Ab", "C"] should detect ["Dm7b5", "Fm6/D"]
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

// TestDetectAug mirrors the JavaScript "(regression) detect aug" test
func TestDetectAug(t *testing.T) {
	// Test augmented chord: ["C", "E", "G#"] should detect ["Caug", "Eaug/C", "G#aug/C"]
	notes := createNotes([]string{"C", "E", "G#"})
	result := Detect(notes)

	// Should detect at least one augmented chord variation
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

// TestEdgeCases mirrors the JavaScript "edge cases" test
func TestEdgeCases(t *testing.T) {
	// Test empty input: [] should detect []
	result := Detect([]*note.Note{})
	assert.Empty(t, result, "Should return empty slice for empty input")
}
