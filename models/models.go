package models

import (
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

type User struct {
	gorm.Model
	ID       string
	Username string
	Password string // Хранится только хеш!
}
