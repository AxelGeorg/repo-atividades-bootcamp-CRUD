package internal

import "errors"

var (
	ErrRepositoryWarehouseNotFound = errors.New("repository: warehouse not found")
)

type RepositoryWarehouse interface {
	FindAll() ([]Warehouse, error)
	FindById(id int) (w Warehouse, err error)
	Save(w *Warehouse) (err error)
}
