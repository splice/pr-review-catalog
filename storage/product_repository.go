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

type ProductsParams struct {
	Limit int
	Tag   string
}

// FindAllProducts returns zero or more products.
func (pr *ProductRepository) FindAllProducts(ctx context.Context, params ProductsParams) ([]*model.Product, error) {
	var products []*model.Product

	absPath, _ := filepath.Abs(pr.StorageLoc)
	data, err := ioutil.ReadFile(absPath)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &products); err != nil {
		return nil, err
	}

	var returnedProducts []*model.Product
	if params.Tag != "" {
		for _, product := range products {
			var shouldReturnProduct bool
			for _, tag := range product.Tags {
				if tag == params.Tag {
					shouldReturnProduct = true
				}
			}
			if shouldReturnProduct {
				returnedProducts = append(returnedProducts, product)
			}
		}
	} else {
		returnedProducts = products
	}

	if len(products) < params.Limit {
		return returnedProducts, nil
	}

	return returnedProducts[:params.Limit], nil
}
