package server

import (
	"encoding/json"
	"time"

	"github.com/splice/catalog-interview/model"
)

type productPresenter struct {
	UUID          string     `json:"uuid"`
	Name          string     `json:"name"`
	Cost          int        `json:"cost"`
	MainGenreUUID string     `json:"main_genre_uuid,omitempty"`
	Description   string     `json:"description,omitempty"`
	CoverImageURL string     `json:"cover_image_url,omitempty"`
	DemoAudioURL  string     `json:"demo_audio_url,omitempty"`
	Permalink     string     `json:"permalink,omitempty"`
	SKU           string     `json:"sku,omitempty"`
	Tags          []string   `json:"tags,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty"`
}

type productsResponse struct {
	products []*model.Product
}

type getProductsPresenter struct {
	Products  []*productPresenter
	TotalCost int
}

// MarshalJSON presents the structure in JSON format.
func (p *productsResponse) MarshalJSON() ([]byte, error) {
	var getProductsResponse getProductsPresenter
	var products []*productPresenter

	for _, product := range p.products {
		presenter := &productPresenter{
			UUID:          product.UUID,
			Name:          product.Name,
			Cost:          product.Cost,
			Description:   product.Description,
			CoverImageURL: product.CoverImageURL,
			DemoAudioURL:  product.DemoAudioURL,
			Permalink:     product.Permalink,
			SKU:           product.SKU,
			CreatedAt:     product.CreatedAt,
			UpdatedAt:     product.UpdatedAt,
			Tags:          product.Tags,
		}
		if product.DeletedAt != nil {
			presenter.DeletedAt = product.DeletedAt
		}
		getProductsResponse.TotalCost += product.Cost
		products = append(products, presenter)
	}

	if len(products) > 0 {
		getProductsResponse.Products = products
	}

	return json.Marshal(getProductsResponse)
}
