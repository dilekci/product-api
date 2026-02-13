package request

import "product-app/services/product/internal/usecase/model"

// AddProductRequest represents the request payload used to create a new product.
// It is typically populated from an incoming HTTP JSON request.
type AddProductRequest struct {
	// Product name
	Name string `json:"name"`

	// Product price
	Price float32 `json:"price"`

	// Product description
	Description string `json:"description"`

	// Discount amount or rate applied to the product
	Discount float32 `json:"discount"`

	// Store or seller name
	Store string `json:"store"`

	// List of product image URLs
	ImageUrls []string `json:"image_urls"`

	// Category identifier of the product
	CategoryID int64 `json:"category_id"`
}

// ToModel converts AddProductRequest to ProductCreate domain model.
func (addProductRequest AddProductRequest) ToModel() model.ProductCreate {
	return model.ProductCreate{
		Name:        addProductRequest.Name,
		Price:       addProductRequest.Price,
		Description: addProductRequest.Description,
		Discount:    addProductRequest.Discount,
		Store:       addProductRequest.Store,
		ImageUrls:   addProductRequest.ImageUrls,
		CategoryID:  addProductRequest.CategoryID,
	}
}
