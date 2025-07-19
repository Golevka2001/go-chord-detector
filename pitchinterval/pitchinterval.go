// Reference: https://github.com/tonaljs/tonal/tree/main/packages/pitch-interval/index.ts
package pitchinterval

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
)

type Quaility string

const (
	dddd Quaility = "dddd"
	ddd  Quaility = "ddd"
	dd   Quaility = "dd"
	d    Quaility = "d"
	m    Quaility = "m"
	M    Quaility = "M"
	P    Quaility = "P"
	A    Quaility = "A"
	AA   Quaility = "AA"
	AAA  Quaility = "AAA"
	AAAA Quaility = "AAAA"
)

type Type string

const (
	perfectable Type = "perfectable"
	majorable   Type = "majorable"
)

type Interval struct {
	// Pitch
	// NamedPitch
	Empty bool
	Name  string
	Num   int
	Q     Quaility
	T     Type
	Step  int
	Alt   int
	// Dir Direction
	Simple    int
	Semitones int
	Chroma    int
	// Coord IntervalCoordinates
	Oct int
}

var Nointerval = Interval{
	Empty:     true,
	Name:      "",
	Num:       0,
	Q:         "",
	T:         "",
	Step:      0,
	Alt:       0,
	Simple:    0,
	Semitones: 0,
	Chroma:    0,
	Oct:       0,
}

// shorthand tonal notation (with quality after number)
const interbalTonalRegex = `([-+]?\d+)(d{1,4}|m|M|P|A{1,4})`

// standard shorthand notation (with quality before number)
const intervalShorthandRegex = `(AA|A|P|M|m|d|dd)([-+]?\d+)`

var intervalRegex = regexp.MustCompile(`^` + interbalTonalRegex + `|` + intervalShorthandRegex + `$`)

func TokenizeInterval(str string) [2]string {
	matches := intervalRegex.FindStringSubmatch(str)
	if matches == nil {
		return [2]string{"", ""}
	}
	if matches[1] != "" {
		return [2]string{matches[1], matches[2]}
	}
	return [2]string{matches[4], matches[3]}
}

var sizes = []int{0, 2, 4, 5, 7, 9, 11}
var types = []string{"P", "M", "M", "P", "P", "M", "M"}

func Parse(str string) Interval {
	tokens := TokenizeInterval(str)
	if tokens[0] == "" {
		return Nointerval
	}
	num, err := strconv.Atoi(tokens[0])
	if err != nil {
		return Nointerval
	}
	q := Quaility(tokens[1])
	step := int(math.Abs(float64(num))-1) % 7
	t := types[step]
	if t == "M" && q == "P" {
		return Nointerval
	}
	var tp Type
	if t == "M" {
		tp = majorable
	} else {
		tp = perfectable
	}
	name := fmt.Sprintf("%d%s", num, q)
	var dir int
	if num < 0 {
		dir = -1
	} else {
		dir = 1
	}
	var simple int
	if num == 8 || num == -8 {
		simple = num
	} else {
		simple = dir * (step + 1)
	}
	alt := qToAlt(tp, q)
	oct := int(math.Floor((math.Abs(float64(num)) - 1) / 7))
	semitones := dir * (sizes[step] + alt + 12*oct)
	chroma := (((dir * (sizes[step] + alt)) % 12) + 12) % 12
	// coord := pitch.coordinates(pitch.Pitch{
	// 	Step: step,
	// 	Alt:  alt,
	// 	Oct:  oct,
	// 	Dir:  dir,
	// })
	return Interval{
		Empty: false,
		Name:  name,
		Num:   num,
		Q:     q,
		T:     tp,
		Step:  step,
		Alt:   alt,
		// Dir:    dir,
		Simple:    simple,
		Semitones: semitones,
		Chroma:    chroma,
		// Coord:     coord,
		Oct: oct,
	}
}

func qToAlt(t Type, q Quaility) int {
	if (q == "M" && t == majorable) || (q == "P" && t == perfectable) {
		return 0
	}

	if q == "m" && t == majorable {
		return -1
	}

	qStr := string(q)
	if matched, _ := regexp.MatchString("^A+$", qStr); matched {
		return len(qStr)
	}

	if matched, _ := regexp.MatchString("^d+$", qStr); matched {
		if t == perfectable {
			return -len(qStr)
		} else {
			return -(len(qStr) + 1)
		}
	}

	return 0
}
