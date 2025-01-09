package repository

import (
	"database/sql"
	"log"

	"app/internal"
)

func NewCustomersMySQL(db *sql.DB, storage internal.StorageCustomer) *CustomersMySQL {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM customers").Scan(&count)
	if err != nil && storage != nil {
		log.Printf("Error checking customers table: %v", err)
		return &CustomersMySQL{db}
	}

	if count > 0 {
		return &CustomersMySQL{db}
	}

	if storage != nil {
		customers, err := storage.FindAll()
		if err == nil {
			for _, cust := range customers {
				_, err := db.Exec(
					"INSERT INTO customers (`first_name`, `last_name`, `condition`) VALUES (?, ?, ?)",
					cust.FirstName, cust.LastName, cust.Condition,
				)
				if err != nil {
					log.Printf("Error inserting customer %v: %v", cust, err)
				}
			}
		} else {
			log.Printf("Error fetching customers from storage: %v", err)
		}
	}

	return &CustomersMySQL{db}
}

// CustomersMySQL is the MySQL repository implementation for customer entity.
type CustomersMySQL struct {
	// db is the database connection.
	db *sql.DB
}

// FindAll returns all customers from the database.
func (r *CustomersMySQL) FindAll() (c []internal.Customer, err error) {
	// execute the query
	rows, err := r.db.Query("SELECT `id`, `first_name`, `last_name`, `condition` FROM customers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// iterate over the rows
	for rows.Next() {
		var cs internal.Customer
		// scan the row into the customer
		err := rows.Scan(&cs.Id, &cs.FirstName, &cs.LastName, &cs.Condition)
		if err != nil {
			return nil, err
		}
		// append the customer to the slice
		c = append(c, cs)
	}
	err = rows.Err()
	if err != nil {
		return
	}

	return
}

func (r *CustomersMySQL) GetTotalValues() (totalValues []internal.CustomerTotalValue, err error) {
	rows, err := r.db.Query(`
		SELECT 
			c.condition, 
			ROUND(SUM(s.quantity * p.price), 2) AS total_value
		FROM 
			customers c
		JOIN 
			invoices i ON c.id = i.customer_id
		JOIN 
			sales s ON i.id = s.invoice_id
		JOIN 
			products p ON s.product_id = p.id
		GROUP BY 
			c.condition;
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var tv internal.CustomerTotalValue
		err := rows.Scan(&tv.Condition, &tv.TotalValue)
		if err != nil {
			return nil, err
		}
		totalValues = append(totalValues, tv)
	}

	err = rows.Err()
	if err != nil {
		return
	}

	return
}

func (r *CustomersMySQL) GetSpentMoreMoney() (spentMoreMoney []internal.CustomerSpentMoreMoney, err error) {
	rows, err := r.db.Query(`
		SELECT
		    c.first_name,
		    c.last_name,
		    ROUND(SUM(s.quantity * p.price), 2) AS total_spent
		FROM
		    customers c
		JOIN
		    invoices i ON c.id = i.customer_id
		JOIN
		    sales s ON i.id = s.invoice_id
		JOIN
		    products p ON s.product_id = p.id
		WHERE
    		c.condition = 1 
		GROUP BY
		    c.first_name, c.last_name
		ORDER BY
		    total_spent DESC
		LIMIT 5;
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var cust internal.CustomerSpentMoreMoney
		err := rows.Scan(&cust.FirstName, &cust.LastName, &cust.Amount)
		if err != nil {
			return nil, err
		}
		spentMoreMoney = append(spentMoreMoney, cust)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return spentMoreMoney, nil
}

// Save saves the customer into the database.
func (r *CustomersMySQL) Save(c *internal.Customer) (err error) {
	// execute the query
	res, err := r.db.Exec(
		"INSERT INTO customers (`first_name`, `last_name`, `condition`) VALUES (?, ?, ?)",
		(*c).FirstName, (*c).LastName, (*c).Condition,
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
	(*c).Id = int(id)

	return
}
