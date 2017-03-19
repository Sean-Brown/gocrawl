package gocrawl

/*
UrlData - The URL and the depth of the URL from the starting URL
*/
type UrlData struct {
	URL   string
	Depth int
}

/*
InitUrlData - Construct a new UrlData instance
*/
func InitUrlData(url string, depth int) UrlData {
	return UrlData{URL: url, Depth: depth}
}
