package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

// PageControllerParams contains the params for the PageController
type PageControllerParams struct {
	Products productRepository
}

// NewPageController initializes the underlying page controller and
// returns a pointer to the HTTP router.
func NewPageController(params *PageControllerParams) *mux.Router {
	pc := &pageController{
		mux:      mux.NewRouter(),
		products: params.Products,
		dec:      schema.NewDecoder(),
	}
	pc.dec.IgnoreUnknownKeys(true)

	pc.mux.HandleFunc("/products", pc.getProducts).Methods(http.MethodGet)

	return pc.mux
}
