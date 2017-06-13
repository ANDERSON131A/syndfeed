package feedparser

import (
	"html"
	"io"

	"github.com/antchfx/xquery/xml"
)

func parseAtomAuthorElement(elem *xmlquery.Node) *Person {
	author := new(Person)
	for elem := elem.FirstChild; elem != nil; elem = elem.NextSibling {
		switch elem.Data {
		case "name":
			author.Name = elem.InnerText()
		case "email":
			author.Email = elem.InnerText()
		case "uri":
			author.URL = elem.InnerText()
		}
	}
	return author
}

func parseAtom(r io.Reader) (*Feed, error) {
	doc, err := xmlquery.Parse(r)
	if err != nil {
		return nil, err
	}

	root := doc.SelectElement("feed")
	feed := &Feed{FeedType: AtomType}
	feed.FeedVersion = root.SelectAttr("version")

	for elem := root.FirstChild; elem != nil; elem = elem.NextSibling {
		switch elem.Data {
		case "title":
			feed.Title = elem.InnerText()
		case "updated":
			if t, err := ParseDate(elem.InnerText()); err == nil {
				feed.Updated = &t
			}
		case "author":
			feed.Authors = append(feed.Authors, parseAtomAuthorElement(elem))
		case "link":
			if elem.SelectAttr("rel") == "self" {
				feed.FeedLink = elem.SelectAttr("href")
			} else {
				feed.Link = elem.SelectAttr("href")
			}
		case "category":
			if term := elem.SelectAttr("category"); term != "" {
				feed.Categories = append(feed.Categories, term)
			}
		case "generator":
			feed.Generator = elem.InnerText()
		case "logo":
			feed.ImageURL = elem.InnerText()
		case "rights":
			feed.Copyright = elem.InnerText()
		case "subtitle":
			feed.Description = elem.InnerText()
		case "entry":
			item := new(Item)
			for elem := elem.FirstChild; elem != nil; elem = elem.NextSibling {
				switch elem.Data {
				case "id":
					item.GUID = elem.InnerText()
				case "title":
					item.Title = elem.InnerText()
				case "updated":
					if t, err := ParseDate(elem.InnerText()); err == nil {
						item.Updated = &t
					}
				case "author":
					item.Authors = append(item.Authors, parseAtomAuthorElement(elem))
				case "content":
					item.Content = html.UnescapeString(elem.InnerText())
				case "link":
					if elem.SelectAttr("rel") == "alternate" {
						item.Link = elem.SelectAttr("href")
					} else if item.Link == "" {
						item.Link = elem.SelectAttr("href")
					}
				case "summary":
					item.Description = html.UnescapeString(elem.InnerText())
				}
			}
			feed.Items = append(feed.Items, item)
		}
	}
	return feed, nil
}
