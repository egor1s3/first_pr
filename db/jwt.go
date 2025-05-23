package db_m

import (
	"fmt"
	"log"
	"main/models"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var users = map[string]models.User{}

func Get_env() []byte {
	return []byte(os.Getenv("JWT_PASSWORD"))
}

func Login(users []models.User, username, password string) (string, error) {

	for _, name := range users {
		log.Println(users)
		log.Println(name.Username, username)
		log.Println(name.Password, password)
		if name.Username == username {
			if err := bcrypt.CompareHashAndPassword([]byte(name.Password), []byte(password)); err != nil {
				return "", fmt.Errorf("invalid password")
			}
			token, err := GenerateJWT(name.ID)
			Checkjwt(token)
			if err != nil {
				return "", err
			}
			return token, nil
		}

	}

	return "", fmt.Errorf("user not found")

	// Проверка пароля

	// Генерация JWT

}

func GenerateJWT(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString(Get_env())
}

func Checkjwt(token string) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return Get_env(), nil
	})
	if err != nil {
		log.Fatal("JWT verification failed:", err)
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		fmt.Println("Valid JWT. User ID:", claims["sub"])
	} else {
		fmt.Println("Invalid JWT")
	}
}

func Lmain() {
	CreateDBIfNotExists()
	db := DBInit()
	// Регистрация пользователя
	user, err := Register(db, "alice", "password123")
	fmt.Println(users)
	fmt.Println(user)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Registered: %+v\n", user)

	// Логин и получение JWT
	//	token, err := Login("alice", "password123")
	if err != nil {
		log.Fatal(err)
	}
	//	fmt.Println("JWT Token:", token)

	// Проверка токена
	//	Checkjwt(token)

}
