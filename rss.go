package feedparser

import (
	"io"

	"github.com/antchfx/xquery/xml"
)

// parsing an RSS reader.
func parseRSS(r io.Reader) (*Feed, error) {
	doc, err := xmlquery.Parse(r)
	if err != nil {
		return nil, err
	}

	// extension modules for this feed.
	var extensions = make(map[string]string)
	root := doc.SelectElement("rss")
	for _, attr := range root.Attr {
		if attr.Name.Space == "xmlns" {
			extensions[attr.Name.Local] = attr.Value
		}
	}

	feed := &Feed{FeedType: RSSType}
	feed.FeedVersion = root.SelectAttr("version")
	for elem := root.SelectElement("channel").FirstChild; elem != nil; elem = elem.NextSibling {
		if v, ok := extensions[elem.Prefix]; ok {
			handler, ok := modules[v]
			if ok {
				handler.ParseElement(elem, feed, FeedElement)
			}
			continue
		}
		switch elem.Data {
		case "title":
			feed.Title = elem.InnerText()
		case "description":
			feed.Description = elem.InnerText()
		case "link":
			feed.Link = elem.InnerText()
		case "category":
			feed.Categories = append(feed.Categories, elem.InnerText())
		case "copyright":
			feed.Copyright = elem.InnerText()
		case "generator":
			feed.Generator = elem.InnerText()
		case "image":
			if elem := elem.SelectElement("url"); elem != nil {
				feed.ImageURL = elem.InnerText()
			}
		case "language":
			feed.Language = elem.InnerText()
		case "lastBuildDate":
			if t, err := ParseDate(elem.InnerText()); err == nil {
				feed.Updated = &t
			}
		case "pubDate":
			if t, err := ParseDate(elem.InnerText()); err == nil {
				feed.Published = &t
			}
		case "item":
			item := new(Item)
			for elem := elem.FirstChild; elem != nil; elem = elem.NextSibling {
				if v, ok := extensions[elem.Prefix]; ok {
					handler, ok := modules[v]
					if ok {
						handler.ParseElement(elem, item, ItemElement)
					}
					continue
				}
				switch elem.Data {
				case "author":
					item.Authors = append(item.Authors, &Person{Name: elem.InnerText()})
				case "category":
					item.Categories = append(item.Categories, elem.InnerText())
				case "comments":
					item.CommentURL = elem.InnerText()
				case "description":
					item.Description = elem.InnerText()
				case "enclosure":
					item.Enclosures = append(item.Enclosures, &Enclosure{
						URL:    elem.SelectAttr("url"),
						Type:   elem.SelectAttr("type"),
						Length: elem.SelectAttr("length"),
					})
				case "guid":
					item.GUID = elem.InnerText()
				case "link":
					item.Link = elem.InnerText()
				case "pubDate":
					if t, err := ParseDate(elem.InnerText()); err == nil {
						item.Published = &t
					}
				case "title":
					item.Title = elem.InnerText()
				}
			}
			feed.Items = append(feed.Items, item)
		}
	}
	return feed, nil
}
