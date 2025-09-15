package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Modelo da tabela
type CreditCard struct {
	ID         uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	CardNumber string `json:"cardNumber"`
	UserId     string `json:"userId"`
}

// Inicializa conexão com o PostgreSQL
func initDB() *gorm.DB {
	dsn := "host=localhost user=postgres password=postgres dbname=creditCard port=5423 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Erro ao conectar no banco: ", err)
	}

	// Cria a tabela se não existir
	err = db.AutoMigrate(&CreditCard{})
	if err != nil {
		log.Fatal("Erro ao migrar schema: ", err)
	}

	return db
}

func main() {
	db := initDB()
	router := gin.Default()

	// Listar todos os cartões
	router.GET("/creditCards", func(c *gin.Context) {
		var cards []CreditCard
		if err := db.Find(&cards).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, cards)
	})

	// Buscar cartão por ID
	router.GET("/creditCards/:id", func(c *gin.Context) {
		id := c.Param("id")
		var card CreditCard
		if err := db.First(&card, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Cartão não encontrado"})
			return
		}
		c.JSON(http.StatusOK, card)
	})

	// Criar cartão
	router.POST("/creditCards", func(c *gin.Context) {
		var newCard CreditCard
		if err := c.BindJSON(&newCard); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
			return
		}
		if err := db.Create(&newCard).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, newCard)
	})

	// Atualizar cartão
	router.PUT("/creditCards/:id", func(c *gin.Context) {
		id := c.Param("id")
		var card CreditCard
		if err := db.First(&card, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Cartão não encontrado"})
			return
		}

		var updateData CreditCard
		if err := c.BindJSON(&updateData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
			return
		}

		// Atualiza apenas os campos enviados
		if err := db.Model(&card).Updates(updateData).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, card)
	})

	// Deletar cartão
	router.DELETE("/creditCards/:id", func(c *gin.Context) {
		id := c.Param("id")
		if err := db.Delete(&CreditCard{}, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusNoContent)
	})

	router.Run("localhost:5000")
}
