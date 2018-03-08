package syndfeed

import "github.com/antchfx/xmlquery"

// http://www.disobey.com/detergent/2002/extendingrss2/

// Module is an interface is used to handle an extension element
// which not specified in either the Atom 1.0 or RSS 2.0 specifications.
type Module interface {
	ParseElement(n *xmlquery.Node, v interface{})
}

// ModuleHandlerFunc is a utility function that wrapper a function
// as Module object.
type ModuleHandlerFunc func(*xmlquery.Node, interface{})

func (f ModuleHandlerFunc) ParseElement(n *xmlquery.Node, v interface{}) {
	f(n, v)
}

// http://web.resource.org/rss/1.0/modules/content/
var rssContentModule = ModuleHandlerFunc(func(n *xmlquery.Node, v interface{}) {
	if n.Data == "encoded" {
		v.(*Item).Content = n.InnerText()
	}
})

// http://web.resource.org/rss/1.0/modules/syndication/
var rssSyndicationModule = ModuleHandlerFunc(func(n *xmlquery.Node, v interface{}) {
	if n.Data == "updateBase" {
		if t, err := parseDateString(n.InnerText()); err == nil {
			v.(*Feed).LastUpdatedTime = t
		}
	}
})

// http://web.resource.org/rss/1.0/modules/dc/
var rssDublinCoreModule = ModuleHandlerFunc(func(n *xmlquery.Node, v interface{}) {
	switch v2 := v.(type) {
	case *Feed:
		switch n.Data {
		case "title":
			v2.Title = n.InnerText()
		case "creator":
			v2.Authors = append(v2.Authors, &Person{Name: n.InnerText()})
		case "subject":
		case "description":
			v2.Description = n.InnerText()
		case "publisher":
		case "contributor":
			v2.Contributors = append(v2.Contributors, &Person{Name: n.InnerText()})
		case "date":
			if t, err := parseDateString(n.InnerText()); err == nil {
				v2.LastUpdatedTime = t
			}
		case "language":
			v2.Language = n.InnerText()
		case "rights":
			v2.Copyright = n.InnerText()
		}
	case *Item:
		switch n.Data {
		case "title":
			v2.Title = n.InnerText()
		case "creator":
			v2.Authors = append(v2.Authors, &Person{Name: n.InnerText()})
		case "subject":
		case "description":
			v2.Summary = n.InnerText()
		case "publisher":
		case "contributor":
			v2.Contributors = append(v2.Contributors, &Person{Name: n.InnerText()})
		case "date":
			if t, err := parseDateString(n.InnerText()); err == nil {
				v2.PublishDate = t
			}
		case "rights":
			v2.Copyright = n.InnerText()
		}
	}
})

var modules = make(map[string]Module)

// RegisterExtensionModule registers Module with the specified XML namespace.
func RegisterExtensionModule(nsURL string, m Module) {
	modules[nsURL] = m
}

func init() {
	RegisterExtensionModule("http://purl.org/dc/elements/1.1/", rssDublinCoreModule)
	RegisterExtensionModule("https://purl.org/dc/elements/1.1/", rssDublinCoreModule)
	RegisterExtensionModule("http://purl.org/rss/1.0/modules/content/", rssContentModule)
	RegisterExtensionModule("https://purl.org/rss/1.0/modules/content/", rssContentModule)
	RegisterExtensionModule("http://purl.org/rss/1.0/modules/syndication/", rssSyndicationModule)
	RegisterExtensionModule("https://purl.org/rss/1.0/modules/syndication/", rssSyndicationModule)
}
