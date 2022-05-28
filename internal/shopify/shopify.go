package shopify

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ProductDetails struct {
	ID                   int                     `json:"id"`
	Title                string                  `json:"title"`
	Handle               string                  `json:"handle"`      // unique identifier slug
	Description          string                  `json:"description"` // HTML description
	TimePublished        time.Time               `json:"published_at"`
	TimeCreated          time.Time               `json:"created_at"`
	Vendor               string                  `json:"vendor"`
	Type                 string                  `json:"type"`
	Tags                 []string                `json:"tags"`
	Price                int                     `json:"price"`
	PriceMin             int                     `json:"price_min"`
	PriceMax             int                     `json:"price_max"`
	Available            bool                    `json:"available"`
	PriceVaries          bool                    `json:"price_varies"`
	CompareAtPrice       *int                    `json:"compare_at_price"`
	CompareAtPriceMin    int                     `json:"compare_at_price_min"`
	CompareAtPriceMax    int                     `json:"compare_at_price_max"`
	CompareAtPriceVaries bool                    `json:"compare_at_price_varies"`
	Variants             []ProductDetailsVariant `json:"variants"`
	Images               []string                `json:"images"`         // image URLs, without http(s): prefix
	FeaturedImage        string                  `json:"featured_image"` // also without http(s): prefix
	Options              []ProductOption         `json:"options"`
	URL                  string                  `json:"url"` // path after domain
	// TODO: media (do we care?)
}

type ProductOption struct {
	Name     string   `json:"name"`     // name of option (example: "Size")
	Position int      `json:"position"` // ordering?
	Values   []string `json:"values"`   // all possible enum values (example: ["XS","S","M","L","XL","Custom sizing"])
}

type ProductVariant struct {
	ID               int           `json:"id"`                // [LD]
	Title            string        `json:"title"`             // [LD]
	Option1          string        `json:"option1"`           // [LD]
	Option2          string        `json:"option2"`           // [LD]
	Option3          string        `json:"option3"`           // [LD]
	SKU              string        `json:"sku"`               // [LD]
	RequiresShipping bool          `json:"requires_shipping"` // [LD]
	Taxable          bool          `json:"taxable"`           // [LD]
	FeaturedImage    *ProductImage `json:"featured_image"`    // can be null [LD]
	Available        bool          `json:"available"`         // [LD]
	Options          []string      `json:"options"`
	// CompareAtPrice // [LD] null on both, need another example to see what it is
}

type ProductDetailsVariant struct {
	ProductVariant
	Price               int    `json:"price"` // Price in cent units as int, seems more reasonable
	Weight              int    `json:"weight"`
	Name                string `json:"name"`                 // full name product+variant[D]
	PublicTitle         string `json:"public_title"`         // [D]
	InventoryQuantity   int    `json:"inventory_quantity"`   // can be negative, auto decrements on sale?
	InventoryManagement string `json:"inventory_management"` // example: "shopify"
	InventoryPolicy     string `json:"inventory_policy"`     // "deny", ?
	Barcode             string `json:"barcode"`              // "" ??
}

type ProductImage struct {
	ID          int       `json:"id"`         // image ID
	ProductID   int       `json:"product_id"` // product ID (equal to parent product)
	Position    int       `json:"position"`   // ordering?
	TimeCreated time.Time `json:"created_at"` // timestamp
	TimeUpdated time.Time `json:"updated_at"` // timestamp
	Alt         string    `json:"alt"`        // alt tag text
	Width       int       `json:"width"`      // width in px
	Height      int       `json:"height"`     // height in px
	Src         string    `json:"src"`        // URL
	VariantIDs  []int     `json:"variant_ids"`
}

// ProductDetails fetches detailed info about a product with handle
func FetchProductDetails(ctx context.Context, storeDomain, productHandle string) (*ProductDetails, error) {
	uri := fmt.Sprintf("https://%s/products/%v.js", storeDomain, productHandle)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}

	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code fetching %v from %v: %v", productHandle, storeDomain, res.StatusCode)
	}

	var d ProductDetails
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&d)
	return &d, err
}
