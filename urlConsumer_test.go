package crawler

import (
	"testing"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"net/url"
)

const urlConsumerDom = `<html>
<head>
<title>woooo</title>
<meta a="b">
<link rel="stylesheet" type="text/css" href="theme.css">
<base href="www.base.com">
</head>
<body>
<div>
<p>Some text <a href="www.a.com">ref</a> <a href="www.b.com">b</a></p>
<h1 class="header">header1</h1>
<h2 class="header">header2</h2>
</div>
</html>`
var _url = &url.URL{Path: "www.a.com/somepage"}

func makeNewDoc(t *testing.T) *goquery.Document {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(urlConsumerDom))
	if err != nil {
		t.Fatal(err)
	}
	doc.Url = _url
	return doc
}
func initURLConsumer(sameDomain bool) *URLConsumer {
	return &URLConsumer{
		/* make a buffered channel so the go routines won't freeze */
		urls: make(chan string, 2),
		rules: InitURLParsingRules(sameDomain),
	}
}
func assertLinksFound(t *testing.T, urls chan string, expected int) {
	found := len(urls)
	if found != expected {
		t.Fatal("Failed to find all the urls. Expected: ", expected, ", Actual: ", found)
	}
}

func TestParsesAllLinksWhenAllDomainsAreAllowed(t *testing.T) {
	c := initURLConsumer(false)
	doc := makeNewDoc(t)
	c.parseLinks(doc)
	assertLinksFound(t, c.urls, 2)
}

func TestDoesNotParseLinksWhenOnlyLinksInTheSameDomainAreAllowed(t *testing.T) {
	c := initURLConsumer(true)
	doc := makeNewDoc(t)
	c.parseLinks(doc)
	assertLinksFound(t, c.urls, 1)
}
