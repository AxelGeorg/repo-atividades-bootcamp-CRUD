package repository

import (
	"app/internal"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func NewRepositoryProductDB(db *sql.DB) *RepositoryProductDB {
	return &RepositoryProductDB{db: db}
}

type RepositoryProductDB struct {
	db *sql.DB
}

func (r *RepositoryProductDB) FindById(id int) (p internal.Product, err error) {
	query := "SELECT id, name, quantity, code_value, is_published, expiration, price FROM products WHERE id = ?"
	row := r.db.QueryRow(query, id)

	var isPublishedStr string
	var expirationBytes []byte
	err = row.Scan(&p.Id,
		&p.Name,
		&p.Quantity,
		&p.CodeValue,
		&isPublishedStr,
		&expirationBytes,
		&p.Price)

	if err != nil {
		if err == sql.ErrNoRows {
			return p, fmt.Errorf("product not found with id: %d", id)
		}
		return p, err
	}

	p.IsPublished = (isPublishedStr == "1")

	expirationString := string(expirationBytes)
	expirationTime, err := time.Parse("2006-01-02", expirationString)
	if err != nil {
		return p, fmt.Errorf(err.Error(), id)
	}

	p.Expiration = expirationTime

	return p, nil
}

func (r *RepositoryProductDB) Save(p *internal.Product) (err error) {
	var lastID int

	// Primeiro, buscar o maior ID atual
	query := "SELECT COALESCE(MAX(id), 0) FROM products"
	err = r.db.QueryRow(query).Scan(&lastID)
	if err != nil {
		return err
	}

	// Incrementa o último ID
	p.Id = lastID + 1

	// Prepare o comando de inserção
	insertQuery := "INSERT INTO products (id, name, quantity, code_value, is_published, expiration, price) VALUES (?, ?, ?, ?, ?, ?, ?)"
	isPublishedStr := "0" // padrão para não publicado
	if p.IsPublished {
		isPublishedStr = "1"
	}

	// Inserindo o produto no banco de dados
	_, err = r.db.Exec(insertQuery,
		p.Id, // Agora o ID é definido corretamente
		p.Name,
		p.Quantity,
		p.CodeValue,
		isPublishedStr,
		p.Expiration,
		p.Price)
	return err
}

func (r *RepositoryProductDB) UpdateOrSave(p *internal.Product) (err error) {
	// Check if the product exists
	_, err = r.FindById(p.Id)
	if err != nil {
		// If the product does not exist, save it
		return r.Save(p)
	}
	// Otherwise, update the existing product
	return r.Update(p)
}

func (r *RepositoryProductDB) Update(p *internal.Product) (err error) {
	query := "UPDATE products SET name = ?, quantity = ?, code_value = ?, is_published = ?, expiration = ?, price = ? WHERE id = ?"
	isPublishedStr := "0"
	if p.IsPublished {
		isPublishedStr = "1"
	}

	_, err = r.db.Exec(query,
		p.Name,
		p.Quantity,
		p.CodeValue,
		isPublishedStr,
		p.Expiration,
		p.Price,
		p.Id)
	return
}

func (r *RepositoryProductDB) Delete(id int) (err error) {
	query := "DELETE FROM products WHERE id = ?"
	_, err = r.db.Exec(query, id)
	return
}
