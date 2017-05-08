package feedparser

import (
	"time"

	"github.com/antchfx/xquery/xml"
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
	// Type is type of feed, its RSS or Atom.
	Type Type
	// Title is the name of the channel.
	Title string `json:"title"`
	// Link is the URL to the HTML website corresponding to the channel.
	Link string `json:"link"`
	// Description is sentence describing for channel.
	// atom: subtitle.
	Description string `json:"description"`
	// Category is list of category name.
	Category []string `json:"categories,omitempty"`
	// Copyright notice for content in the channel.
	// atom: rights.
	Copyright string `json:"copyright,omitempty"`
	// Logo is the image logo of channel.
	// rss: image
	// atom: logo
	Logo string `json:"logo,omitempty"`
	// Updated is the last time the content of the channel changed.
	// rss: lastBuildDate
	// atom: updated
	Updated time.Time `json:"pubDate"`
	// Items is a list of article.
	Items []*Item `json:"items"`

	doc *xmlquery.Node
}

// Document returns XML document object.
func (f *Feed) Document() *xmlquery.Node {
	return f.doc
}

// Item represents a single entry in the feed.
type Item struct {
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	Link        string    `json:"link"`
	Author      []string  `json:"authors,omitempty"`
	Category    []string  `json:"categories,omitempty"`
	Published   time.Time `json:"pubDate"`
	Description string    `json:"description"`
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
