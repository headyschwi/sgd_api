package controllers

import (
	"net/http"
	"sgd_api/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProductsController struct {
	db *gorm.DB
}

func NewProductController(db *gorm.DB) *ProductsController {
	return &ProductsController{db: db}
}

func (pc *ProductsController) CreateProduct(c *gin.Context) {
	var product models.Product

	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	if err := pc.db.Create(&product).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"data": product})
}

func (pc *ProductsController) GetProducts(c *gin.Context) {
	var products []models.Product

	if err := pc.db.Find(&products).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"data": products})
}

func (pc *ProductsController) GetProduct(c *gin.Context) {
	var product models.Product

	if err := pc.db.First(&product, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"data": product})
}

func (pc *ProductsController) UpdateProduct(c *gin.Context) {

	var product models.Product
	if err := pc.db.First(&product, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := pc.db.Save(&product).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": product})
}

func (pc *ProductsController) DeleteProduct(c *gin.Context) {
	var product models.Product

	if err := pc.db.First(&product, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	if err := pc.db.Delete(&product).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"data": product})
}

func (pc *ProductsController) GetOrderWithProduct(c *gin.Context) {

	var product models.Product
	var orders []models.Order

	if err := pc.db.First(&product, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	if err := pc.db.Model(&product).Association("Orders").Find(&orders); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"data": orders})
}
