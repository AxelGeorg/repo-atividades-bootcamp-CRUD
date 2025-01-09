package service

import "app/internal"

// NewCustomersDefault creates new default service for customer entity.
func NewCustomersDefault(rp internal.RepositoryCustomer) *CustomersDefault {
	return &CustomersDefault{rp}
}

// CustomersDefault is the default service implementation for customer entity.
type CustomersDefault struct {
	// rp is the repository for customer entity.
	rp internal.RepositoryCustomer
}

// FindAll returns all customers.
func (s *CustomersDefault) FindAll() (c []internal.Customer, err error) {
	c, err = s.rp.FindAll()
	return
}

func (s *CustomersDefault) GetTotalValues() (t []internal.CustomerTotalValue, err error) {
	t, err = s.rp.GetTotalValues()
	return
}

func (s *CustomersDefault) GetSpentMoreMoney() (t []internal.CustomerSpentMoreMoney, err error) {
	t, err = s.rp.GetSpentMoreMoney()
	return
}

// Save saves the customer.
func (s *CustomersDefault) Save(c *internal.Customer) (err error) {
	err = s.rp.Save(c)
	return
}
