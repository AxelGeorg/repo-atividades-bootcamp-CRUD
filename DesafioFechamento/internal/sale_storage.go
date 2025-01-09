package internal

// StorageSale is the interface that wraps the basic StorageSale methods.
type StorageSale interface {
	// FindAll returns all sales.
	FindAll() (s []Sale, err error)
}
