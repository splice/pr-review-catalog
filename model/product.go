package model

import (
	"time"
)

// A Product is a potentially sellable item on the Splice marketplace.
type Product struct {
	ID            int        `json:"id"`
	UUID          string     `json:"uuid"`
	Name          string     `json:"name"`
	Cost          int        `json:"cost"`
	Description   string     `json:"description"`
	CoverImageURL string     `json:"cover_image_url"`
	DemoAudioURL  string     `json:"demo_audio_url"`
	Permalink     string     `json:"permalink"`
	SKU           string     `json:"sku"`
	Tags          []string   `json:"tags"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at"`
}
