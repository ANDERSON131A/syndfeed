package syndfeed

import (
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/antchfx/xmlquery"
)

// SyndFeed is top-level feed object, <feed> in Atom 1.0 and
// <rss> in RSS 2.0.
type SyndFeed struct {
	Authors      []*SyndPerson
	BaseURL      string
	Categories   []string
	Contributors []*SyndPerson
	Copyright    string
	Namespace    map[string]string // map[namespace-prefix]namespace-url
	Description  string
	Generator    string
	Id           string
	ImageURL     string
	Items        []*SyndItem
	Language     string
	// LastUpdatedTime is the feed was last updated time.
	LastUpdatedTime   time.Time
	Title             string
	Links             []*SyndLink
	Version           string
	ElementExtensions []*SyndElementExtension
}

// SyndLink represents a link within a syndication
// feed or item.
type SyndLink struct {
	MediaType string
	URL       string
	Title     string
	RelType   string
}

// SyndItem is a feed item.
type SyndItem struct {
	BaseURL      string
	Authors      []*SyndPerson
	Contributors []*SyndPerson
	Categories   []string
	Content      string
	Copyright    string
	Id           string
	// LastUpdatedTime is the feed item last updated time.
	LastUpdatedTime time.Time
	Links           []*SyndLink
	// PublishDate is the feed item publish date.
	PublishDate       time.Time
	Summary           string
	Title             string
	ElementExtensions []*SyndElementExtension
	//CommentURL      string
}

// SyndPerson is an author or contributor of the feed content.
type SyndPerson struct {
	Name  string
	URL   string
	Email string
}

// SyndElementExtension is an syndication element extension.
type SyndElementExtension struct {
	Name, Namespace, Value string
}

// Parse parses a syndication feed(RSS,Atom).
func Parse(r io.Reader) (*SyndFeed, error) {
	doc, err := xmlquery.Parse(r)
	if err != nil {
		return nil, err
	}
	if doc.SelectElement("rss") != nil {
		return rss.parse(doc)
	} else if doc.SelectElement("feed") != nil {
		return atom.parse(doc)
	}
	return nil, errors.New("invalid syndication feed without <rss> or <feed> element")
}

// LoadURL loads a syndication feed URL.
func LoadURL(url string) (*SyndFeed, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return Parse(res.Body)
}
