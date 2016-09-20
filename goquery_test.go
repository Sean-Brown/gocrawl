package crawler

import (
	"testing"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

const dom = `<html>
<head>
<title>woooo</title>
<meta a="b">
<link rel="stylesheet" type="text/css" href="theme.css">
<base href="www.base.com">
</head>
<body>
<div>
<p>Some text <a href="www.a.com">ref</a></p>
<h1 class="header">header1</h1>
<h2 class="header">header2</h2>
</div>
</html>`

func Disabled_TestFindingLinks(t *testing.T) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(dom))
	if err != nil {
		t.Fatal(err)
	}
	atag := doc.Find(a)
	atag.Each(func(_ int, tag *goquery.Selection) {
		attr, exists := tag.Attr(href)
		if exists {
			t.Log(attr)
		}
	})
}
