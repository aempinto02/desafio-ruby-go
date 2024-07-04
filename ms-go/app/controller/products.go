package controller

import (
	"encoding/json"
	"fmt"
	"ms-go/app/helpers"
	"ms-go/app/models"
	"ms-go/app/services/kafka/producers"
	"ms-go/app/services/products"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func IndexProducts(c *gin.Context) {
	fmt.Println("IndexProducts method entered")
	all, err := products.ListAll()

	if err != nil {
		fmt.Printf("Error in ListAll: %v\n", err)
		switch e := err.(type) {
		case *helpers.GenericError:
			c.JSON(e.Code, gin.H{"message": e.Error()})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
	}

	jsonData, jsonErr := json.Marshal(all)
	if jsonErr != nil {
		fmt.Printf("Error in JSON Marshalling: %v\n", jsonErr)
		c.JSON(http.StatusInternalServerError, gin.H{"message": jsonErr.Error()})
		return
	}
	fmt.Printf("JSON Response: %s", string(jsonData))

	c.JSON(http.StatusOK, gin.H{"data": all})
}

func ShowProducts(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid product ID"})
		return
	}

	product, err := products.Details(models.Product{ID: id})

	if err != nil {
		switch e := err.(type) {
		case *helpers.GenericError:
			c.JSON(e.Code, gin.H{"message": e.Error()})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": product})
}

func CreateProducts(c *gin.Context) {
	var params models.Product

	if err := c.BindJSON(&params); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	product, err := products.UpsertProduct(params)

	// Kafka send message to Rails topic
	if producers.EnableKafkaProducer {
		producers.SendToRailsKafkaTopic(params)
	}

	if err != nil {
		switch e := err.(type) {
		case *helpers.GenericError:
			c.JSON(e.Code, gin.H{"message": e.Error()})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
	}
	c.JSON(http.StatusCreated, gin.H{"data": product})
}

func UpdateProducts(c *gin.Context) {
	var params models.Product

	if err := c.BindJSON(&params); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))

	params.ID = id

	product, err := products.UpsertProduct(params)

	// Kafka send message to Rails topic
	if producers.EnableKafkaProducer {
		producers.SendToRailsKafkaTopic(params)
	}

	if err != nil {
		switch e := err.(type) {
		case *helpers.GenericError:
			c.JSON(e.Code, gin.H{"message": e.Error()})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": product})
}
