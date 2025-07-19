// A pitch class set is a set (no repeated) of pitch classes (notes without octaves).
// Pitch classes are useful to identify musical structures (if two chords are related, for example).
// Reference: https://github.com/tonaljs/tonal/tree/main/packages/pcset/index.ts
package pcset

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Golevka2001/go-chord-detector/pitchinterval"

	"github.com/go-music-theory/music-theory/note"
)

// Pcset defines the properties of a pitch class set.
//
// SetNum is a number between 1 and 4095 (both included) that uniquely identifies
// the set. It's the decimal number of the chroma.
//
// Chroma is a string representation of the set: a 12-char string with either "1"
// or "0" as characters, representing a pitch class or not for the given position
// in the octave. For example, a "1" at index 0 means 'C', a "1" at index 2 means 'D', and so on.
//
// Normalized is the chroma but shifted to the first 1.
//
// Intervals are the intervals of the pitch class set *starting from C*.
type Pcset struct {
	Empty      bool
	Name       string
	SetNum     int
	Chroma     string
	Normalized string
	Intervals  []string
}

var EmptyPcset = Pcset{
	Empty:      true,
	Name:       "",
	SetNum:     0,
	Chroma:     "000000000000",
	Normalized: "000000000000",
	Intervals:  []string{},
}

func setNumToChroma(num int) string {
	binary := fmt.Sprintf("%012b", num)
	return binary
}

func chromaToNumber(chroma string) int {
	num, _ := strconv.ParseInt(chroma, 2, 64)
	return int(num)
}

var cache = map[string]Pcset{
	EmptyPcset.Chroma: EmptyPcset,
}

// NotesToPcset replaces the original `get(src: Set)` method.
func NotesToPcset(set []*note.Note) Pcset {
	chroma := notesToChroma(set)

	if cached, exists := cache[chroma]; exists {
		return cached
	}

	pcset := chromaToPcset(chroma)
	cache[chroma] = pcset
	return pcset
}

// IntervalsToPcset replaces the original `get(src: string[])` method.
func IntervalsToPcset(set []string) Pcset {
	chroma := intervalsToChroma(set)

	if cached, exists := cache[chroma]; exists {
		return cached
	}

	pcset := chromaToPcset(chroma)
	cache[chroma] = pcset
	return pcset
}

var ivls = []string{
	"1P", "2m", "2M", "3m", "3M", "4P", "5d", "5P", "6m", "6M", "7m", "7M",
}

// chromaToIntervals returns the intervals of a pcset *starting from C*.
//
// Parameters:
//   - chroma: the chroma of the pcset.
//
// Returns an array of interval names or an empty array if not a valid pitch class set.
func chromaToIntervals(chroma string) []string {
	var result []string
	for i, char := range chroma {
		if char == '1' {
			result = append(result, ivls[i])
		}
	}
	return result
}

// Modes returns all rotations of the chroma, optionally discarding ones that start with "0".
// This is used, for example, to get all the modes of a scale.
//
// Paramerters:
//   - set: the list of notes or pitchChr of the set.
//   - normalize: remove all the rotations that starts with "0".
//
// Returns an array with all the modes of the chroma.
func Modes(set []*note.Note, normalize bool) []string {
	pcset := NotesToPcset(set)
	binary := strings.Split(pcset.Chroma, "")

	var modes []string
	for i := 0; i < len(binary); i++ {
		rotated := rotate(binary, i)
		rotatedStr := strings.Join(rotated, "")

		if !normalize || rotatedStr[0] == '1' {
			modes = append(modes, rotatedStr)
		}
	}
	return compact(modes)
}

func chromaToPcset(chroma string) Pcset {
	setNum := chromaToNumber(chroma)

	rotations := chromaRotations(chroma)
	normalizedNum := setNum
	for _, rotation := range rotations {
		if rotation[0] == '1' {
			num := chromaToNumber(rotation)
			if num < normalizedNum || normalizedNum < 2048 {
				normalizedNum = num
			}
		}
	}
	if normalizedNum < 2048 {
		normalizedNum = setNum
	}
	normalized := setNumToChroma(normalizedNum)

	intervals := chromaToIntervals(chroma)

	return Pcset{
		Empty:      false,
		Name:       "",
		SetNum:     setNum,
		Chroma:     chroma,
		Normalized: normalized,
		Intervals:  intervals,
	}
}

func chromaRotations(chroma string) []string {
	binary := strings.Split(chroma, "")
	var rotations []string

	for i := 0; i < len(binary); i++ {
		rotated := rotate(binary, i)
		rotations = append(rotations, strings.Join(rotated, ""))
	}
	return rotations
}

// notesToChroma replaces the original `listToChroma(set: any[])` method.
func notesToChroma(set []*note.Note) string {
	if len(set) == 0 {
		return EmptyPcset.Chroma
	}

	binary := [12]byte{'0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0'}
	for _, n := range set {
		if n.Class != note.Nil {
			index := (int(n.Class) - 1) % 12
			binary[index] = '1'
		}
	}

	return string(binary[:])
}

// intervalsToChroma replaces the original `intervalsToChroma(set: string[])` method.
func intervalsToChroma(set []string) string {
	if len(set) == 0 {
		return EmptyPcset.Chroma
	}

	binary := [12]byte{'0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '0'}
	for _, interval := range set {
		index := pitchinterval.Parse(interval).Chroma
		if index != -1 {
			binary[index] = '1'
		}
	}
	return string(binary[:])
}

func rotate(slice []string, n int) []string {
	if len(slice) == 0 {
		return slice
	}
	length := len(slice)
	n = ((n % length) + length) % length
	result := make([]string, length)
	copy(result, slice[n:])
	copy(result[length-n:], slice[:n])
	return result
}

func compact(arr []string) []string {
	result := make([]string, 0)
	for _, item := range arr {
		if item != "" {
			result = append(result, item)
		}
	}
	return result
}
