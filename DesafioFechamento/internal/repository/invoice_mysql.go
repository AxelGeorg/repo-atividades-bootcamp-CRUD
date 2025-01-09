package repository

import (
	"database/sql"
	"log"

	"app/internal"
)

// NewInvoicesMySQL creates new mysql repository for invoice entity.
func NewInvoicesMySQL(db *sql.DB, storage internal.StorageInvoice) *InvoicesMySQL {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM invoices").Scan(&count)
	if err != nil && storage != nil {
		log.Printf("Error checking invoices table: %v", err)
		return &InvoicesMySQL{db}
	}

	if count > 0 {
		return &InvoicesMySQL{db}
	}

	if storage != nil {
		invoices, err := storage.FindAll()
		if err == nil {
			for _, invo := range invoices {
				_, err := db.Exec(
					"INSERT INTO invoices (`datetime`, `total`, `customer_id`) VALUES (?, ?, ?)",
					invo.Datetime, invo.Total, invo.CustomerId,
				)
				if err != nil {
					log.Printf("Error inserting invoice %v: %v", invo, err)
				}
			}
		} else {
			log.Printf("Error fetching invoices from storage: %v", err)
		}
	}

	return &InvoicesMySQL{db}
}

// InvoicesMySQL is the MySQL repository implementation for invoice entity.
type InvoicesMySQL struct {
	// db is the database connection.
	db *sql.DB
}

// FindAll returns all invoices from the database.
func (r *InvoicesMySQL) FindAll() (i []internal.Invoice, err error) {
	// execute the query
	rows, err := r.db.Query("SELECT `id`, `datetime`, `total`, `customer_id` FROM invoices")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// iterate over the rows
	for rows.Next() {
		var iv internal.Invoice
		// scan the row into the invoice
		err := rows.Scan(&iv.Id, &iv.Datetime, &iv.Total, &iv.CustomerId)
		if err != nil {
			return nil, err
		}
		// append the invoice to the slice
		i = append(i, iv)
	}
	err = rows.Err()
	if err != nil {
		return
	}

	return
}

// Save saves the invoice into the database.
func (r *InvoicesMySQL) Save(i *internal.Invoice) (err error) {
	// execute the query
	res, err := r.db.Exec(
		"INSERT INTO invoices (`datetime`, `total`, `customer_id`) VALUES (?, ?, ?)",
		(*i).Datetime, (*i).Total, (*i).CustomerId,
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
	(*i).Id = int(id)

	return
}
