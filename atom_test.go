package syndfeed

import (
	"fmt"
	"os"
	"testing"

	"github.com/antchfx/xmlquery"
)

// The following example is parse an Atom feed.
func ExampleAtom() {
	f, err := os.Open("./a.atom")
	if err != nil {
		panic(err)
	}
	feed, err := ParseAtom(f)
	if err != nil {
		panic(err)
	}
	fmt.Println(feed.Title)
	fmt.Println(feed.Version)
}

func TestAtomFeed(t *testing.T) {
	f, _ := os.Open("./_samples/atom1.atom")
	feed, err := ParseAtom(f)
	if err != nil {
		t.Fatal(err)
	}
	if feed.Title == "" {
		t.Fatal("<feed.title> is nil")
	}
	if e, g := 2, len(feed.Links); e != g {
		t.Fatalf("<feed.link> element count expected %d but got %d", e, g)
	}
	if len(feed.Items) == 0 {
		t.Fatal("feed.Items count is zero")
	}
}

func TestitunesAtomFeed(t *testing.T) {
	itunesModule := ModuleHandlerFunc(func(n *xmlquery.Node, v interface{}) {
		switch n.Data {
		case "releaseDate":
			t, _ := parseDateString(n.InnerText())
			v.(*Item).PublishDate = t
		case "artist":
			v.(*Item).Authors = append(v.(*Item).Authors, &Person{Name: n.InnerText()})
		case "image":
		}
	})
	RegisterExtensionModule("https://rss.itunes.apple.com", itunesModule)
	f, _ := os.Open("./_samples/itunes.atom")
	feed, err := ParseAtom(f)
	if err != nil {
		t.Fatal(err)
	}
	if e, g := "iTunes Store : Hot Tracks", feed.Title; e != g {
		t.Errorf("<feed.title> expected %s, but got %s", e, g)
	}
	if e, g := "https://rss.itunes.apple.com/api/v1/us/apple-music/hot-tracks/all/10/explicit.atom", feed.Id; e != g {
		t.Errorf("<feed.id> expected %s, but got %s", e, g)
	}
	if e, g := 2, len(feed.Links); e != g {
		t.Errorf("<feed.link> count expected %d but %d", e, g)
	}
	if feed.ImageURL == "" {
		t.Error("feed.ImageURL is nil")
	}
	entry := feed.Items[0]
	if e, g := "Hillsong Worship", entry.Authors[0].Name; e != g {
		t.Errorf("<entry.author> expected %s,but got %s", e, g)
	}
	if entry.Content == "" {
		t.Error("<entry.content> is nil")
	}
}
