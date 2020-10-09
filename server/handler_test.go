package server

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/splice/catalog-interview/model"
	"github.com/splice/catalog-interview/storage"
	"github.com/stretchr/testify/assert"
)

func TestGetProducts(t *testing.T) {
	productRepo := productRepoMock{}
	findAllProductsMock = func(ctx context.Context, p storage.ProductsParams) ([]*model.Product, error) {
		return []*model.Product{{ID: 1, Name: "abc", Tags: []string{"drums"}}, {ID: 2, Name: "xyz"}}, nil
	}

	controller := NewPageController(&PageControllerParams{Products: productRepo})

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/products", nil)
	controller.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	response := &getProductsPresenter{}
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))

	assert.Len(t, response.Products, 2)
	assert.Equal(t, "abc", response.Products[0].Name)
	assert.Equal(t, "xyz", response.Products[1].Name)
	assert.Equal(t, []string{"drums"}, response.Products[0].Tags)
}

var (
	findAllProductsMock func(context.Context, storage.ProductsParams) ([]*model.Product, error)
)

type productRepoMock struct{}

func (prm productRepoMock) FindAllProducts(ctx context.Context, p storage.ProductsParams) ([]*model.Product, error) {
	return findAllProductsMock(ctx, p)
}
