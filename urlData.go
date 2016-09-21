package gocrawl

type URLData struct {
	url string
	depth int
}

func InitURLData(url string, depth int) URLData {
	return URLData{url:url, depth:depth}
}
