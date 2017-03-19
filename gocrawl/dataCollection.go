package gocrawl

import (
	"time"

	"github.com/PuerkitoBio/goquery"
)

/*
DomQuery - The URL and the DOM elements that were parsed from that site
*/
type DomQuery struct {
	/* the url fetched */
	URL string
	/* the time the url was fetched */
	Fetched time.Time
	/* the parsed data */
	DOM *goquery.Document
}

/*
InitDomQuery - DomQuery constructor
*/
func InitDomQuery(url string, dom *goquery.Document) DomQuery {
	return DomQuery{URL: url, Fetched: time.Now().UTC(), DOM: dom}
}
