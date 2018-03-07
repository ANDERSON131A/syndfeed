package syndfeed

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestLoadURL(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f, _ := os.Open("./_samples/rss2.rss")
		defer f.Close()
		w.Header().Set("Content-Type", "text/xml")
		io.Copy(w, f)
	}))
	defer ts.Close()

	_, err := LoadURL(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
}

func TestInvalidFeed(t *testing.T) {
	var s = `<?xml version="1.0" encoding="utf-8"?>
	<books>
		<book></book>
	</books>
	`
	_, err := Parse(strings.NewReader(s))
	if err == nil {
		t.Fatal("err expected not nil but got nil")
	}
}

func TestParse(t *testing.T) {
	s := `<?xml version="1.0" encoding="utf-8"?>
<feed xmlns="http://www.w3.org/2005/Atom">
	<title>Example Feed</title>
	<subtitle>Insert witty or insightful remark here</subtitle>
	<link href="http://example.org/"/>
	<updated>2003-12-13T18:30:02Z</updated>
	<author>
			<name>John Doe</name>
			<email>johndoe@example.com</email>
	</author>
	<id>urn:uuid:60a76c80-d399-11d9-b93C-0003939e0af6</id>
	<entry>
			<title>Atom-Powered Robots Run Amok</title>
			<link href="http://example.org/2003/12/13/atom03"/>
			<id>urn:uuid:1225c695-cfb8-4ebb-aaaa-80da344efa6a</id>
			<updated>2003-12-13T18:30:02Z</updated>
			<summary>Some text.</summary>
	</entry>
</feed>
`
	feed, err := Parse(strings.NewReader(s))
	if err != nil {
		t.Fatal(err)
	}
	if feed.Title == "" {
		t.Fatal("<feed.title> is nil")
	}
}
