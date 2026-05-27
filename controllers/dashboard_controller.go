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
	search := c.Query("search")
	tipe := c.Query("tipe")
	jenis := c.Query("jenis_hewan")
	qPengawas := c.Query("pengawas_id")

	query := config.DB.Preload("Pengawas")

	if qPengawas != "" {
		query = query.Where("pengawas_id = ?", qPengawas)
	}
	if tipe != "" {
		query = query.Where("tipe = ?", tipe)
	}
	if jenis != "" {
		query = query.Where("jenis_hewan = ?", jenis)
	}
	if search != "" {
		term := "%" + search + "%"
		query = query.Where(
			config.DB.Where("kode_hewan LIKE ?", term).
				Or("nama_sohibul LIKE ?", term).
				Or("catatan LIKE ?", term),
		)
	}

	query = query.Order(`CASE
		WHEN waktu_mulai_jagal IS NOT NULL AND waktu_selesai_kuliti IS NULL THEN 1
		WHEN waktu_mulai_jagal IS NULL THEN 2
		ELSE 3 END ASC`)
	query = query.Order("waktu_mulai_jagal ASC").Order("kode_hewan ASC")

	if err := query.Find(&hewan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data hewan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Berhasil", "data": hewan})
}
