package controllers

import (
	"net/http"

	"qurban/config"
	"qurban/models"

	"github.com/gin-gonic/gin"
)

// calculateDashboardSummary aggregates real-time stats from the database.
// Called by the HTTP handler and by SSE broadcast goroutines.
func calculateDashboardSummary() models.DashboardSummary {
	var summary models.DashboardSummary

	config.DB.Model(&models.Hewan{}).Count(&summary.TotalHewan)
	config.DB.Model(&models.Hewan{}).Where("waktu_selesai_kuliti IS NOT NULL").Count(&summary.TotalHewanSelesai)

	config.DB.Model(&models.Hewan{}).Select("COALESCE(SUM(kantong_packing), 0)").Scan(&summary.TotalKantongPacking)
	config.DB.Model(&models.Distribusi{}).Select("COALESCE(SUM(jumlah_kantong), 0)").Scan(&summary.TotalKantongDistribusi)

	var countMulai int64
	config.DB.Model(&models.Hewan{}).Where("waktu_mulai_jagal IS NOT NULL").Count(&countMulai)
	if countMulai > 0 {
		config.DB.Model(&models.Hewan{}).Select("MIN(waktu_mulai_jagal)").Scan(&summary.WaktuMulai)
	}

	var belumSelesai int64
	config.DB.Model(&models.Hewan{}).Where("waktu_selesai_kuliti IS NULL").Count(&belumSelesai)
	if belumSelesai == 0 && summary.TotalHewan > 0 {
		config.DB.Model(&models.Hewan{}).Select("MAX(waktu_selesai_kuliti)").Scan(&summary.WaktuSelesai)
	}

	return summary
}

func GetDashboardSummary(c *gin.Context) {
	summary := calculateDashboardSummary()
	c.JSON(http.StatusOK, gin.H{"data": summary})
}

func GetPublicHewan(c *gin.Context) {
	var hewan []models.Hewan

	if err := config.DB.Preload("Pengawas").Find(&hewan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data hewan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Berhasil", "data": hewan})
}
