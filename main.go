package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"kozan/db"
	"kozan/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func main() {
	// .env dosyasını yükle
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	// Veritabanı bağlantısını başlat
	db.InitDB()

	// Uploads klasörünü oluştur
	if err := os.MkdirAll("uploads", 0755); err != nil {
		log.Fatal(err)
	}

	// Gin router'ı oluştur
	r := gin.Default()

	// Statik dosya sunucusu
	r.Static("/uploads", "./uploads")

	// CORS ayarları
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// API rotaları
	api := r.Group("/api")
	{
		// Auth rotaları
		auth := api.Group("/auth")
		{
			auth.POST("/login", loginHandler)
		}

		// Admin rotaları
		admin := api.Group("/admin")
		admin.Use(authMiddleware())
		{
			// Dosya yükleme
			admin.POST("/upload", uploadHandler)

			// Ürünler
			admin.GET("/products", getProductsHandler)
			admin.POST("/products", createProductHandler)
			admin.PUT("/products/:id", updateProductHandler)
			admin.DELETE("/products/:id", deleteProductHandler)

			// Hizmetler
			admin.GET("/services", getServicesHandler)
			admin.POST("/services", createServiceHandler)
			admin.PUT("/services/:id", updateServiceHandler)
			admin.DELETE("/services/:id", deleteServiceHandler)

			// Hakkımızda
			admin.GET("/about", getAboutHandler)
			admin.PUT("/about", updateAboutHandler)

			// İletişim
			admin.GET("/contact", getContactHandler)
			admin.PUT("/contact", updateContactHandler)
		}

		// Public rotaları
		api.GET("/products", getPublicProductsHandler)
		api.GET("/services", getPublicServicesHandler)
		api.GET("/about", getPublicAboutHandler)
		api.GET("/contact", getPublicContactHandler)
	}

	// Sunucuyu başlat
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}

// Auth işlemleri
func loginHandler(c *gin.Context) {
	var loginData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var user models.User
	err := db.DB.QueryRow(
		"SELECT id, username FROM users WHERE username = $1 AND password = $2",
		loginData.Username,
		loginData.Password,
	).Scan(&user.ID, &user.Username)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// JWT token oluştur
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
		"user":  user,
	})
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// Ürün işlemleri
func getProductsHandler(c *gin.Context) {
	rows, err := db.DB.Query("SELECT id, name, description, price, image FROM products")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Image); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		products = append(products, p)
	}

	c.JSON(http.StatusOK, products)
}

func createProductHandler(c *gin.Context) {
	var product models.Product
	if err := c.BindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := db.DB.QueryRow(
		"INSERT INTO products (name, description, price, image) VALUES ($1, $2, $3, $4) RETURNING id",
		product.Name, product.Description, product.Price, product.Image,
	).Scan(&product.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, product)
}

func updateProductHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var product models.Product
	if err := c.BindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = db.DB.Exec(
		"UPDATE products SET name = $1, description = $2, price = $3, image = $4 WHERE id = $5",
		product.Name, product.Description, product.Price, product.Image, id,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	product.ID = id
	c.JSON(http.StatusOK, product)
}

func deleteProductHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	result, err := db.DB.Exec("DELETE FROM products WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted"})
}

// Hizmet işlemleri
func getServicesHandler(c *gin.Context) {
	rows, err := db.DB.Query("SELECT id, title, description, image FROM services")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var services []models.Service
	for rows.Next() {
		var s models.Service
		if err := rows.Scan(&s.ID, &s.Title, &s.Description, &s.Image); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		services = append(services, s)
	}

	c.JSON(http.StatusOK, services)
}

func createServiceHandler(c *gin.Context) {
	var service models.Service
	if err := c.BindJSON(&service); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := db.DB.QueryRow(
		"INSERT INTO services (title, description, image) VALUES ($1, $2, $3) RETURNING id",
		service.Title, service.Description, service.Image,
	).Scan(&service.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, service)
}

func updateServiceHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var service models.Service
	if err := c.BindJSON(&service); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = db.DB.Exec(
		"UPDATE services SET title = $1, description = $2, image = $3 WHERE id = $4",
		service.Title, service.Description, service.Image, id,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	service.ID = id
	c.JSON(http.StatusOK, service)
}

func deleteServiceHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	result, err := db.DB.Exec("DELETE FROM services WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Service deleted"})
}

// Hakkımızda işlemleri
func getAboutHandler(c *gin.Context) {
	var about models.About
	err := db.DB.QueryRow("SELECT id, title, content, image FROM about ORDER BY id LIMIT 1").
		Scan(&about.ID, &about.Title, &about.Content, &about.Image)

	if err != nil {
		c.JSON(http.StatusOK, models.About{}) // Veri yoksa boş döndür
		return
	}

	c.JSON(http.StatusOK, about)
}

func updateAboutHandler(c *gin.Context) {
	var about models.About
	if err := c.BindJSON(&about); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingID int
	err := db.DB.QueryRow("SELECT id FROM about ORDER BY id LIMIT 1").Scan(&existingID)
	if err != nil {
		// Kayıt yoksa yeni ekle
		err = db.DB.QueryRow(
			"INSERT INTO about (title, content, image) VALUES ($1, $2, $3) RETURNING id",
			about.Title, about.Content, about.Image,
		).Scan(&about.ID)
	} else {
		// Varolan kaydı güncelle
		_, err = db.DB.Exec(
			"UPDATE about SET title = $1, content = $2, image = $3 WHERE id = $4",
			about.Title, about.Content, about.Image, existingID,
		)
		about.ID = existingID
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, about)
}

// İletişim işlemleri
func getContactHandler(c *gin.Context) {
	var contact models.Contact
	err := db.DB.QueryRow("SELECT id, title, phone, email, address FROM contact ORDER BY id LIMIT 1").
		Scan(&contact.ID, &contact.Title, &contact.Phone, &contact.Email, &contact.Address)

	if err != nil {
		c.JSON(http.StatusOK, models.Contact{}) // Veri yoksa boş döndür
		return
	}

	c.JSON(http.StatusOK, contact)
}

func updateContactHandler(c *gin.Context) {
	var contact models.Contact
	if err := c.BindJSON(&contact); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingID int
	err := db.DB.QueryRow("SELECT id FROM contact ORDER BY id LIMIT 1").Scan(&existingID)
	if err != nil {
		// Kayıt yoksa yeni ekle
		err = db.DB.QueryRow(
			"INSERT INTO contact (title, phone, email, address) VALUES ($1, $2, $3, $4) RETURNING id",
			contact.Title, contact.Phone, contact.Email, contact.Address,
		).Scan(&contact.ID)
	} else {
		// Varolan kaydı güncelle
		_, err = db.DB.Exec(
			"UPDATE contact SET title = $1, phone = $2, email = $3, address = $4 WHERE id = $5",
			contact.Title, contact.Phone, contact.Email, contact.Address, existingID,
		)
		contact.ID = existingID
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, contact)
}

// Public handlers
func getPublicProductsHandler(c *gin.Context) {
	getProductsHandler(c)
}

func getPublicServicesHandler(c *gin.Context) {
	getServicesHandler(c)
}

func getPublicAboutHandler(c *gin.Context) {
	getAboutHandler(c)
}

func getPublicContactHandler(c *gin.Context) {
	getContactHandler(c)
}

// Dosya yükleme işleyicisi
func uploadHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dosya yüklenemedi"})
		return
	}

	// Benzersiz dosya adı oluştur
	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
	filepath := fmt.Sprintf("uploads/%s", filename)

	// Dosyayı kaydet
	if err := c.SaveUploadedFile(file, filepath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Dosya kaydedilemedi"})
		return
	}

	// Dosya URL'ini döndür
	fileURL := fmt.Sprintf("/uploads/%s", filename)
	c.JSON(http.StatusOK, gin.H{"url": fileURL})
}
