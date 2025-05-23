package db_m

import (
	"fmt"
	"log"
	"main/models"
	"os"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

/*
services: для докера и бд
  db:
    image: postgres
    environment:
      POSTGRES_USER: myuser
      POSTGRES_PASSWORD: mypassword
      POSTGRES_DB: mydb
    ports:
      - "5432:5432"

*/

func Register(db *gorm.DB, username, password string) (models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, err
	}

	user := models.User{
		ID:       fmt.Sprintf("%d", time.Now().UnixNano()),
		Username: username,
		Password: string(hashedPassword),
	}

	db.Create(&user)
	return user, nil
}

func CreateDBIfNotExists() {
	// Подключение к служебной БД postgres
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Env_err:", err)
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=postgres port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Admin connection failed:", err)
	}

	// Проверка существования БД
	var exists bool
	db.Raw("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = ?)", os.Getenv("DB_NAME")).Scan(&exists)

	if !exists {
		if err := db.Exec("CREATE DATABASE " + os.Getenv("DB_NAME")).Error; err != nil {
			log.Fatal("Failed to create database:", err)
		}
		log.Println("Database created successfully")
	}
}

func DBInit() *gorm.DB {
	// Строка подключения
	//	dsn := "host=localhost user=username password=password dbname=dbname port=5432 sslmode=disable TimeZone=Europe/Moscow"
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Env_err:", err)
	}

	requiredEnv := []string{"DB_HOST", "DB_USER", "DB_PASSWORD", "DB_NAME", "DB_PORT", "DB_SSL_MODE"}
	for _, env := range requiredEnv {
		if os.Getenv(env) == "" {
			log.Fatalf("Environment variable %s is not set", env)
		}
	}

	// Формируем DSN строку
	dsn := fmt.Sprintf("host=%s user=%s password=%s database=%s port=%s sslmode=%s TimeZone=Europe/Moscow",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSL_MODE"),
	)
	fmt.Println("DSN:", dsn)
	// Подключение к БД
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	fmt.Println("Successfully connected to PostgreSQL with GORM!")

	// Настройка пула соединений
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	// Автомиграция (создание таблиц)
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("AutoMigrate failed:", err)
	}

	/*	var fetchedUser User
		db.First(&fetchedUser, "Username = ?", "alice")
		fmt.Printf("User: %+v\n", fetchedUser)

		// Обновление записи
		db.Model(&fetchedUser).Update("Username", "ali")
	*/
	// Удаление записи
	//	db.Delete(&user)
	return db
}
