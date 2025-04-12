package products

import (
	"context"
	"projectgrom/internal/model"
	p "projectgrom/internal/repository/products"
)

type ProductsService struct {
	db *p.ProductsDB
}

func InitProductsService(data string) (*ProductsService, error) {
	productDb, err := p.NewProducts(data)
	if err != nil {
		return nil, err
	}
	return &ProductsService{productDb}, nil
}

func (p *ProductsService) Add(name, desc string, price float64) error {
	return p.db.Add(name, desc, price)
}

func (p *ProductsService) Update(name string, price float64) error {
	return p.db.Update(name, price)
}

func (p *ProductsService) GetByName(name string) (model.Product, error) {
	return p.db.GetByName(name)
}

func (p *ProductsService) GetAll(ctx context.Context) ([]model.Product, error) {
	return p.db.GetAll(ctx)
}

func (p *ProductsService) Delete(name string) error {
	return p.db.Delete(name)
}
