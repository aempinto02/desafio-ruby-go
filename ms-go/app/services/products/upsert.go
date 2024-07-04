package products

import (
	"context"
	"ms-go/app/helpers"
	"ms-go/app/models"
	mockdb "ms-go/app/test/mock"
	"ms-go/db"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func UpsertProduct(data models.Product) (*models.Product, error) {
	if useMock {
		// Handle mock data scenario
		if data.ID > 0 {
			// Find the product in mock data
			for _, p := range mockdb.MockData {
				if p.ID == data.ID {
					// Update the product in mock data
					if data.Name != "" {
						p.Name = data.Name
					}
					if data.Brand != "" {
						p.Brand = data.Brand
					}
					if data.Price != 0 {
						p.Price = data.Price
					}
					if data.Description != "" {
						p.Description = data.Description
					}
					if data.Stock != 0 {
						p.Stock = data.Stock
					}
					return &p, nil
				}
			}
			// If product not found in mock data, return an error
			return nil, &helpers.GenericError{Msg: "Product with ID not found in mock data", Code: http.StatusNotFound}
		}

		// Insert new product in mock data
		// (Assuming mock data supports insertion)
		newProduct := models.Product{
			ID:          len(mockdb.MockData) + 1, // Simulate auto-increment
			Name:        data.Name,
			Brand:       data.Brand,
			Price:       data.Price,
			Description: data.Description,
			Stock:       data.Stock,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		mockdb.MockData = append(mockdb.MockData, newProduct)
		return &newProduct, nil
	}

	if data.ID > 0 {
		// Check if the product with the given ID exists
		existingProduct := models.Product{}
		err := db.Connection().FindOne(context.TODO(), bson.M{"id": data.ID}).Decode(&existingProduct)
		if err == nil {
			// Product with given ID exists, update it
			update := bson.M{}
			if data.Name != "" {
				update["name"] = data.Name
			}
			if data.Brand != "" {
				update["brand"] = data.Brand
			}
			if data.Price != 0 {
				update["price"] = data.Price
			}
			if data.Description != "" {
				update["description"] = data.Description
			}
			if data.Stock != 0 {
				update["stock"] = data.Stock
			}

			filter := bson.M{"id": data.ID}
			opts := options.Update().SetUpsert(false)
			_, err = db.Connection().UpdateOne(context.TODO(), filter, bson.M{"$set": update}, opts)
			if err != nil {
				return nil, &helpers.GenericError{Msg: err.Error(), Code: http.StatusInternalServerError}
			}

			// Fetch the updated document
			err = db.Connection().FindOne(context.TODO(), bson.M{"id": data.ID}).Decode(&existingProduct)
			if err != nil {
				return nil, &helpers.GenericError{Msg: err.Error(), Code: http.StatusInternalServerError}
			}

			return &existingProduct, nil
		} else if err.Error() == "mongo: no documents in result" {
			// Product with given ID does not exist, proceed to insert
		} else {
			return nil, &helpers.GenericError{Msg: err.Error(), Code: http.StatusInternalServerError}
		}
	}

	// Find the document with the highest id
	var maxProduct models.Product
	opts := options.FindOne().SetSort(bson.D{{Key: "id", Value: -1}})
	err := db.Connection().FindOne(context.TODO(), bson.D{}, opts).Decode(&maxProduct)
	if err != nil && err.Error() != "mongo: no documents in result" {
		return nil, &helpers.GenericError{Msg: err.Error(), Code: http.StatusInternalServerError}
	}

	// Set the new id
	if maxProduct.ID == 0 {
		data.ID = 1 // if no documents exist, start from 1
	} else {
		data.ID = maxProduct.ID + 1
	}

	if err := data.Validate(); err != nil {
		return nil, &helpers.GenericError{Msg: err.Error(), Code: http.StatusUnprocessableEntity}
	}

	// Set timestamps
	data.CreatedAt = time.Now()
	data.UpdatedAt = data.CreatedAt

	// Insert the new document
	_, err = db.Connection().InsertOne(context.TODO(), data)
	if err != nil {
		return nil, &helpers.GenericError{Msg: err.Error(), Code: http.StatusInternalServerError}
	}

	defer db.Disconnect()

	return &data, nil
}
