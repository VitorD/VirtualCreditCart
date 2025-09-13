package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type creditCard struct {
	ID         string `json:"id"`
	CardNumber string `json:"cardNumber"`
	UserId     string `json:"userId"`
}

var creditCards = []creditCard{
	{ID: "1", CardNumber: "1234656727", UserId: "123"},
}

func getAllVirtualCreditCards(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, creditCards)
}

func getVirtualCreditCardById(c *gin.Context) {
	id := c.Param("id")

	for _, a := range creditCards {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"messsage": "credit card not found"})
}

func creatVirtualCreditCard(c *gin.Context) {
	var newCreditCard creditCard

	if err := c.BindJSON(&newCreditCard); err != nil {
		return
	}
	creditCards = append(creditCards, newCreditCard)

	c.IndentedJSON(http.StatusNoContent, nil)
}

func deleteVirtualCreditCard(c *gin.Context) {
	var newCreditCard creditCard

	if err := c.BindJSON(&newCreditCard); err != nil {
		return
	}
	creditCards = append(creditCards, newCreditCard)

	c.IndentedJSON(http.StatusNoContent, nil)
}

func main() {
	router := gin.Default()
	router.GET("/creditCards", getAllVirtualCreditCards)
	router.GET("/creditCards/:id", getVirtualCreditCardById)
	router.POST("/creditCards", creatVirtualCreditCard)
	router.DELETE("/creditCards", deleteVirtualCreditCard)

	router.Run("localhost:5000")
}
