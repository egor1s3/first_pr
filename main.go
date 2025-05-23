package main

import (
	"log"
	db_m "main/db"
	"main/models"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func getcook(coname string) gin.HandlerFunc {
	return func(c *gin.Context) {

		token, err := c.Cookie(coname)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// Здесь должна быть проверка JWT токена
		if token != "generated_jwt_token_here" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		c.Next()
	}
}
func Logging(useragent string, url *url.URL, method string, msg string) {
	logrus.WithFields(logrus.Fields{
		"user":    useragent,
		"url":     url,
		"method":  method,
		"user_id": 123, // сделать id из файла
		"action":  msg,
	}).Info("Пользователь вошел в систему")
}

func RegMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var users []models.User
		_ = db.Find(&users)
		username := c.PostForm("username")
		password := c.PostForm("password")
		log.Println(users)
		for _, name := range users {
			log.Println(name.Username, username)
			log.Println(name.Password, password)
			if name.Username == username {
				c.JSON(http.StatusOK, gin.H{"message": "Такой пользователь уже есть"})
				return
			}
		}
		_, err := db_m.Register(db, username, password)
		c.JSON(http.StatusOK, gin.H{"username": username, "password": password})
		if err != nil {
			log.Fatal(err)
		}
		Logging(c.Request.UserAgent(), c.Request.URL, c.Request.Method, "register_post")

	}

}
func AuthMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")
		var users []models.User
		_ = db.Find(&users)

		token, err := db_m.Login(users, username, password)
		if err != nil {
			log.Fatal("login error", err)
		}
		if token == "" {
			c.JSON(http.StatusOK, gin.H{"message": "Login unsuccessful:("})
			return
		}
		// Устанавливаем cookie
		c.SetCookie("auth_token", token, 3600, "/", "", false, true)

		c.JSON(http.StatusOK, gin.H{"message": "Login successful"})

		//		c.JSON(http.StatusOK, gin.H{"users": users, "username": username, "password": password})
		// Проверка учётных данных (в реальном приложении - проверка в БД)
		/*		if username == "admin" && password == "123" { //здесь перенести данные из другой папки  для проверки
					token := "generated_jwt_token_here" // В реальности генерируем JWT

					// Устанавливаем cookie
					c.SetCookie("auth_token", token, 3600, "/", "", false, true)

					c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
					return
				}

				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		*/
		Logging(c.Request.UserAgent(), c.Request.URL, c.Request.Method, "login_post")
	}
}

func main() {
	db_m.CreateDBIfNotExists()
	db := db_m.DBInit()
	router := gin.Default()

	// Раздаём статику из templates по URL /static
	router.Static("/index", "./templates")
	// Загружаем HTML
	router.LoadHTMLGlob("templates/*.html")

	// Роуты
	router.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{

			"PageCSS": "/static/style_main.css",
		})

		Logging(c.Request.UserAgent(), c.Request.URL, c.Request.Method, "mainpage")
	})

	router.GET("/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", gin.H{
			"PageCSS": "/static/style_au.css",
		})
		Logging(c.Request.UserAgent(), c.Request.URL, c.Request.Method, "registerpage")
	})

	router.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"PageCSS": "/static/style_au.css",
		})
		Logging(c.Request.UserAgent(), c.Request.URL, c.Request.Method, "loginpage")
	})

	router.POST("/register", RegMiddleware(db)) // 1сделать просто обновление нового пользователя
	// 2 подумать про случай когда уже существует пользователь

	router.POST("/login", AuthMiddleware(db))

	router.GET("/protected", getcook("auth_token"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to protected area!"})
	})

	router.Run(":8080")
}
