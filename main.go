package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type creditCard struct {
	ID         string `json:"id"`
	CardNumber string `json:"cardNumber"`
	UserId     string `json:"userId"`
}

func getAllVirtualCreditCards(collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {

		filter := bson.D{}

		cursor, _ := collection.Find(context.TODO(), filter)
		var creditCards []creditCard

		if err := cursor.All(context.TODO(), &creditCards); err != nil {
			panic(err)
		}

		c.IndentedJSON(http.StatusOK, creditCards)
	}
}

func getVirtualCreditCardById(collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		filter := bson.D{{Key: "id", Value: id}}

		result := collection.FindOne(context.TODO(), filter)
		var creditCard creditCard

		if err := result.Decode(&creditCard); err != nil {
			panic(err)
		}

		c.IndentedJSON(http.StatusOK, creditCard)
	}
}

func updateVirtualCreditCard(collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var creditCard creditCard
		if err := c.BindJSON(&creditCard); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		}

		filter := bson.D{{Key: "id", Value: id}}

		update := bson.D{{Key: "$set", Value: bson.D{{Key: "cardnumber", Value: creditCard.CardNumber}}}}

		result := collection.FindOneAndUpdate(context.TODO(), filter, update)

		c.IndentedJSON(http.StatusOK, result)
	}
}

func createVirtualCreditCard(collection *mongo.Collection) gin.HandlerFunc {

	return func(c *gin.Context) {
		var newCreditCard creditCard

		if err := c.BindJSON(&newCreditCard); err != nil {
			return
		}

		_, err := collection.InsertOne(context.TODO(), newCreditCard)
		if err != nil {
			log.Fatal("Could not save the virtual credit card id: " + newCreditCard.ID)
			return
		}

		c.IndentedJSON(http.StatusNoContent, nil)
	}

}

func deleteVirtualCreditCard(collection *mongo.Collection) gin.HandlerFunc {

	return func(c *gin.Context) {
		id := c.Param("id")

		filter := bson.D{{Key: "id", Value: id}}

		_, err := collection.DeleteOne(context.TODO(), filter)
		if err != nil {
			panic(err)
		}

		c.IndentedJSON(http.StatusNoContent, nil)
	}
}

func main() {
	uri := os.Getenv("MONGODB_URI")
	docs := "www.mongodb.com/docs/drivers/go/current/"

	if uri == "" {
		log.Fatal("Set your 'MONGODB_URI' environment variable" +
			"See: " + docs +
			"usage-examples/#environment-variable")
	}

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	coll := client.Database("db").Collection("creditCards")

	router := gin.Default()
	router.GET("/creditCards", getAllVirtualCreditCards(coll))
	router.GET("/creditCards/:id", getVirtualCreditCardById(coll))
	router.POST("/creditCards", createVirtualCreditCard(coll))
	router.DELETE("/creditCards/:id", deleteVirtualCreditCard(coll))
	router.PUT("/creditCards/:id", updateVirtualCreditCard(coll))

	router.Run("localhost:5000")
}
