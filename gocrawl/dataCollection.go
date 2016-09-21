package gocrawl

import (
	"github.com/PuerkitoBio/goquery"
	"time"
)

type DataCollection struct {
	/* the url fetched */
	URL string
	/* the time the url was fetched */
	Fetched time.Time
	/* the parsed data */
	DOM *goquery.Document
}

func InitDataCollection(url string, dom *goquery.Document) DataCollection {
	return DataCollection{URL: url, Fetched: time.Now().UTC(), DOM: dom}
}
