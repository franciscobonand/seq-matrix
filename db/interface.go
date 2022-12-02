package db

type Database interface {
	// Get returns the total number of items registered, and how many of them are valid
	Get() (int64, int64, error)
	// Set creates a new entry on the database, with an unique ID
	Set(seq []string, valid bool) error
}
