package repository

import (
	"app/internal"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func NewRepositoryWarehouseDB(db *sql.DB) *RepositoryWarehouseDB {
	return &RepositoryWarehouseDB{db: db}
}

type RepositoryWarehouseDB struct {
	db *sql.DB
}

func (r *RepositoryWarehouseDB) FindAll() ([]internal.Warehouse, error) {
	query := "SELECT id, name, address, telephone, capacity FROM warehouses"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var warehouses []internal.Warehouse
	for rows.Next() {
		var w internal.Warehouse
		if err := rows.Scan(&w.Id, &w.Name, &w.Address, &w.Telephone, &w.Capacity); err != nil {
			return nil, err
		}
		warehouses = append(warehouses, w)
	}

	return warehouses, nil
}

func (r *RepositoryWarehouseDB) FindById(id int) (w internal.Warehouse, err error) {
	query := "SELECT id, name, address, telephone, capacity FROM warehouses WHERE id = ?"
	row := r.db.QueryRow(query, id)

	err = row.Scan(&w.Id,
		&w.Name,
		&w.Address,
		&w.Telephone,
		&w.Capacity)

	if err != nil {
		if err == sql.ErrNoRows {
			return w, fmt.Errorf("warehouse not found with id: %d", id)
		}
		return w, err
	}

	return w, nil
}

func (r *RepositoryWarehouseDB) Save(w *internal.Warehouse) (err error) {
	var lastID int

	query := "SELECT COALESCE(MAX(id), 0) FROM warehouses;"
	err = r.db.QueryRow(query).Scan(&lastID)
	if err != nil {
		return err
	}

	w.Id = lastID + 1

	insertQuery := "INSERT INTO warehouses (id, name, address, telephone, capacity) VALUES (?, ?, ?, ?, ?)"

	_, err = r.db.Exec(insertQuery,
		w.Id,
		w.Name,
		w.Address,
		w.Telephone,
		w.Capacity)
	return err
}
