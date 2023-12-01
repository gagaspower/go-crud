package userController

import (
	"learn-go/config"
	"learn-go/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Index(c *gin.Context) {
	var users []models.Users

	config.DB.Select([]string{"id", "nama_user", "email"}).Find(&users)
	c.JSON(http.StatusOK, gin.H{"message": "Success", "data": users})
}

func Create(c *gin.Context) {
	var user models.Users

	if err := c.BindJSON(&user); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Enkripsi password sebelum disimpan
	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to hash password"})
		return
	}

	user.Password = hashedPassword

	config.DB.Create(&user)
	c.JSON(http.StatusOK, gin.H{"message": "Success", "data": user})
}

func Show(c *gin.Context) {
	var user models.Users
	id := c.Param("id")

	if err := config.DB.First(&user, id).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Data tidak ditemukan"})
			return
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Success", "data": user})
}

func Update(c *gin.Context) {
	var user models.Users

	id := c.Param("id")

	// Ambil data user dari database berdasarkan ID
	// Ambil data user dari database berdasarkan ID
	if err := config.DB.First(&user, id).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Data not exists"})
			return
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
	}

	// Simpan password yang sudah ada di database sebelum bind data baru
	existingPassword := user.Password

	// Bind input dari request ke variabel user
	if err := c.BindJSON(&user); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// jika password tidak kosong atau bukan empty string
	if user.Password != "" {
		hashedPassword, err := hashPassword(user.Password)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to hash password"})
			return
		}

		user.Password = hashedPassword
	} else {
		user.Password = existingPassword
	}

	err := config.DB.Model(&user).Where("id = ? ", id).Updates(user).Error

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	user.Password = ""
	c.JSON(http.StatusOK, gin.H{"message": "Update success", "data": user})

}

func Destroy(c *gin.Context) {
	var user models.Users
	id := c.Param("id")
	// Ambil data user dari database berdasarkan ID
	if err := config.DB.First(&user, id).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Data not found!"})
			return
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
	}

	err := config.DB.Where("id=?", id).Delete(&user).Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed delete data", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data has been deleted!"})

}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}
