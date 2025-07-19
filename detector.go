package detector

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/Golevka2001/go-chord-detector/chordtype"
	"github.com/Golevka2001/go-chord-detector/pcset"
	"github.com/go-music-theory/music-theory/note"
)

type FoundChord struct {
	Weight float64
	Name   string
}

type DetectOptions struct {
	AssumePerfectFifth bool
}

func Detect(notes []*note.Note) []string {
	return DetectWithOptions(notes, DetectOptions{})
}

func DetectWithOptions(source []*note.Note, options DetectOptions) []string {
	if len(source) == 0 {
		return make([]string, 0)
	}

	found := findMatches(source, 1.0, options)

	var result []string
	for _, chord := range found {
		if chord.Weight > 0 {
			result = append(result, chord.Name)
		}
	}

	sort.Slice(result, func(i, j int) bool {
		var weightI, weightJ float64
		for _, f := range found {
			if f.Name == result[i] {
				weightI = f.Weight
			}
			if f.Name == result[j] {
				weightJ = f.Weight
			}
		}
		return weightI > weightJ
	})

	return result
}

const (
	// 3m 000100000000
	// 3M 000010000000
	AnyThirdsMask = 384 // 0b000110000000 = 384
	// 5P 000000010000
	PerfectFifthMask = 16 // 0b000000010000 = 16
	// 5d 000000100000
	// 5A 000000001000
	NonPerfectFifthsMask = 40 // 0b000000101000 = 40
	// 7m 000000000010
	// 7M 000000000001
	AnySeventhMask = 3 // 0b000000000011 = 3
)

func hasAnyThird(chromaNumber int) bool {
	return chromaNumber&AnyThirdsMask != 0
}

func hasPerfectFifth(chromaNumber int) bool {
	return chromaNumber&PerfectFifthMask != 0
}

func hasAnySeventh(chromaNumber int) bool {
	return chromaNumber&AnySeventhMask != 0
}

func hasNonPerfectFifth(chromaNumber int) bool {
	return chromaNumber&NonPerfectFifthsMask != 0
}

func hasAnyThirdAndPerfectFifthAndAnySeventh(chordType chordtype.ChordType) bool {
	chromaNumber, _ := strconv.ParseInt(chordType.Chroma, 2, 64)
	return hasAnyThird(int(chromaNumber)) && hasPerfectFifth(int(chromaNumber)) && hasAnySeventh(int(chromaNumber))
}

func withPerfectFifth(chroma string) string {
	chromaNumber, _ := strconv.ParseInt(chroma, 2, 64)
	if hasNonPerfectFifth(int(chromaNumber)) {
		return chroma
	}
	newChroma := int(chromaNumber) | PerfectFifthMask
	return chromaToString(newChroma)
}

func chromaToString(chroma int) string {
	binary := make([]byte, 12)
	for i := 0; i < 12; i++ {
		if chroma&(1<<(11-i)) != 0 {
			binary[i] = '1'
		} else {
			binary[i] = '0'
		}
	}
	return string(binary)
}

func findMatches(notes []*note.Note, weight float64, options DetectOptions) []FoundChord {
	if len(notes) == 0 {
		return make([]FoundChord, 0)
	}

	tonic := notes[0]
	tonicChroma := (int(tonic.Class) - 1) % 12

	// We need to test all notes to get the correct baseNote
	allModes := pcset.Modes(notes, false)

	var found []FoundChord
	for index, mode := range allModes {
		modeWithPerfectFifth := mode
		if options.AssumePerfectFifth {
			modeWithPerfectFifth = withPerfectFifth(mode)
		}

		// Some chords could have the same chroma but different interval spelling
		var chordTypes []chordtype.ChordType
		for _, chordType := range chordtype.All() {
			if options.AssumePerfectFifth && hasAnyThirdAndPerfectFifthAndAnySeventh(chordType) {
				if chordType.Chroma == modeWithPerfectFifth {
					chordTypes = append(chordTypes, chordType)
				}
			} else {
				if chordType.Pcset.Chroma == mode {
					chordTypes = append(chordTypes, chordType)
				}
			}
		}

		for _, chordType := range chordTypes {
			chordName := ""
			if len(chordType.Aliases) > 0 {
				chordName = chordType.Aliases[0]
			}

			if index >= int(note.B) {
				continue
			}
			baseNote := note.Class(index + 1).String(note.Sharp) // `0` is defined as `Nil`
			isInversion := index != tonicChroma

			if isInversion {
				found = append(found, FoundChord{
					Weight: 0.5 * weight,
					Name:   fmt.Sprintf("%s%s/%s", baseNote, chordName, tonic.Class.String(note.Sharp)),
				})
			} else {
				found = append(found, FoundChord{
					Weight: 1.0 * weight,
					Name:   fmt.Sprintf("%s%s", baseNote, chordName),
				})
			}
		}
	}

	return found
}
