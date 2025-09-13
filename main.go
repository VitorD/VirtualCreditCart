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

func getCreditCards(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, creditCards)
}

func main() {
	router := gin.Default()
	router.GET("/creditCards", getCreditCards)

	router.Run("localhost:5000")
}
