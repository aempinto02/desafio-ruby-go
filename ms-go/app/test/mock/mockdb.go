package mock

import (
	"ms-go/app/models"
)

var MockData = []models.Product{
	{ID: 1, Name: "Product 1", Brand: "Brand A", Price: 100.0, Description: "Desc 1", Stock: 10},
	{ID: 2, Name: "Product 2", Brand: "Brand B", Price: 200.0, Description: "Desc 2", Stock: 20},
}

// GetMockProductByID returns a mock product by ID
func GetMockProductByID(id int) (*models.Product, error) {
	for _, product := range MockData {
		if product.ID == id {
			return &product, nil
		}
	}
	return nil, nil // Return nil or an error if product with ID not found in mock data
}
