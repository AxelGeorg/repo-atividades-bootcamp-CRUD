package repository

import (
	"database/sql"
	"log"

	"app/internal"
)

// NewSalesMySQL creates new mysql repository for sale entity.
func NewSalesMySQL(db *sql.DB, storage internal.StorageSale) *SalesMySQL {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM sales").Scan(&count)
	if err != nil && storage != nil {
		log.Printf("Error checking sales table: %v", err)
		return &SalesMySQL{db}
	}

	if count > 0 {
		return &SalesMySQL{db}
	}

	if storage != nil {
		sales, err := storage.FindAll()
		if err == nil {
			for _, sale := range sales {
				_, err := db.Exec(
					"INSERT INTO sales (`quantity`, `product_id`, `invoice_id`) VALUES (?, ?, ?)",
					sale.Quantity, sale.ProductId, sale.InvoiceId,
				)
				if err != nil {
					log.Printf("Error inserting sale %v: %v", sale, err)
				}
			}
		} else {
			log.Printf("Error fetching sales from storage: %v", err)
		}
	}

	return &SalesMySQL{db}
}

// SalesMySQL is the MySQL repository implementation for sale entity.
type SalesMySQL struct {
	// db is the database connection.
	db *sql.DB
}

// FindAll returns all sales from the database.
func (r *SalesMySQL) FindAll() (s []internal.Sale, err error) {
	// execute the query
	rows, err := r.db.Query("SELECT `id`, `quantity`, `product_id`, `invoice_id` FROM sales")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// iterate over the rows
	for rows.Next() {
		var sa internal.Sale
		// scan the row into the sale
		err := rows.Scan(&sa.Id, &sa.Quantity, &sa.ProductId, &sa.InvoiceId)
		if err != nil {
			return nil, err
		}
		// append the sale to the slice
		s = append(s, sa)
	}
	err = rows.Err()
	if err != nil {
		return
	}

	return
}

// Save saves the sale into the database.
func (r *SalesMySQL) Save(s *internal.Sale) (err error) {
	// execute the query
	res, err := r.db.Exec(
		"INSERT INTO sales (`quantity`, `product_id`, `invoice_id`) VALUES (?, ?, ?)",
		(*s).Quantity, (*s).ProductId, (*s).InvoiceId,
	)
	if err != nil {
		return err
	}

	// get the last inserted id
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	// set the id
	(*s).Id = int(id)

	return
}
