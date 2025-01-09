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

func (r *RepositoryProductDB) FindAll() ([]internal.Product, error) {
	query := "SELECT id, name, quantity, code_value, is_published, expiration, price, id_warehouse FROM products"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []internal.Product
	for rows.Next() {
		var p internal.Product
		var isPublishedStr string
		var expirationBytes []byte

		// Escaneie os dados retornados, incluindo a coluna expiration como []byte
		if err := rows.Scan(&p.Id, &p.Name, &p.Quantity, &p.CodeValue, &isPublishedStr, &expirationBytes, &p.Price, &p.IdWarehouse); err != nil {
			return nil, err
		}

		// Trate a coluna is_published
		p.IsPublished = (isPublishedStr == "1")

		// Converta o expirationBytes para time.Time
		if len(expirationBytes) > 0 {
			expirationString := string(expirationBytes)
			expirationTime, err := time.Parse("2006-01-02", expirationString)
			if err != nil {
				return nil, fmt.Errorf("invalid expiration format for product ID %d: %v", p.Id, err)
			}
			p.Expiration = expirationTime
		}

		products = append(products, p)
	}

	return products, nil
}

func (r *RepositoryProductDB) FindById(id int) (p internal.Product, err error) {
	query := "SELECT id, name, quantity, code_value, is_published, expiration, price, id_warehouse FROM products WHERE id = ?"
	row := r.db.QueryRow(query, id)

	var isPublishedStr string
	var expirationBytes []byte
	err = row.Scan(&p.Id,
		&p.Name,
		&p.Quantity,
		&p.CodeValue,
		&isPublishedStr,
		&expirationBytes,
		&p.Price,
		&p.IdWarehouse)

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

func (r *RepositoryProductDB) CountProductsByWarehouseID(id int) (count int, err error) {
	query := `
        SELECT COUNT(p.id) 
        FROM warehouses w 
        LEFT JOIN products p ON w.id = p.id_warehouse 
        WHERE w.id = ?`

	err = r.db.QueryRow(query, id).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
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
	insertQuery := "INSERT INTO products (id, name, quantity, code_value, is_published, expiration, price, id_warehouse) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
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
		p.Price,
		p.IdWarehouse)
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
	query := "UPDATE products SET name = ?, quantity = ?, code_value = ?, is_published = ?, expiration = ?, price = ?, id_warehouse = ? WHERE id = ?"
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
		p.IdWarehouse,
		p.Id)
	return
}

func (r *RepositoryProductDB) Delete(id int) (err error) {
	query := "DELETE FROM products WHERE id = ?"
	_, err = r.db.Exec(query, id)
	return
}
