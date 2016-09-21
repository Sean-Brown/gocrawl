package gocrawl

import (
	"github.com/PuerkitoBio/goquery"
	"time"
)

type DataCollection struct {
	/* the url fetched */
	url string
	/* the time the url was fetched */
	fetched time.Time
	/* the parsed data */
	dom *goquery.Document
}

func InitDataCollection(url string, dom *goquery.Document) DataCollection {
	return DataCollection{url: url, fetched: time.Now().UTC(), dom: dom}
}
