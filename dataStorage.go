package gocrawl

type DataStorage interface {
	/* Store data associated with a specific url in the data store */
	Store(url string, data string)
}
