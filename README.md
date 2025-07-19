# go-chord-detector

[![Go Report Card](https://goreportcard.com/badge/github.com/Golevka2001/go-chord-detector)](https://goreportcard.com/report/github.com/Golevka2001/go-chord-detector)
[![go.dev reference](https://godoc.org/github.com/Golevka2001/go-chord-detector?status.svg)](https://pkg.go.dev/github.com/Golevka2001/go-chord-detector)

A Go library for detecting and describing musical chords, inspired by [Tonal.js](https://github.com/tonaljs/tonal) and built on [go-music-theory](https://github.com/go-music-theory/music-theory).

## Usage

```go
import (
    detector "github.com/Golevka2001/go-chord-detector"
    "github.com/go-music-theory/music-theory/note"
)

detector.Detect([]*note.Note{
    note.Named("D"),
    note.Named("F#"),
    note.Named("A"),
    note.Named("C"),
})  // => ["D7"]

// You can also use the note class directly
detector.Detect([]*note.Note{
    {Class: note.E},
    {Class: note.Gs},
    {Class: note.B},
    {Class: note.Cs},
})  // => ["E6", "C#m7/E"]
```

**Options**

- `AssumePerfectFifth`: if `true`, the detector will assume that any chord with a third is also a perfect fifth. This is useful for detecting chords with a missing fifth, but can lead to false positives. Default: `false`.

```go
DetectWithOptions([]*note.Note{
        note.Named("D"),
        note.Named("F"),
        note.Named("C"),
    },
    DetectOptions{AssumePerfectFifth: true},
)  // => ["Dm7"]

DetectWithOptions([]*note.Note{
        note.Named("D"),
        note.Named("F"),
        note.Named("C"),
    },
    DetectOptions{AssumePerfectFifth: false},
)  // => []
```

## License

[MIT License](LICENSE)
