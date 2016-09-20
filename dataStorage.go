package crawler

type DataStorage interface {
	/* Store data in the data store */
	Store()
	Retrieve() string
}
