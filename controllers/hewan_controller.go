package controllers

import (
	"net/http"
	"qurban/config"
	"qurban/models"

	"github.com/gin-gonic/gin"
)

func CreateHewan(c *gin.Context) {
	var input struct {
		KodeHewan   string   `json:"kode_hewan" binding:"required"`
		Tipe        string   `json:"tipe" binding:"required,oneof=qurban sedekah"`
		JenisHewan  string   `json:"jenis_hewan" binding:"required,oneof=sapi kambing"`
		NamaSohibul []string `json:"nama_sohibul" binding:"required"`
		PengawasID  uint     `json:"pengawas_id" binding:"required"`
		Catatan     string   `json:"catatan"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid"})
		return
	}

	hewan := models.Hewan{
		KodeHewan: input.KodeHewan, Tipe: input.Tipe,
		JenisHewan: input.JenisHewan, NamaSohibul: input.NamaSohibul,
		PengawasID: input.PengawasID, Catatan: input.Catatan,
	}

	if err := config.DB.Create(&hewan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data hewan"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Data hewan berhasil diinput", "data": hewan})
}

func GetHewan(c *gin.Context) {
	var hewan []models.Hewan
	userRole := c.GetString("userRole")
	userID := c.GetUint("userID")
	search := c.Query("search")
	tipe := c.Query("tipe")
	jenis := c.Query("jenis_hewan")
	qPengawas := c.Query("pengawas_id")

	query := config.DB.Preload("Pengawas")

	if userRole == string(models.RolePengawas) {
		query = query.Where("pengawas_id = ?", userID)
	} else if qPengawas != "" {
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
		WHEN waktu_mulai_jagal IS NOT NULL AND waktu_selesai_jagal IS NULL THEN 1
		WHEN waktu_mulai_jagal IS NULL THEN 2
		ELSE 3 END ASC`)
	query = query.Order("jenis_hewan DESC").Order("tipe ASC").Order("kode_hewan ASC")

	if err := query.Find(&hewan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data hewan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Berhasil", "data": hewan})
}

func UpdateHewan(c *gin.Context) {
	id := c.Param("id")
	var hewan models.Hewan

	if err := config.DB.First(&hewan, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data hewan tidak ditemukan"})
		return
	}

	var input struct {
		KodeHewan   string   `json:"kode_hewan"`
		Tipe        string   `json:"tipe"`
		JenisHewan  string   `json:"jenis_hewan"`
		NamaSohibul []string `json:"nama_sohibul"`
		PengawasID  uint     `json:"pengawas_id"`
		Catatan     string   `json:"catatan"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid"})
		return
	}

	hewan.KodeHewan = input.KodeHewan
	hewan.Tipe = input.Tipe
	hewan.JenisHewan = input.JenisHewan
	hewan.NamaSohibul = input.NamaSohibul
	hewan.PengawasID = input.PengawasID
	hewan.Catatan = input.Catatan

	if err := config.DB.Save(&hewan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate data hewan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data hewan berhasil diperbarui"})
}

func DeleteHewan(c *gin.Context) {
	id := c.Param("id")
	var hewan models.Hewan

	if err := config.DB.First(&hewan, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data hewan tidak ditemukan"})
		return
	}

	if hewan.WaktuMulaiJagal != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Hewan sudah diproses, data tidak dapat dihapus"})
		return
	}

	if err := config.DB.Delete(&hewan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus data hewan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data hewan berhasil dihapus"})
}