package server

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/splice/catalog-interview/model"
	"github.com/stretchr/testify/assert"
)

func TestGetProducts(t *testing.T) {
	productRepo := productRepoMock{}
	findAllProductsMock = func(ctx context.Context, limit int) ([]*model.Product, error) {
		return []*model.Product{{ID: 1, Name: "abc"}, {ID: 2, Name: "xyz"}}, nil
	}

	controller := NewPageController(&PageControllerParams{Products: productRepo})

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/products", nil)
	controller.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	response := []*productPresenter{}
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))

	assert.Len(t, response, 2)
	assert.Equal(t, "abc", response[0].Name)
	assert.Equal(t, "xyz", response[1].Name)
}

var (
	findAllProductsMock func(context.Context, int) ([]*model.Product, error)
)

type productRepoMock struct{}

func (prm productRepoMock) FindAllProducts(ctx context.Context, limit int) ([]*model.Product, error) {
	return findAllProductsMock(ctx, limit)
}
