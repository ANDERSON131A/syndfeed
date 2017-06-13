package feedparser

import (
	"encoding/json"
	"time"
)

// FeedType represents feed types.
type FeedType int

const (
	// AtomType represents feed type is Atom.
	AtomType FeedType = iota
	// RSSType represents feed type is RSS.
	RSSType
)

// Feed is the web feed data for the Atom 1.0 or RSS 2.0.
type Feed struct {
	Title       string     `json:"title,omitempty"`
	Description string     `json:"description,omitempty"`
	Link        string     `json:"link,omitempty"`
	FeedLink    string     `json:"feedLink,omitempty"`
	Updated     *time.Time `json:"updated,omitempty"`
	Published   *time.Time `json:"published,omitempty"`
	Authors     []*Person  `json:"authors,omitempty"`
	Language    string     `json:"language,omitempty"`
	ImageURL    string     `json:"imageUrl,omitempty"`
	Generator   string     `json:"generator,omitempty"`
	Copyright   string     `json:"copyright,omitempty"`
	Categories  []string   `json:"categories,omitempty"`
	FeedType    FeedType   `json:"feedType"`
	FeedVersion string     `json:"feedVersion"`
	Items       []*Item    `json:"items"`
}

func (f Feed) ToString() string {
	json, _ := json.MarshalIndent(f, "", "    ")
	return string(json)
}

// Item represents a single entry in the feed.
type Item struct {
	Title       string       `json:"title,omitempty"`
	Description string       `json:"description,omitempty"`
	Content     string       `json:"content,omitempty"`
	Link        string       `json:"link,omitempty"`
	Updated     *time.Time   `json:"updated,omitempty"`
	Published   *time.Time   `json:"published,omitempty"`
	Authors     []*Person    `json:"author,omitempty"`
	CommentURL  string       `json:"commentUrl,omitempty"`
	GUID        string       `json:"guid,omitempty"`
	ImageURL    string       `json:"imageUrl,omitempty"`
	Categories  []string     `json:"categories,omitempty"`
	Enclosures  []*Enclosure `json:"enclosures,omitempty"`
}

// Enclosure is a file associated with a given Item.
type Enclosure struct {
	URL    string `json:"url,omitempty"`
	Length string `json:"length,omitempty"`
	Type   string `json:"type,omitempty"`
}

// Person is an author or contributor of the feed content.
type Person struct {
	Name  string `json:"name,omitempty"`
	URL   string `json:"url,omitempty"`
	Email string `json:"email,omitempty"`
}
