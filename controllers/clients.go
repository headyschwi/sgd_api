package controllers

import (
	"net/http"
	"sgd_api/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ClientController struct {
	db *gorm.DB
}

func NewClientController(db *gorm.DB) *ClientController {
	return &ClientController{db: db}
}

func (cc ClientController) CreateClient(c *gin.Context) {
	var input struct {
		Name  string `json:"name" binding:"required"`
		Email string `json:"email" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := cc.db.Begin()

	client := models.Client{Name: input.Name, Email: input.Email}
	if err := tx.Create(&client).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cart := models.Cart{ClientID: client.ID}
	if err := tx.Create(&cart).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx.Model(&client).Association("Cart").Append(&cart)
	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"data": client})

}

func (cc ClientController) GetClients(c *gin.Context) {
	var clients []models.Client

	if err := cc.db.Find(&clients).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"data": clients})
}

func (cc ClientController) GetClient(c *gin.Context) {
	var client models.Client

	if err := cc.db.Preload("Orders").Preload("Cart").First(&client, c.Param("client_id")).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": client})
}

func (cc ClientController) UpdateClient(c *gin.Context) {
	var client models.Client

	if err := cc.db.First(&client, c.Param("client_id")).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.ShouldBindJSON(&client); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := cc.db.Save(&client).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": client})
}

func (cc ClientController) DeleteClient(c *gin.Context) {
	var client models.Client

	if err := cc.db.First(&client, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := cc.db.Delete(&client).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": "Client deleted successfully!"})
}
