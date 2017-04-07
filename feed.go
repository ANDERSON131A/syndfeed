package feedparser

import (
	"io"
	"time"
)

// Type represents feed types.
type Type int

const (
	// TypeUnknown represents feed type is  unknown.
	TypeUnknown Type = iota
	// TypeAtom represents feed type is Atom.
	TypeAtom
	// TypeRSS represents feed type is RSS.
	TypeRSS
)

func (typ Type) String() string {
	switch typ {
	case TypeAtom:
		return "atom"
	case TypeRSS:
		return "rss"
	default:
		return "unknown"
	}
}

// Feed is represents a RSS(Atom) feed output.
type Feed struct {
	Author      string
	Title       string
	Link        string
	Updated     time.Time
	Image       string
	Language    string
	Description string
	Copyright   string
	Version     string
	Type        Type
	Items       []*Item
}

// Item represents a single entry in a given feed.
type Item struct {
	Title       string
	Body        io.Reader
	Link        string
	Author      []string
	Published   time.Time
	Category    []string
	Description string
}

// ItemSlice provides sorting Item Slice by Published field.
type ItemSlice []*Item

func (slice ItemSlice) Len() int {
	return len(slice)
}

func (slice ItemSlice) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func (slice ItemSlice) Less(i, j int) bool {
	return slice[i].Published.Sub(slice[j].Published) < 0
}
