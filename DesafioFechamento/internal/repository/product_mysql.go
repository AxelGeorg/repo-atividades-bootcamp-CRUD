package repository

import (
	"database/sql"
	"log"

	"app/internal"
)

// NewProductsMySQL creates new mysql repository for product entity.
func NewProductsMySQL(db *sql.DB, storage internal.StorageProduct) *ProductsMySQL {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM products").Scan(&count)
	if err != nil && storage != nil {
		log.Printf("Error checking products table: %v", err)
		return &ProductsMySQL{db}
	}

	if count > 0 {
		return &ProductsMySQL{db}
	}

	if storage != nil {
		products, err := storage.FindAll()
		if err == nil {
			for _, prod := range products {
				_, err := db.Exec(
					"INSERT INTO products (`description`, `price`) VALUES (?, ?)",
					prod.Description, prod.Price,
				)
				if err != nil {
					log.Printf("Error inserting product %v: %v", prod, err)
				}
			}
		} else {
			log.Printf("Error fetching products from storage: %v", err)
		}
	}

	return &ProductsMySQL{db}
}

// ProductsMySQL is the MySQL repository implementation for product entity.
type ProductsMySQL struct {
	// db is the database connection.
	db *sql.DB
}

// FindAll returns all products from the database.
func (r *ProductsMySQL) FindAll() (p []internal.Product, err error) {
	// execute the query
	rows, err := r.db.Query("SELECT `id`, `description`, `price` FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// iterate over the rows
	for rows.Next() {
		var pr internal.Product
		// scan the row into the product
		err := rows.Scan(&pr.Id, &pr.Description, &pr.Price)
		if err != nil {
			return nil, err
		}
		// append the product to the slice
		p = append(p, pr)
	}
	err = rows.Err()
	if err != nil {
		return
	}

	return
}

func (r *ProductsMySQL) GetBestSelling() (p []internal.ProductBestSelling, err error) {
	rows, err := r.db.Query(`
		SELECT 
			p.description, 
			SUM(s.quantity) AS total_sold
		FROM 
			products p
		JOIN 
			sales s ON p.id = s.product_id
		GROUP BY 
			p.id
		ORDER BY 
			total_sold DESC
		LIMIT 5;
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var bestSelling internal.ProductBestSelling
		err := rows.Scan(&bestSelling.Description, &bestSelling.Total)
		if err != nil {
			return nil, err
		}

		p = append(p, bestSelling)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return p, nil
}

// Save saves the product into the database.
func (r *ProductsMySQL) Save(p *internal.Product) (err error) {
	// execute the query
	res, err := r.db.Exec(
		"INSERT INTO products (`description`, `price`) VALUES (?, ?)",
		(*p).Description, (*p).Price,
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
	(*p).Id = int(id)

	return
}
