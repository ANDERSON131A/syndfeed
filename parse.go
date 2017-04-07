package feedparser

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io"
	"strings"
	"time"

	"github.com/antchfx/xquery/xml"
)

var DefaultHandler = HandlerFunc(func(r io.Reader) (*Feed, error) {
	return parse(r)
})

// Handler is an interface provides parsing Atom and RSS feeds.
type Handler interface {
	Parse(io.Reader) (*Feed, error)
}

// HandlerFunc is an adapter to allow the use of ordinary functions as feed handlers.
type HandlerFunc func(io.Reader) (*Feed, error)

func (f HandlerFunc) Parse(r io.Reader) (*Feed, error) {
	return f(r)
}

// Parse parsing a give reader.
func Parse(r io.Reader) (*Feed, error) {
	return DefaultHandler.Parse(r)
}

func parse(r io.Reader) (*Feed, error) {
	preview := make([]byte, 1024)
	n, err := io.ReadFull(r, preview)
	switch {
	case err == io.ErrUnexpectedEOF:
		preview = preview[:n]
	case err != nil:
		return nil, err
	}

	typ := detectFeedType(bytes.NewReader(preview))
	r = io.MultiReader(bytes.NewReader(preview), r)
	switch typ {
	case TypeRSS:
		return parseRSS(r)
	case TypeAtom:
		return parseAtom(r)
	}
	return nil, errors.New("unknown feed type")
}

// https://validator.w3.org/feed/docs/atom.html
func parseAtom(r io.Reader) (*Feed, error) {
	doc, err := xmlquery.Parse(r)
	if err != nil {
		return nil, err
	}
	feed := &Feed{Type: TypeAtom}
	for node := doc.SelectElement("feed").FirstChild; node != nil; node = node.NextSibling {
		name := node.Data
		if node.Prefix != "" {
			name = node.Prefix + ":" + name
		}
		switch name {
		case "title":
			feed.Title = node.InnerText()
		case "link":
			feed.Link = node.SelectAttr("href")
		case "updated":
			if t, err := parseDate(node.InnerText()); err == nil {
				feed.Updated = t
			}
		case "rights":
			feed.Copyright = node.InnerText()
		case "logo":
			feed.Image = node.InnerText()
		case "entry":
			item := &Item{}
			for node := node.FirstChild; node != nil; node = node.NextSibling {
				name := node.Data
				if node.Prefix != "" {
					name = node.Prefix + ":" + name
				}
				switch name {
				case "title":
					item.Title = node.InnerText()
				case "link":
					switch node.SelectAttr("rel") {
					case "", "alternate":
						item.Link = node.SelectAttr("href")
					}
				case "summary":
					item.Description = node.InnerText()
				case "content":
					item.Body = strings.NewReader(node.InnerText())
				case "published":
					if t, err := parseDate(node.InnerText()); err == nil {
						item.Published = t
					}
				case "category":
					item.Category = append(item.Category, node.SelectAttr("term"))
				case "author":
					if node := node.SelectElement("name"); node != nil {
						item.Author = append(item.Author, node.InnerText())
					}
				}
			}
			feed.Items = append(feed.Items, item)
		}
	}
	return feed, nil
}

// https://validator.w3.org/feed/docs/rss2.html
func parseRSS(r io.Reader) (*Feed, error) {
	doc, err := xmlquery.Parse(r)
	if err != nil {
		return nil, err
	}

	version := xmlquery.FindOne(doc, "rss").SelectAttr("version")
	feed := &Feed{Type: TypeRSS, Version: version}
	for node := doc.SelectElement("rss/channel").FirstChild; node != nil; node = node.NextSibling {
		name := node.Data
		if node.Prefix != "" {
			name = node.Prefix + ":" + name
		}
		switch name {
		case "title":
			feed.Title = node.InnerText()
		case "link":
			// `atom:link` score > `link`
			if feed.Link == "" {
				feed.Link = node.InnerText()
			}
		case "atom:link":
			// <atom:link href="" rel="self"></atom:link>
			if node.SelectAttr("rel") == "self" {
				feed.Link = node.SelectAttr("href")
			}
		case "description":
			feed.Description = node.InnerText()
		case "copyright":
			feed.Copyright = node.InnerText()
		case "language":
			feed.Language = node.InnerText()
		case "image":
			if node := node.SelectElement("url"); node != nil {
				feed.Image = node.InnerText()
			}
		case "item":
			item := &Item{}
			for node := node.FirstChild; node != nil; node = node.NextSibling {
				name := node.Data
				if node.Prefix != "" {
					name = node.Prefix + ":" + name
				}
				switch name {
				case "category":
					item.Category = append(item.Category, node.InnerText())
				case "description":
					item.Description = node.InnerText()
				case "link":
					item.Link = node.InnerText()
				case "title":
					item.Title = node.InnerText()
				case "pubDate":
					if t, err := parseDate(node.InnerText()); err == nil {
						item.Published = t
					}
				case "dc:creator", "author":
					item.Author = append(item.Author, node.InnerText())
				case "content:encoded":
					item.Body = strings.NewReader(node.InnerText())
				}
			}
			feed.Items = append(feed.Items, item)
		}
	}
	return feed, nil
}

