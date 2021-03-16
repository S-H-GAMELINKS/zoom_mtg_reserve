package main

import (
	"fmt"
	"os"
	"time"
	"github.com/joho/godotenv"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func getJWT() (tokenString string, err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": os.Getenv("ZOOM_API_KEY"),
		"exp": fmt.Sprintf("%d", time.Now().Add(1*time.Minute).Unix()),
	})
	
	tokenString, err = token.SignedString([]byte(os.Getenv("ZOOM_API_SECRET")))
	if err != nil {
		fmt.Println("Error Can not create JWT!")
		tokenString = ""
		return "", err
	}
	
	fmt.Println(tokenString)
	return tokenString, err
}

func main() {
	err := godotenv.Load()
	if err != nil {
	  fmt.Println("Error loading .env file")
	}

	r := gin.Default()

	r.GET("/reserve/zoom/mtgs", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "All reserved Zoom MTG's!",
		})
	})

	r.POST("/reserve/zoom/mtg", func(c *gin.Context) {
		jwt, err := getJWT()
		if err != nil {
			fmt.Println("Error Can not reserved Zoom MTG!")
			jwt = "can not reserved!"
		}
		c.JSON(200, gin.H{
			"message": jwt,
		})
	})
	r.Run()
}
