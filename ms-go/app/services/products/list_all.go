package products

import (
	"context"
	"ms-go/app/helpers"
	"ms-go/app/models"
	mockdb "ms-go/app/test/mock"
	"ms-go/db"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

func ListAll() ([]models.Product, error) {
	if useMock {
		return mockdb.MockData, nil
	}

	var products = []models.Product{}

	data, err := db.Connection().Find(context.TODO(), bson.D{})

	if err != nil {
		return nil, &helpers.GenericError{Msg: err.Error(), Code: http.StatusInternalServerError}
	}

	for data.Next(context.TODO()) {
		var product models.Product

		err := data.Decode(&product)

		if err != nil {
			return nil, &helpers.GenericError{Msg: err.Error(), Code: http.StatusInternalServerError}
		}

		products = append(products, product)
	}

	defer db.Disconnect()

	return products, nil
}
