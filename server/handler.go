package server

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/pkg/errors"
	"github.com/splice/catalog-interview/model"
)

const (
	maxLimit = 10
)

// A pageController handles dispatch of HTTP requests to handlers.
type pageController struct {
	mux      *mux.Router
	dec      *schema.Decoder // safe to share across requests
	products productRepository
}

// productRepository is implemented by any type that can provide access
// to products and their related tables.
type productRepository interface {
	FindAllProducts(ctx context.Context, limit int) ([]*model.Product, error)
}

type getProductsArgs struct {
	// Limit is the number of items per page
	Limit int `schema:"limit"`
}

// getProducts returns a paginated list of products.
// getAssets returns all assets paginated
func (pc *pageController) getProducts(w http.ResponseWriter, r *http.Request) {
	var params getProductsArgs
	if err := pc.parseForm(r, &params); err != nil {
		pc.badRequestError(w, errors.Wrap(err, errParsingArgs))
		return
	}

	if params.Limit == 0 || params.Limit > maxLimit {
		params.Limit = maxLimit
	}

	products, err := pc.products.FindAllProducts(r.Context(), params.Limit)
	if err != nil {
		pc.serverError(w, errors.Wrap(err, errFetchingProduct))
		return
	}

	pc.renderJSON(w, &productsResponse{
		products: products,
	})
}
