package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"kozan/db"
	"kozan/models"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize database
	db.InitDB()

	// Initialize Gin
	r := gin.Default()

	// Static file server
	r.Static("/uploads", "./uploads")
	if err := os.MkdirAll("uploads", 0755); err != nil {
		log.Fatal(err)
	}

	// CORS configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://admin.localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Options handler for CORS preflight requests
	r.OPTIONS("/*path", func(c *gin.Context) {
		c.Status(204)
	})

	// Public routes
	r.POST("/api/auth/login", loginHandler)

	// Admin routes (protected)
	admin := r.Group("/api/admin")
	admin.Use(authMiddleware())
	{
		// Products
		admin.GET("/products", getProductsHandler)
		admin.POST("/products", createProductHandler)
		admin.PUT("/products/:id", updateProductHandler)
		admin.DELETE("/products/:id", deleteProductHandler)

		// Services
		admin.GET("/services", getServicesHandler)
		admin.POST("/services", createServiceHandler)
		admin.PUT("/services/:id", updateServiceHandler)
		admin.DELETE("/services/:id", deleteServiceHandler)

		// About
		admin.GET("/about", getAboutHandler)
		admin.PUT("/about", updateAboutHandler)

		// Contact
		admin.GET("/contact", getContactHandler)
		admin.PUT("/contact", updateContactHandler)

		// Hero
		admin.GET("/hero", getHeroHandler)
		admin.PUT("/hero", updateHeroHandler)

		// Footer
		admin.GET("/footer", getFooterHandler)
		admin.PUT("/footer", updateFooterHandler)

		// Users routes
		admin.GET("/users", getUsersHandler)
		admin.POST("/users", createUserHandler)
		admin.PUT("/users/:id", updateUserHandler)
		admin.DELETE("/users/:id", deleteUserHandler)
	}

	// Public API routes
	api := r.Group("/api")
	{
		api.GET("/products", getProductsHandler)
		api.GET("/services", getServicesHandler)
		api.GET("/about", getAboutHandler)
		api.GET("/contact", getContactHandler)
	}

	// Start server
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz istek"})
		return
	}

	var user models.User
	var hashedPassword string
	err := db.DB.QueryRow("SELECT id, username, password, role_id FROM users WHERE username = $1", loginData.Username).
		Scan(&user.ID, &user.Username, &hashedPassword, &user.RoleID)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Kullanıcı adı veya şifre hatalı"})
		return
	}

	// Şifre kontrolü
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(loginData.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Kullanıcı adı veya şifre hatalı"})
		return
	}

	// JWT token oluştur
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"user_id":  user.ID,
		"role_id":  user.RoleID,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token oluşturulamadı"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"role_id":  user.RoleID,
		},
	})
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
			c.Abort()
			return
		}

		// "Bearer " prefix'ini kaldır
		tokenString := authHeader
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenString = authHeader[7:]
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("username", claims["username"])
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}
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
	err := db.DB.QueryRow("SELECT id, title, email, address, weekday_hours, saturday_hours, sunday_hours FROM contact ORDER BY id LIMIT 1").
		Scan(&contact.ID, &contact.Title, &contact.Email, &contact.Address, &contact.WeekdayHours, &contact.SaturdayHours, &contact.SundayHours)

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
			"INSERT INTO contact (title, email, address, weekday_hours, saturday_hours, sunday_hours) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
			contact.Title, contact.Email, contact.Address, contact.WeekdayHours, contact.SaturdayHours, contact.SundayHours,
		).Scan(&contact.ID)
	} else {
		// Varolan kaydı güncelle
		_, err = db.DB.Exec(
			"UPDATE contact SET title = $1, email = $2, address = $3, weekday_hours = $4, saturday_hours = $5, sunday_hours = $6 WHERE id = $7",
			contact.Title, contact.Email, contact.Address, contact.WeekdayHours, contact.SaturdayHours, contact.SundayHours, existingID,
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

func getHeroHandler(c *gin.Context) {
	// Implementation of getHeroHandler
}

func updateHeroHandler(c *gin.Context) {
	// Implementation of updateHeroHandler
}

func getFooterHandler(c *gin.Context) {
	// Implementation of getFooterHandler
}

func updateFooterHandler(c *gin.Context) {
	// Implementation of updateFooterHandler
}

// Users endpoints
func getUsersHandler(c *gin.Context) {
	rows, err := db.DB.Query("SELECT id, username, email, role_id FROM users")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.RoleID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		users = append(users, u)
	}

	c.JSON(http.StatusOK, users)
}

func createUserHandler(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz kullanıcı verisi"})
		return
	}

	// Zorunlu alanları kontrol et
	if user.Username == "" || user.Email == "" || user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kullanıcı adı, e-posta ve şifre zorunludur"})
		return
	}

	// Varsayılan role_id değerini ayarla
	if user.RoleID == 0 {
		user.RoleID = 2 // Varsayılan kullanıcı rolü
	}

	// Şifreyi hashle
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Şifre işlenirken hata oluştu"})
		return
	}

	// Kullanıcıyı veritabanına ekle
	err = db.DB.QueryRow(
		"INSERT INTO users (username, email, password, role_id) VALUES ($1, $2, $3, $4) RETURNING id",
		user.Username, user.Email, string(hashedPassword), user.RoleID,
	).Scan(&user.ID)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			if strings.Contains(err.Error(), "users_username_key") {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Bu kullanıcı adı zaten kullanılıyor"})
			} else if strings.Contains(err.Error(), "users_email_key") {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Bu e-posta adresi zaten kullanılıyor"})
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Bu bilgiler zaten kullanılıyor"})
			}
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kullanıcı oluşturulurken hata oluştu: " + err.Error()})
		return
	}

	// Şifreyi response'dan temizle
	user.Password = ""
	c.JSON(http.StatusCreated, user)
}

func updateUserHandler(c *gin.Context) {
	id := c.Param("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz kullanıcı ID"})
		return
	}

	var updateData struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		RoleID   int    `json:"role_id"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz güncelleme verisi"})
		return
	}

	// Mevcut kullanıcıyı kontrol et
	var user models.User
	err = db.DB.QueryRow("SELECT id, username, email, role_id FROM users WHERE id = $1", userID).
		Scan(&user.ID, &user.Username, &user.Email, &user.RoleID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kullanıcı bulunamadı"})
		return
	}

	// Kullanıcıyı güncelle
	_, err = db.DB.Exec(
		"UPDATE users SET username = $1, email = $2, role_id = $3 WHERE id = $4",
		updateData.Username, updateData.Email, updateData.RoleID, userID,
	)

	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			if strings.Contains(err.Error(), "users_username_key") {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Bu kullanıcı adı zaten kullanılıyor"})
			} else if strings.Contains(err.Error(), "users_email_key") {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Bu e-posta adresi zaten kullanılıyor"})
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Bu bilgiler zaten kullanılıyor"})
			}
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kullanıcı güncellenirken hata oluştu"})
		return
	}

	// Güncellenmiş kullanıcı bilgilerini getir
	err = db.DB.QueryRow("SELECT id, username, email, role_id FROM users WHERE id = $1", userID).
		Scan(&user.ID, &user.Username, &user.Email, &user.RoleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Güncellenmiş kullanıcı bilgileri alınamadı"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func deleteUserHandler(c *gin.Context) {
	id := c.Param("id")

	result, err := db.DB.Exec("DELETE FROM users WHERE id = $1", id)
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
		c.JSON(http.StatusNotFound, gin.H{"error": "Kullanıcı bulunamadı"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Kullanıcı silindi"})
}
