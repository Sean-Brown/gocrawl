package gocrawl

import (
	"github.com/PuerkitoBio/goquery"
	"net/url"
	"strings"
	"testing"
	"github.com/Sean-Brown/gocrawl/config"
	"github.com/stretchr/testify/assert"
	"fmt"
)

const UrlConsumerDom = `<html>
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

var _url = &url.URL{Path: "www.a.com/somepage", Host: "a"}

func makeNewDoc(t *testing.T) *goquery.Document {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(UrlConsumerDom))
	if err != nil {
		t.Fatal(err)
	}
	doc.Url = _url
	return doc
}
func initUrlConsumer(sameDomain bool) *UrlConsumer {
	return &UrlConsumer{
		/* make a buffered channel so the go routines won't freeze */
		urls:  make(chan UrlData, 2),
		rules: config.InitURLParsingRules(sameDomain, 10),
	}
}
func assertLinksFound(t *testing.T, urls chan UrlData, expected int) {
	found := len(urls)
	assert.Equal(t, found, expected, fmt.Sprintf("Failed to find all the urls. Expected: %d, Actual: %d", expected, found))
}

func TestParsesAllLinksWhenAllDomainsAreAllowed(t *testing.T) {
	c := initUrlConsumer(false)
	doc := makeNewDoc(t)
	c.parseLinks(doc, 1)
	assertLinksFound(t, c.urls, 2)
}

func TestDoesNotParseLinksWhenOnlyLinksInTheSameDomainAreAllowed(t *testing.T) {
	c := initUrlConsumer(true)
	doc := makeNewDoc(t)
	c.parseLinks(doc, 1)
	assertLinksFound(t, c.urls, 1)
}
