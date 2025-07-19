// A dictionary of musical chords.
// Reference: https://github.com/tonaljs/tonal/tree/main/packages/chord-type/index.ts
package chordtype

import (
	"sort"
	"strconv"
	"strings"

	"github.com/Golevka2001/go-chord-detector/pcset"
)

type ChordQuality string

const (
	Major      ChordQuality = "Major"
	Minor      ChordQuality = "Minor"
	Augmented  ChordQuality = "Augmented"
	Diminished ChordQuality = "Diminished"
	Unknown    ChordQuality = "Unknown"
)

type ChordType struct {
	pcset.Pcset
	Name      string
	Quality   ChordQuality
	Aliases   []string
	Intervals []string
}

var NoChordType = ChordType{
	Pcset:     pcset.EmptyPcset,
	Name:      "",
	Quality:   Unknown,
	Intervals: []string{},
	Aliases:   []string{},
}

var dictionary []ChordType
var index map[string]ChordType

func init() {
	dictionary = make([]ChordType, 0)
	index = make(map[string]ChordType)

	for _, data := range chords {
		if len(data) >= 3 {
			intervals := strings.Split(data[0], " ")
			fullName := data[1]
			aliases := strings.Split(data[2], " ")
			Add(intervals, aliases, fullName)
		}
	}

	sort.Slice(dictionary, func(i, j int) bool {
		return dictionary[i].Pcset.SetNum < dictionary[j].Pcset.SetNum
	})
}

// Get retrieves a chord type by name, alias, chroma, or setNum.
func Get(typeName string) ChordType {
	if chord, exists := index[typeName]; exists {
		return chord
	}
	return NoChordType
}

// Names returns all chord (long) names.
func Names() []string {
	var names []string
	for _, chord := range dictionary {
		if chord.Name != "" {
			names = append(names, chord.Name)
		}
	}
	return names
}

// Symbols returns all chord symbols.
func Symbols() []string {
	var symbols []string
	for _, chord := range dictionary {
		if len(chord.Aliases) > 0 {
			symbols = append(symbols, chord.Aliases[0])
		}
	}
	return symbols
}

// Keys returns all the keys used to reference chord types
func Keys() []string {
	var keys []string
	for key := range index {
		keys = append(keys, key)
	}
	return keys
}

// All return a list of all chord types.
func All() []ChordType {
	return dictionary
}

// RemoveAll clears the dictionary and index.
func RemoveAll() {
	dictionary = make([]ChordType, 0)
	index = make(map[string]ChordType)
}

// Add adds a chord to the dictionary.
func Add(intervals []string, aliases []string, fullName string) {
	quality := getQuality(intervals)

	chord := ChordType{
		Pcset:     pcset.IntervalsToPcset(intervals),
		Name:      fullName,
		Quality:   quality,
		Intervals: intervals,
		Aliases:   aliases,
	}

	dictionary = append(dictionary, chord)
	if chord.Name != "" {
		index[chord.Name] = chord
	}
	index[strconv.Itoa(chord.Pcset.SetNum)] = chord
	index[chord.Pcset.Chroma] = chord

	for _, alias := range chord.Aliases {
		AddAlias(chord, alias)
	}
}

func AddAlias(chord ChordType, alias string) {
	index[alias] = chord
}

func getQuality(intervals []string) ChordQuality {
	for _, interval := range intervals {
		switch interval {
		case "5A":
			return Augmented
		case "3M":
			return Major
		case "5d":
			return Diminished
		case "3m":
			return Minor
		}
	}
	return Unknown
}
