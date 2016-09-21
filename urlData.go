package gocrawl

type UrlData struct {
	url string
	depth int
}

func InitUrlData(url string, depth int) UrlData {
	return UrlData{url:url, depth:depth}
}
