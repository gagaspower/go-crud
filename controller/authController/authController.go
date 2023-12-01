package authcontroller

import (
	"fmt"
	"learn-go/config"
	"learn-go/models"
	token "learn-go/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(c *gin.Context) {

	var input LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := LoginCheck(input.Email, input.Password)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email or password is incorrect."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})

}

func LoginCheck(email string, password string) (string, error) {
	var err error
	var user models.Users

	err = config.DB.Where("email = ?", email).First(&user).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			fmt.Println("User not found:", email)
			return "", fmt.Errorf("user with email %s not found", email)
		}
		fmt.Println("Error querying database:", err)
		return "", err
	}

	fmt.Println("Found user:", user.Email, user.Password)

	err = VerifyPassword(password, user.Password)

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		fmt.Println("Password mismatch")
		return "", err
	}

	fmt.Println("Authentication successful")

	token, err := token.GenerateToken(uint(user.Id))

	if err != nil {
		fmt.Println("Error generating token:", err)
		return "", err
	}

	return token, nil
}

func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
