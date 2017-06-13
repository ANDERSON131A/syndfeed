FeedParser 
===
[![GoDoc](https://godoc.org/github.com/zhengchun/feedparser?status.svg)](https://godoc.org/github.com/zhengchun/feedparser)

The `FeedParser` package is a feed parser that supports parsing both RSS and Atom feeds. 

Examples
===
```go
data := `<?xml version="1.0" encoding="UTF-8" ?>
<rss version="2.0" xmlns:dc="http://purl.org/dc/elements/1.1/" >
<channel>
<title>Sample Feed</title>
</channel>
</rss>`
feed, err := feedparser.Parse(strings.NewReader(data))	
fmt.Println(feed.Title)
fmt.Println(feed.FeedType)
```