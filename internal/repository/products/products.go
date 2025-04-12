package products

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"projectgrom/internal/model"
)

type ProductsDB struct {
	*sql.DB
}

var (
	OpenError   = errors.New("failed to open database")
	AddError    = errors.New("failed to add product")
	UpdateError = errors.New("failed to update product")
)

func NewProducts(data string) (*ProductsDB, error) {
	db, err := sql.Open("postgres", data)
	if err != nil {
		return nil, fmt.Errorf("error to Open db: %w, error:%s", OpenError, err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error to ping db: %s", err)
	}
	return &ProductsDB{db}, nil
}

func (p *ProductsDB) Add(name, desc string, price float64) error {
	data, err := p.Exec("INSERT INTO products(name,description,price) VALUES($1,$2,$3);", name, desc, price)
	if err != nil {
		return fmt.Errorf("error to add product: %w", AddError)
	}
	val, err := data.RowsAffected()
	if err != nil || val != 1 {
		return fmt.Errorf("error to add product: %w", AddError)
	}
	return nil
}

func (p *ProductsDB) GetAll(ctx context.Context) ([]model.Product, error) {
	rows, err := p.Query("SELECT name,description,price FROM products;")
	if err != nil {
		return []model.Product{}, fmt.Errorf("error to get all products: %w", err)
	}
	defer rows.Close()
	var m []model.Product
	var product model.Product
	for rows.Next() {
		if err := rows.Scan(&product.Name, &product.Description, &product.Price); err != nil {
			return []model.Product{}, fmt.Errorf("error to get all products: %w", err)
		}
		m = append(m, product)
	}
	return m, nil
}

func (p *ProductsDB) Update(name string, price float64) error {
	val, err := p.Exec("UPDATE products SET price = $1 WHERE name = $2; ", price, name)
	if err != nil {
		return fmt.Errorf("error to update product: %w", UpdateError)
	}
	effect, err := val.RowsAffected()
	if effect != 1 || err != nil {
		return fmt.Errorf("error to update product: %w", UpdateError)
	}
	return nil
}

func (p *ProductsDB) GetByName(name string) (model.Product, error) {
	row, err := p.Query("SELECT name,description,price FROM products WHERE name=$1;", name)
	if err != nil {
		return model.Product{}, fmt.Errorf("error to get product by name: %w", err)
	}
	defer row.Close()
	var m model.Product
	if err := row.Scan(&m.Name, &m.Description, &m.Price); err != nil {
		return model.Product{}, fmt.Errorf("error to get product by name: %w", err)
	}
	return m, nil
}

func (p *ProductsDB) Delete(name string) error {
	rows, err := p.Exec("DELETE FROM products WHERE name=$1;", name)
	if err != nil {
		return fmt.Errorf("delete error")
	}
	val, err := rows.RowsAffected()
	if err != nil || val != 1 {
		return fmt.Errorf("delete error")
	}
	return nil
}
