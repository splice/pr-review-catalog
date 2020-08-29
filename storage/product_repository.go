package storage

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	"github.com/splice/catalog-interview/model"
)

// ProductRepository provides access to product and related tables.
type ProductRepository struct {
	StorageLoc string
}

// FindAllProducts returns zero or more products.
func (pr *ProductRepository) FindAllProducts(ctx context.Context, limit int) ([]*model.Product, error) {
	var products []*model.Product

	absPath, _ := filepath.Abs(pr.StorageLoc)
	data, err := ioutil.ReadFile(absPath)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &products); err != nil {
		return nil, err
	}

	if len(products) < limit {
		return products, nil
	}

	return products[:limit], nil
}
