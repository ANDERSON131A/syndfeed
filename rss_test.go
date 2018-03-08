package syndfeed

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

// The following example is parse an RSS feed.
func ExampleRSS() {
	f, err := os.Open("./a.rss")
	if err != nil {
		panic(err)
	}
	feed, err := ParseRSS(f)
	if err != nil {
		panic(err)
	}
	fmt.Println(feed.Title)
	fmt.Println(feed.Version)
}

func TestRSSFeed(t *testing.T) {
	f, _ := os.Open("./_samples/engadget.rss")
	feed, err := ParseRSS(f)
	if err != nil {
		t.Fatal(err)
	}
	if e, g := "2.0", feed.Version; g != e {
		t.Fatalf("feed version expected %s but %s", e, g)
	}
	if feed.Title == "" {
		t.Fatalf("<feed.title> is nil")
	}
	if e, g := 25, len(feed.Items); e != g {
		t.Fatalf("item count expected %d but %d", e, g)
	}
	var entry = feed.Items[0]
	if len(entry.Authors) == 0 {
		t.Fatal("<item.author> is nil")
	}
	if e, g := "Mallory Locklear", strings.TrimSpace(entry.Authors[0].Name); e != g {
		t.Fatalf("<item.author> expected %s, but got %s", e, g)
	}
	if entry.Id == "" {
		t.Fatal("<item.id> is nil")
	}
	if e, g := 2, len(entry.ElementExtensions); e != g {
		t.Fatalf("expected extension element number is %d but %d", e, g)
	}
}

func TestRSS2(t *testing.T) {
	f, _ := os.Open("./_samples/rss2.rss")
	feed, err := ParseRSS(f)
	if err != nil {
		t.Fatal(err)
	}
	if len(feed.Namespace) == 0 {
		t.Fatal("feed namespace count is zero")
	}
	if feed.Language == "" {
		t.Fatal("feed.Language is nil")
	}
	if feed.Copyright == "" {
		t.Fatal("feed.Copyright is nil")
	}
	if len(feed.Authors) == 0 {
		t.Fatal("feed.Authors is zero")
	}
}

func TestSyndicationModule(t *testing.T) {

}
