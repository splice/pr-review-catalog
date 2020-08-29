package storage

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindAllProducts(t *testing.T) {
	ctx := context.Background()
	repo := &ProductRepository{StorageLoc: "./seeds/products.json"}
	products, err := repo.FindAllProducts(ctx, 1)
	assert.NoError(t, err)
	assert.Len(t, products, 1)
	assert.Equal(t, "Product name 1", products[0].Name)
	assert.False(t, products[0].CreatedAt.IsZero())
}
