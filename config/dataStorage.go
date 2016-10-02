package gocrawl

type DataStorage interface {
	/** Store data associated with a specific url in the data store */
	Store(url string, data string)
	/** Get the data associated with a specific url */
	Get(url string) string
	/** Get the number of items in the data store */
	NumItems() int
}
