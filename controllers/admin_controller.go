package controllers

import (
	"net/http"
	"qurban/config"
	"qurban/models"

	"github.com/gin-gonic/gin"
)

func CreateUser(c *gin.Context) {
	var input struct {
		NamaLengkap string      `json:"nama_lengkap" binding:"required"`
		Username    string      `json:"username" binding:"required"`
		Password    string      `json:"password" binding:"required"`
		Role        models.Role `json:"role" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid"})
		return
	}

	// Whitelist role validation
	if !models.ValidRoles[input.Role] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role tidak valid"})
		return
	}

	user := models.User{
		NamaLengkap: input.NamaLengkap,
		Username:    input.Username,
		Password:    input.Password,
		Role:        input.Role,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat user baru"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User berhasil didaftarkan", "user_id": user.ID})
}

func GetUsers(c *gin.Context) {
	var users []models.User

	roleFilter := c.Query("role")
	searchFilter := c.Query("search")

	query := config.DB

	if roleFilter != "" {
		query = query.Where("role = ?", roleFilter)
	}

	if searchFilter != "" {
		query = query.Where("username LIKE ?", "%"+searchFilter+"%")
	}

	if err := query.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Berhasil", "data": users})
}

func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	if err := config.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
		return
	}

	var input struct {
		Username string      `json:"username"`
		Password string      `json:"password"`
		Role     models.Role `json:"role"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid"})
		return
	}

	user.Username = input.Username
	user.Role = input.Role

	if input.Password != "" {
		user.Password = input.Password
	}

	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User berhasil diperbarui"})
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	if err := config.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
		return
	}

	if err := config.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User berhasil dihapus"})
}
