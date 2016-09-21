package gocrawl

type URLData struct {
	URL   string
	Depth int
}

func InitURLData(url string, depth int) URLData {
	return URLData{URL: url, Depth: depth}
}
