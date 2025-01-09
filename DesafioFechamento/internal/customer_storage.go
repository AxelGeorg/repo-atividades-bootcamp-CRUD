package internal

// StorageCustomer is the interface that wraps the basic methods that a customer storage should implement.
type StorageCustomer interface {
	// FindAll returns all customers
	FindAll() (c []Customer, err error)
}
