package controllers

import (
	"context"
	"net/http"
	"productManagmentBackend/database"
	"productManagmentBackend/models"
	"productManagmentBackend/pkg/cache"
	"productManagmentBackend/pkg/logger"
	"productManagmentBackend/pkg/queue"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateProduct(c *gin.Context) {
	var product models.Product
	ctx := context.Background()

	if err := c.ShouldBindJSON(&product); err != nil {
		logger.Log.WithError(err).Error("Failed to bind JSON")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the user exists
	var user models.User
	if err := database.DB.First(&user, product.UserID).Error; err != nil {
		logger.Log.WithError(err).Error("User not found")
		c.JSON(http.StatusBadRequest, gin.H{"error": "User does not exist"})
		return
	}

	// Insert product into the database
	if result := database.DB.Create(&product); result.Error != nil {
		logger.Log.WithError(result.Error).Error("Failed to create product")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}

	// Queue image processing job
	if len(product.ProductImages) > 0 {
		job := queue.ImageProcessingJob{
			ProductID: product.ID,
			Images:    product.ProductImages,
		}
		if err := queue.PublishImageJob(ctx, job); err != nil {
			logger.Log.WithError(err).Error("Failed to queue image processing job")
		}
	}

	c.JSON(http.StatusCreated, product)
}

func GetProduct(c *gin.Context) {
	id := c.Param("id")
	productID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		logger.Log.WithError(err).Error("Invalid product ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	// Try to get from cache first
	if product, err := cache.GetProduct(uint(productID)); err == nil {
		c.JSON(http.StatusOK, product)
		return
	}

	// If not in cache, get from database
	var product models.Product
	if result := database.DB.First(&product, productID); result.Error != nil {
		logger.Log.WithError(result.Error).Error("Product not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// Cache the product
	if err := cache.SetProduct(&product); err != nil {
		logger.Log.WithError(err).Error("Failed to cache product")
	}

	c.JSON(http.StatusOK, product)
}

func GetProducts(c *gin.Context) {
	var products []models.Product
	query := database.DB

	// Filter by user_id if provided
	if userID := c.Query("user_id"); userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	// Filter by price range if provided
	if minPrice := c.Query("min_price"); minPrice != "" {
		query = query.Where("product_price >= ?", minPrice)
	}
	if maxPrice := c.Query("max_price"); maxPrice != "" {
		query = query.Where("product_price <= ?", maxPrice)
	}

	// Filter by product name if provided
	if name := c.Query("name"); name != "" {
		query = query.Where("product_name ILIKE ?", "%"+name+"%")
	}

	if result := query.Find(&products); result.Error != nil {
		logger.Log.WithError(result.Error).Error("Failed to fetch products")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}

	c.JSON(http.StatusOK, products)
}