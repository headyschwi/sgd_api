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

	var input struct {
		ClientID uint `json:"client_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	var cart models.Cart
	if err := cc.db.Preload("CartItems").Where("client_id = ?", input.ClientID).First(&cart).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cart not found!"})
	}

	c.JSON(http.StatusOK, gin.H{"data": cart})

}

func (cc CartController) AddToCart(c *gin.Context) {

	var input struct {
		ProductID uint  `json:"product_id"`
		Amount    int64 `json:"amount"`
		ClientID  uint  `json:"client_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	var cart models.Cart
	if err := cc.db.Preload("CartItems").Where("client_id = ?", input.ClientID).First(&cart).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cart not found!"})
		return
	}

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
		Amount int64 `json:"amount"`
	}

	var cart models.Cart
	if err := cc.db.Where("client_id = ?", c.Param("client_id")).First(&cart).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cart not found!"})
	}

	var cartItem models.CartItem
	if err := cc.db.Where("cart_id = ? AND product_id = ?", cart.ID, c.Param("product_id")).First(&cartItem).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cart item not found!"})
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	if err := cc.db.Model(&cartItem).Updates(models.CartItem{Amount: input.Amount}).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"data": cartItem})
}

func (cc CartController) RemoveFromCart(c *gin.Context) {

	var cart models.Cart
	if err := cc.db.Where("client_id = ?", c.Param("client_id")).First(&cart).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cart not found!"})
	}

	var cartItem models.CartItem
	if err := cc.db.Where("cart_id = ? AND product_id = ?", cart.ID, c.Param("product_id")).First(&cartItem).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cart item not found!"})
	}

	if err := cc.db.Delete(&cartItem).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"data": true})
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

		// Busque o produto do banco de dados
		var product models.Product
		if err := tx.Where("id = ?", item.ProductID).First(&product).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Product not found!"})
			return
		}

		// Verifique se há estoque suficiente
		if product.Stock < item.Amount {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Not enough stock for product!"})
			return
		}

		product.Stock -= item.Amount

		// Reduza o estoque do produto
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
