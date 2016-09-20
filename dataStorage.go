package gocrawl

type DataStorage interface {
	/* Store data in the data store */
	Store(data string)
}
