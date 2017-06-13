package feedparser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func TestRSS(t *testing.T) {
	// https://en.wikipedia.org/wiki/RSS
	rss := `<?xml version="1.0" encoding="UTF-8" ?>
<rss version="2.0" xmlns:dc="http://purl.org/dc/elements/1.1/" xmlns:content="http://purl.org/rss/1.0/modules/content/">
<channel>
 <title>RSS Title</title>
 <image>
  <url>http://www.blogsmithmedia.com/cn.engadget.com/media/feedlogo.gif?cachebust=true</url>
 </image>
 <description>This is an example of an RSS feed</description>
 <link>http://www.example.com/main.html</link>
 <lastBuildDate>Mon, 06 Sep 2010 00:01:00 +0000 </lastBuildDate>
 <pubDate>Sun, 06 Sep 2009 16:20:00 +0000</pubDate>
 <ttl>1800</ttl>
 <item>
  <title>Example entry</title>
  <description>Here is some text containing an interesting description.</description>
  <link>http://www.example.com/blog/post/1</link>
  <guid isPermaLink="true">7bd204c6-1655-4c27-aeee-53f933c5395f</guid>
  <dc:creator><![CDATA[Andy Yang]]></dc:creator>
  <content:encoded><![CDATA[Here is some text containing an interesting description.]]></content:encoded>
  <pubDate>Sun, 06 Sep 2009 16:20:00 +0000</pubDate>
 </item>
</channel>
</rss>
`
	feed, err := Parse(strings.NewReader(rss))
	if err != nil {
		t.Fatal(err)
	}
	if feed.FeedType != RSSType {
		t.Fatal("expected feed type is not RSS")
	}
	if len(feed.Items) == 0 {
		t.Fatal("feed.Items is nil")
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(feed); err != nil {
		t.Fatal(err)
	}
	fmt.Println(feed.ToString())
}

func TestAtom(t *testing.T) {
	atom := `<?xml version="1.0" encoding="utf-8"?>
<feed xmlns="http://www.w3.org/2005/Atom">
	<title>Example Feed</title>
	<subtitle>A subtitle.</subtitle>
	<link href="http://example.org/feed/" rel="self" />
	<link href="http://example.org/" />
	<id>urn:uuid:60a76c80-d399-11d9-b91C-0003939e0af6</id>
	<updated>2003-12-13T18:30:02Z</updated>		
	<entry>
		<title>Atom-Powered Robots Run Amok</title>
		<link href="http://example.org/2003/12/13/atom03" />
		<link rel="alternate" type="text/html" href="http://example.org/2003/12/13/atom03.html"/>
		<link rel="edit" href="http://example.org/2003/12/13/atom03/edit"/>
		<id>urn:uuid:1225c695-cfb8-4ebb-aaaa-80da344efa6a</id>
		<updated>2003-12-13T18:30:02Z</updated>
		<summary>Some text.</summary>
		<content type="xhtml">
			<div xmlns="http://www.w3.org/1999/xhtml">
				&lt;p&gt;This is the entry content.&lt;/p&gt;
			</div>
		</content>
		<author>
			<name>John Doe</name>
			<email>johndoe@example.com</email>
		</author>
	</entry>
</feed>
	`
	feed, err := Parse(strings.NewReader(atom))
	if err != nil {
		t.Fatal(err)
	}
	if feed.FeedType != AtomType {
		t.Fatal("expected feed type is not Atom")
	}
	if len(feed.Items) == 0 {
		t.Fatal("feed.Items is nil")
	}
}
