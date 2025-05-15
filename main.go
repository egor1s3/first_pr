package main

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")

		// Проверка учётных данных (в реальном приложении - проверка в БД)
		if username == "admin" && password == "123" { //здесь перенести данные из другой папки  для проверки
			token := "generated_jwt_token_here" // В реальности генерируем JWT

			// Устанавливаем cookie
			c.SetCookie("auth_token", token, 3600, "/", "", false, true)

			c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
			return
		}

		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})

		Logging(c.Request.UserAgent(), c.Request.URL, c.Request.Method, "login_post")
	}
}

func main() {

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

	router.POST("/register", authMiddleware()) // 1сделать просто обновление нового пользователя
	// 2 подумать про случай когда уже существует пользователь

	router.POST("/login", authMiddleware())

	router.GET("/protected", getcook("auth_token"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to protected area!"})
	})

	router.Run(":8080")
}
