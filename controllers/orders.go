package controllers

import (
	"net/http"
	"sgd_api/models"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type OrdersController struct {
	db *gorm.DB
}

func NewOrderController(db *gorm.DB) *OrdersController {
	return &OrdersController{db: db}
}

func (oc OrdersController) GetOrders(c *gin.Context) {

	var orders []models.Order

	if err := oc.db.Preload("OrderItems").Find(&orders).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Orders not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": orders})
}

func (oc OrdersController) UpdateOrderStatus(c *gin.Context) {

	var order models.Order

	if err := oc.db.First(&order, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order not found!"})
		return
	}

	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := oc.db.Save(&order).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order not updated!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": order})
}

func (oc OrdersController) OrderAgain(c *gin.Context) {

	var oldOrder models.Order
	if err := oc.db.First(&oldOrder, c.Param("order_id")).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order not found!"})
		return
	}

	tx := oc.db.Begin()

	orderItems := []models.OrderItem{}
	for _, item := range oldOrder.OrderItems {
		var product models.Product
		if err := oc.db.First(&product, item.ProductID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Product not found!"})
			return
		}

		// Check if there is enough stock
		if product.Stock < item.Amount {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Not enough stock for product!"})
			return
		}

		orderItem := models.OrderItem{
			ProductID: product.ID,
			Amount:    item.Amount,
			Price:     product.Price,
		}

		orderItems = append(orderItems, orderItem)

		// Reduce the stock
		product.Stock -= item.Amount
		if err := tx.Save(&product).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Could not update product stock!"})
			return
		}
	}

	orderPrice := decimal.NewFromFloat(0.0)
	for _, item := range orderItems {
		if err := tx.Create(&item).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Order not created!"})
			return
		}
		orderPrice = orderPrice.Add(item.Price.Mul(decimal.NewFromFloat(float64(item.Amount))))
	}

	newOrder := models.Order{
		ClientID:    oldOrder.ClientID,
		OrderItems:  orderItems,
		TotalPrice:  orderPrice,
		OrderStatus: "pending",
	}

	if err := tx.Create(&newOrder).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order not created!"})
		return
	}

	tx.Model(&newOrder).Association("OrderItems").Append(&orderItems)

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"data": newOrder})
}
