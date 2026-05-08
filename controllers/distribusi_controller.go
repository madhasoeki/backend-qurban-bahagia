package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"qurban/config"
	"qurban/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetDistribusiUser(c *gin.Context) {
	userID := c.Param("user_id")

	var dist models.Distribusi
	err := config.DB.Where("user_id = ?", userID).First(&dist).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Return zero-value placeholder so frontend can render without errors
			c.JSON(http.StatusOK, gin.H{
				"message": "Belum ada distribusi",
				"data":    gin.H{"JumlahKantong": 0, "UserID": userID},
			})
			return
		}
		sendError(c, "Gagal mengambil data distribusi")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Berhasil", "data": dist})
}

func GetAllDistribusi(c *gin.Context) {
	var list []models.Distribusi

	if err := config.DB.Preload("User").Find(&list).Error; err != nil {
		sendError(c, "Gagal mengambil data distribusi")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Berhasil", "data": list})
}

func UpdateDistribusi(c *gin.Context) {
	targetUserIDStr := c.Param("user_id")
	targetUserID, err := strconv.ParseUint(targetUserIDStr, 10, 32)
	if err != nil {
		sendError(c, "ID User tidak valid")
		return
	}

	userRole := c.GetString("userRole")
	currentUserID := c.GetUint("userID")

	// Ownership validation: distribusi role can only update own data
	if userRole == "distribusi" && uint(targetUserID) != currentUserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Anda hanya dapat mengupdate data distribusi Anda sendiri"})
		return
	}

	var input struct {
		Penambahan *int `json:"penambahan" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kirimkan parameter 'penambahan'"})
		return
	}

	var dist models.Distribusi
	if err := config.DB.FirstOrCreate(&dist, models.Distribusi{UserID: uint(targetUserID)}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengakses data distribusi"})
		return
	}

	dist.JumlahKantong += *input.Penambahan
	if dist.JumlahKantong < 0 {
		dist.JumlahKantong = 0
	}

	if err := config.DB.Save(&dist).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data distribusi"})
		return
	}

	Broadcast <- gin.H{"action": "UPDATE_DISTRIBUSI", "data": dist}

	go func() {
		summary := calculateDashboardSummary()
		Broadcast <- gin.H{"action": "UPDATE_DASHBOARD", "data": summary}
	}()

	c.JSON(http.StatusOK, gin.H{"message": "Data distribusi berhasil diupdate", "data": dist})
}
