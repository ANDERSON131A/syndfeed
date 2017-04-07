FeedParser [![GoDoc](https://godoc.org/github.com/zhengchun/feedparser?status.svg)](https://godoc.org/github.com/zhengchun/feedparser)
===
Parse Atom and RSS feeds.

Examples
===
```go
data := `<rss version="2.0">
<channel>
<title>Sample Feed</title>
</channel>
</rss>`
feed, err := feedparser.Parse(strings.NewReader(data))	
fmt.Println(feed.Title)
fmt.Println(feed.Type)
```

Dependencies
===
[xquery](https://github.com/antchfx/xquery) - Extract data from HTML/XML documents using XPATH.

