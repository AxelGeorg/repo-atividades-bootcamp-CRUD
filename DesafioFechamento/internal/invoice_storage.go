package internal

// StorageInvoice is the interface that wraps the basic invoice methods.
type StorageInvoice interface {
	// FindAll returns all invoices.
	FindAll() (i []Invoice, err error)
}
