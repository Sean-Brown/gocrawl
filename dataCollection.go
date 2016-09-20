package gocrawl

import (
	"time"
	"github.com/PuerkitoBio/goquery"
)

type DataCollection struct {
	/* the url fetched */
	url string
	/* the time the url was fetched */
	fetched time.Time
	/* the parsed data */
	dom *goquery.Document
}

func Init(url string, dom *goquery.Document) DataCollection {
	return DataCollection{url:url, fetched:time.Now().UTC(), dom:dom}
}
