package server

import (
	"encoding/json"
	"time"

	"github.com/splice/catalog-interview/model"
)

type productPresenter struct {
	UUID          string     `json:"uuid"`
	Name          string     `json:"name"`
	MainGenreUUID string     `json:"main_genre_uuid,omitempty"`
	Description   string     `json:"description,omitempty"`
	CoverImageURL string     `json:"cover_image_url,omitempty"`
	DemoAudioURL  string     `json:"demo_audio_url,omitempty"`
	Permalink     string     `json:"permalink,omitempty"`
	SKU           string     `json:"sku,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty"`
}

type productsResponse struct {
	products []*model.Product
}

// MarshalJSON presents the structure in JSON format.
func (p *productsResponse) MarshalJSON() ([]byte, error) {
	var products []*productPresenter

	for _, product := range p.products {
		presenter := &productPresenter{
			UUID:          product.UUID,
			Name:          product.Name,
			Description:   product.Description,
			CoverImageURL: product.CoverImageURL,
			DemoAudioURL:  product.DemoAudioURL,
			Permalink:     product.Permalink,
			SKU:           product.SKU,
			CreatedAt:     product.CreatedAt,
			UpdatedAt:     product.UpdatedAt,
		}
		if product.DeletedAt != nil {
			presenter.DeletedAt = product.DeletedAt
		}
		products = append(products, presenter)
	}

	return json.Marshal(products)
}
