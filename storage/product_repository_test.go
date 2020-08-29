package storage

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindAllProducts(t *testing.T) {
	ctx := context.Background()
	repo := &ProductRepository{StorageLoc: "./seeds/products.json"}
	params := ProductsParams{Limit: 1}
	products, err := repo.FindAllProducts(ctx, params)
	assert.NoError(t, err)
	assert.Len(t, products, 1)
	assert.Equal(t, "Product name 1", products[0].Name)
	assert.False(t, products[0].CreatedAt.IsZero())
	assert.Equal(t, 10, products[0].Cost)
	assert.Equal(t, []string{"percussion", "bongo"}, products[0].Tags)
}

func TestFindAllProducts_WithTag(t *testing.T) {
	ctx := context.Background()
	repo := &ProductRepository{StorageLoc: "./seeds/products.json"}
	params := ProductsParams{Limit: 10, Tag: "drums"}
	products, err := repo.FindAllProducts(ctx, params)
	assert.NoError(t, err)
	assert.Len(t, products, 1)
	fmt.Printf("the returned product is %v", products)
	assert.Equal(t, "Product name 2", products[0].Name)
	assert.False(t, products[0].CreatedAt.IsZero())
	assert.Equal(t, 2, products[0].Cost)
	assert.Equal(t, []string{"drums", "snares"}, products[0].Tags)
}