func detectFeedType(r io.Reader) Type {
	decoder := xml.NewDecoder(r)
Loop:
	for {
		tok, err := decoder.Token()
		switch {
		case err == io.EOF:
			break Loop
		case err != nil:
			break Loop
		}
		switch tok := tok.(type) {
		case xml.StartElement:
			switch tok.Name.Local {
			case "rss":
				return TypeRSS
			case "feed":
				return TypeAtom
			default:
				break Loop
			}
		}
	}
	return TypeUnknown
}

func parseDate(ds string) (t time.Time, err error) {
	d := strings.TrimSpace(ds)
	if d == "" {
		return t, errors.New("invlaid date string")
	}
	for _, f := range dateFormats {
		if t, err = time.Parse(f, d); err == nil {
			return
		}
	}
	return t, errors.New("invalid date format")
}

// DateFormats taken from github.com/mjibson/goread
var dateFormats = []string{
	time.RFC822,  // RSS
	time.RFC822Z, // RSS
	time.RFC3339, // Atom
	time.UnixDate,
	time.RubyDate,
	time.RFC850,
	time.RFC1123Z,
	time.RFC1123,
	time.ANSIC,
	"Mon, January 2 2006 15:04:05 -0700",
	"Mon, January 02, 2006, 15:04:05 MST",
	"Mon, January 02, 2006 15:04:05 MST",
	"Mon, Jan 2, 2006 15:04 MST",
	"Mon, Jan 2 2006 15:04 MST",
	"Mon, Jan 2, 2006 15:04:05 MST",
	"Mon, Jan 2 2006 15:04:05 -700",
	"Mon, Jan 2 2006 15:04:05 -0700",
	"Mon Jan 2 15:04 2006",
	"Mon Jan 2 15:04:05 2006 MST",
	"Mon Jan 02, 2006 3:04 pm",
	"Mon, Jan 02,2006 15:04:05 MST",
	"Mon Jan 02 2006 15:04:05 -0700",
	"Monday, January 2, 2006 15:04:05 MST",
	"Monday, January 2, 2006 03:04 PM",
	"Monday, January 2, 2006",
	"Monday, January 02, 2006",
	"Monday, 2 January 2006 15:04:05 MST",
	"Monday, 2 January 2006 15:04:05 -0700",
	"Monday, 2 Jan 2006 15:04:05 MST",
	"Monday, 2 Jan 2006 15:04:05 -0700",
	"Monday, 02 January 2006 15:04:05 MST",
	"Monday, 02 January 2006 15:04:05 -0700",
	"Monday, 02 January 2006 15:04:05",
	"Mon, 2 January 2006 15:04 MST",
	"Mon, 2 January 2006, 15:04 -0700",
	"Mon, 2 January 2006, 15:04:05 MST",
	"Mon, 2 January 2006 15:04:05 MST",
	"Mon, 2 January 2006 15:04:05 -0700",
	"Mon, 2 January 2006",
	"Mon, 2 Jan 2006 3:04:05 PM -0700",
	"Mon, 2 Jan 2006 15:4:5 MST",
	"Mon, 2 Jan 2006 15:4:5 -0700 GMT",
	"Mon, 2, Jan 2006 15:4",
	"Mon, 2 Jan 2006 15:04 MST",
	"Mon, 2 Jan 2006, 15:04 -0700",
	"Mon, 2 Jan 2006 15:04 -0700",
	"Mon, 2 Jan 2006 15:04:05 UT",
	"Mon, 2 Jan 2006 15:04:05MST",
	"Mon, 2 Jan 2006 15:04:05 MST",
	"Mon 2 Jan 2006 15:04:05 MST",
	"mon,2 Jan 2006 15:04:05 MST",
	"Mon, 2 Jan 2006 15:04:05 -0700 MST",
	"Mon, 2 Jan 2006 15:04:05-0700",
	"Mon, 2 Jan 2006 15:04:05 -0700",
	"Mon, 2 Jan 2006 15:04:05",
	"Mon, 2 Jan 2006 15:04",
	"Mon,2 Jan 2006",
	"Mon, 2 Jan 2006",
	"Mon, 2 Jan 15:04:05 MST",
	"Mon, 2 Jan 06 15:04:05 MST",
	"Mon, 2 Jan 06 15:04:05 -0700",
	"Mon, 2006-01-02 15:04",
	"Mon,02 January 2006 14:04:05 MST",
	"Mon, 02 January 2006",
	"Mon, 02 Jan 2006 3:04:05 PM MST",
	"Mon, 02 Jan 2006 15 -0700",
	"Mon,02 Jan 2006 15:04 MST",
	"Mon, 02 Jan 2006 15:04 MST",
	"Mon, 02 Jan 2006 15:04 -0700",
	"Mon, 02 Jan 2006 15:04:05 Z",
	"Mon, 02 Jan 2006 15:04:05 UT",
	"Mon, 02 Jan 2006 15:04:05 MST-07:00",
	"Mon, 02 Jan 2006 15:04:05 MST -0700",
	"Mon, 02 Jan 2006, 15:04:05 MST",
	"Mon, 02 Jan 2006 15:04:05MST",
	"Mon, 02 Jan 2006 15:04:05 MST",
	"Mon , 02 Jan 2006 15:04:05 MST",
	"Mon, 02 Jan 2006 15:04:05 GMT-0700",
	"Mon,02 Jan 2006 15:04:05 -0700",
	"Mon, 02 Jan 2006 15:04:05 -0700",
	"Mon, 02 Jan 2006 15:04:05 -07:00",
	"Mon, 02 Jan 2006 15:04:05 --0700",
	"Mon 02 Jan 2006 15:04:05 -0700",
	"Mon, 02 Jan 2006 15:04:05 -07",
	"Mon, 02 Jan 2006 15:04:05 00",
	"Mon, 02 Jan 2006 15:04:05",
	"Mon, 02 Jan 2006",
	"Mon, 02 Jan 06 15:04:05 MST",
	"January 2, 2006 3:04 PM",
	"January 2, 2006, 3:04 p.m.",
	"January 2, 2006 15:04:05 MST",
	"January 2, 2006 15:04:05",
	"January 2, 2006 03:04 PM",
	"January 2, 2006",
	"January 02, 2006 15:04:05 MST",
	"January 02, 2006 15:04",
	"January 02, 2006 03:04 PM",
	"January 02, 2006",
	"Jan 2, 2006 3:04:05 PM MST",
	"Jan 2, 2006 3:04:05 PM",
	"Jan 2, 2006 15:04:05 MST",
	"Jan 2, 2006",
	"Jan 02 2006 03:04:05PM",
	"Jan 02, 2006",
	"6/1/2 15:04",
	"6-1-2 15:04",
	"2 January 2006 15:04:05 MST",
	"2 January 2006 15:04:05 -0700",
	"2 January 2006",
	"2 Jan 2006 15:04:05 Z",
	"2 Jan 2006 15:04:05 MST",
	"2 Jan 2006 15:04:05 -0700",
	"2 Jan 2006",
	"2.1.2006 15:04:05",
	"2/1/2006",
	"2-1-2006",
	"2006 January 02",
	"2006-1-2T15:04:05Z",
	"2006-1-2 15:04:05",
	"2006-1-2",
	"2006-1-02T15:04:05Z",
	"2006-01-02T15:04Z",
	"2006-01-02T15:04-07:00",
	"2006-01-02T15:04:05Z",
	"2006-01-02T15:04:05-07:00:00",
	"2006-01-02T15:04:05:-0700",
	"2006-01-02T15:04:05-0700",
	"2006-01-02T15:04:05-07:00",
	"2006-01-02T15:04:05 -0700",
	"2006-01-02T15:04:05:00",
	"2006-01-02T15:04:05",
	"2006-01-02 at 15:04:05",
	"2006-01-02 15:04:05Z",
	"2006-01-02 15:04:05 MST",
	"2006-01-02 15:04:05-0700",
	"2006-01-02 15:04:05-07:00",
	"2006-01-02 15:04:05 -0700",
	"2006-01-02 15:04",
	"2006-01-02 00:00:00.0 15:04:05.0 -0700",
	"2006/01/02",
	"2006-01-02",
	"15:04 02.01.2006 -0700",
	"1/2/2006 3:04:05 PM MST",
	"1/2/2006 3:04:05 PM",
	"1/2/2006 15:04:05 MST",
	"1/2/2006",
	"06/1/2 15:04",
	"06-1-2 15:04",
	"02 Monday, Jan 2006 15:04",
	"02 Jan 2006 15:04 MST",
	"02 Jan 2006 15:04:05 UT",
	"02 Jan 2006 15:04:05 MST",
	"02 Jan 2006 15:04:05 -0700",
	"02 Jan 2006 15:04:05",
	"02 Jan 2006",
	"02/01/2006 15:04 MST",
	"02-01-2006 15:04:05 MST",
	"02.01.2006 15:04:05",
	"02/01/2006 15:04:05",
	"02.01.2006 15:04",
	"02/01/2006 - 15:04",
	"02.01.2006 -0700",
	"02/01/2006",
	"02-01-2006",
	"01/02/2006 3:04 PM",
	"01/02/2006 15:04:05 MST",
	"01/02/2006 - 15:04",
	"01/02/2006",
	"01-02-2006",
}
