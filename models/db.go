package models

import (
	"fmt"
	"log"

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

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func DB_INIT() {
	// Строка подключения
	dsn := "host=localhost user=username password=password dbname=dbname port=5432 sslmode=disable TimeZone=Europe/Moscow"

	// Подключение к БД
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	fmt.Println("Successfully connected to PostgreSQL with GORM!")

	// Автомиграция (создание таблиц)
	err = db.AutoMigrate(&Product{})
	if err != nil {
		log.Fatal("AutoMigrate failed:", err)
	}

	// Создание записи
	db.Create(&Product{Code: "D42", Price: 100})

	// Чтение записи
	var product Product
	db.First(&product, "code = ?", "D42") // Найти продукт с кодом D42
	fmt.Printf("Product: %+v\n", product)

	// Обновление записи
	db.Model(&product).Update("Price", 200)

	// Удаление записи
	db.Delete(&product)
}
