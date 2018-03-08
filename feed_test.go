package syndfeed

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func ExampleParse() {
	s := `<?xml version="1.0" encoding="UTF-8" ?>
	<rss version="2.0">
	<channel>
	 <title>RSS Title</title>
	 <description>This is an example of an RSS feed</description>
	 <link>http://www.example.com/main.html</link>
	 <lastBuildDate>Mon, 06 Sep 2010 00:01:00 +0000 </lastBuildDate>
	 <pubDate>Sun, 06 Sep 2009 16:20:00 +0000</pubDate>
	 <ttl>1800</ttl>	
	 <item>
	  <title>Example entry</title>
	  <description>Here is some text containing an interesting description.</description>
	  <link>http://www.example.com/blog/post/1</link>
	  <guid isPermaLink="false">7bd204c6-1655-4c27-aeee-53f933c5395f</guid>
	  <pubDate>Sun, 06 Sep 2009 16:20:00 +0000</pubDate>
	 </item>	
	</channel>
	</rss>`
	feed, err := Parse(strings.NewReader(s))
	if err != nil {
		panic(err)
	}
	fmt.Println(feed.Title)
	fmt.Println(feed.Version)
	fmt.Println(len(feed.Items))
}

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
