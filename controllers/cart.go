package controllers

import (
	"net/http"
	"sgd_api/models"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type CartController struct {
	db *gorm.DB
}

func NewCartController(db *gorm.DB) *CartController {
	return &CartController{db}
}

func (cc CartController) GetCart(c *gin.Context) {

	var cart models.Cart
	if err := cc.db.Preload("CartItems").Where("client_id = ?", c.Param("client_id")).First(&cart).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cart not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": cart})
}

func (cc CartController) AddToCart(c *gin.Context) {

	var input struct {
		ClientID  uint  `json:"client_id"`
		ProductID uint  `json:"product_id"`
		Amount    int64 `json:"amount"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	var cart models.Cart
	if err := cc.db.Preload("CartItems").Where("client_id = ?", input.ClientID).First(&cart).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cart not found!"})
		return
	}

	var existingItem models.CartItem
	if err := cc.db.Where("cart_id = ? AND product_id = ?", cart.ID, input.ProductID).First(&existingItem).Error; err != nil {
		var product models.Product
		if err := cc.db.Where("id = ?", input.ProductID).First(&product).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Product not found!"})
			return
		}

		cartItem := models.CartItem{
			CartID:    cart.ID,
			ProductID: product.ID,
			Amount:    input.Amount,
			Price:     product.Price,
		}

		cart.CartItems = append(cart.CartItems, cartItem)
		cart.TotalPrice = cart.TotalPrice.Add(decimal.NewFromFloat(float64(input.Amount)).Mul(product.Price))

	} else {
		var product models.Product
		if err := cc.db.Where("id = ?", input.ProductID).First(&product).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Product not found!"})
			return
		}

		existingItem.Amount += input.Amount
		existingItem.Price = product.Price
		cart.TotalPrice = cart.TotalPrice.Add(decimal.NewFromFloat(float64(input.Amount)).Mul(product.Price))
	}

	tx := cc.db.Begin()
	if err := tx.Save(&cart).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"data": cart})
}

func (cc CartController) UpdateCartItem(c *gin.Context) {

	var input struct {
		ClientID  uint  `json:"client_id"`
		ProductID uint  `json:"product_id"`
		Amount    int64 `json:"amount"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	var cart models.Cart
	if err := cc.db.Preload("CartItems").Where("client_id = ?", input.ClientID).First(&cart).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cart not found!"})
		return
	}

	var existingItem models.CartItem
	if err := cc.db.Where("cart_id = ? AND product_id = ?", cart.ID, input.ProductID).First(&existingItem).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product not found in cart!"})
		return
	}

	var product models.Product
	if err := cc.db.Where("id = ?", input.ProductID).First(&product).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product not found!"})
		return
	}

	cart.TotalPrice = cart.TotalPrice.Sub(decimal.NewFromFloat(float64(existingItem.Amount)).Mul(product.Price))
	cart.TotalPrice = cart.TotalPrice.Add(decimal.NewFromFloat(float64(input.Amount)).Mul(product.Price))

	existingItem.Amount = input.Amount
	existingItem.Price = product.Price

	tx := cc.db.Begin()
	if err := tx.Save(&cart).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := tx.Save(&existingItem).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"data": cart})
}

func (cc CartController) RemoveFromCart(c *gin.Context) {

	var input struct {
		ClientID  uint `json:"client_id"`
		ProductID uint `json:"product_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	var cart models.Cart
	if err := cc.db.Preload("CartItems").Where("client_id = ?", input.ClientID).First(&cart).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cart not found!"})
		return
	}

	var existingItem models.CartItem
	if err := cc.db.Where("cart_id = ? AND product_id = ?", cart.ID, input.ProductID).First(&existingItem).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product not found in cart!"})
		return
	}

	var product models.Product
	if err := cc.db.Where("id = ?", input.ProductID).First(&product).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product not found!"})
		return
	}

	cart.TotalPrice = cart.TotalPrice.Sub(decimal.NewFromFloat(float64(existingItem.Amount)).Mul(product.Price))

	tx := cc.db.Begin()
	if err := tx.Save(&cart).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := tx.Delete(&existingItem).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"data": cart})
}

func (cc CartController) Checkout(c *gin.Context) {
	var input struct {
		ClientID uint `json:"client_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var cart models.Cart
	if err := cc.db.Preload("CartItems").Where("client_id = ?", input.ClientID).First(&cart).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cart not found!"})
		return
	}

	tx := cc.db.Begin()

	order := models.Order{
		ClientID:    input.ClientID,
		OrderItems:  []models.OrderItem{},
		TotalPrice:  decimal.NewFromFloat(0),
		OrderStatus: "pending",
	}

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, item := range cart.CartItems {
		orderItem := models.OrderItem{
			OrderID:   order.ID,
			ProductID: item.ProductID,
			Amount:    item.Amount,
			Price:     item.Price,
		}

		if err := tx.Create(&orderItem).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var product models.Product
		if err := tx.Where("id = ?", item.ProductID).First(&product).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Product not found!"})
			return
		}

		if product.Stock < item.Amount {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Not enough stock for product!"})
			return
		}

		product.Stock -= item.Amount

		if err := tx.Save(&product).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Could not update product stock!"})
			return
		}

		if err := tx.Delete(&item).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		order.TotalPrice = order.TotalPrice.Add(decimal.NewFromFloat(float64(item.Amount)).Mul(item.Price))
	}

	if err := tx.Save(&order).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cart.TotalPrice = decimal.NewFromFloat(0)
	if err := tx.Save(&cart).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := tx.Where("cart_id = ?", cart.ID).Delete(&models.CartItem{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not clear cart items!"})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"data": order})
}
