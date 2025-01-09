package internal

// StorageProduct is the interface that wraps the basic Product methods.
type StorageProduct interface {
	// FindAll returns all products.
	FindAll() (p []Product, err error)
}
