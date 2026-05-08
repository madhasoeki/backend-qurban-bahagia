package controllers

import (
	"net/http"
	"time"

	"qurban/config"
	"qurban/models"

	"github.com/gin-gonic/gin"
)

func UpdateProgressHewan(c *gin.Context) {
	id := c.Param("id")
	pos := c.Param("pos")
	userRole := c.GetString("userRole")
	userID := c.GetUint("userID")

	var hewan models.Hewan
	if err := config.DB.Preload("Pengawas").First(&hewan, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Hewan tidak ditemukan"})
		return
	}

	if userRole == string(models.RolePengawas) && hewan.PengawasID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak berhak mengelola hewan ini"})
		return
	}

	var input struct {
		Status string `json:"status" binding:"required,oneof=mulai selesai"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status harus 'mulai' atau 'selesai'"})
		return
	}

	// Prerequisite: non-jagal stations require jagal to be completed
	if pos != "jagal" && hewan.WaktuSelesaiJagal == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Proses jagal belum selesai"})
		return
	}

	now := time.Now()

	switch pos {
	case "jagal":
		if input.Status == "mulai" {
			if hewan.WaktuMulaiJagal != nil {
				sendError(c, "Proses jagal sudah dimulai")
				return
			}
			hewan.WaktuMulaiJagal = &now
		} else {
			if hewan.WaktuMulaiJagal == nil || hewan.WaktuSelesaiJagal != nil {
				sendError(c, "Status jagal tidak valid untuk diselesaikan")
				return
			}
			hewan.WaktuSelesaiJagal = &now
		}

	case "kulit":
		if input.Status == "mulai" {
			if hewan.WaktuMulaiKuliti != nil {
				sendError(c, "Proses kulit sudah dimulai")
				return
			}
			hewan.WaktuMulaiKuliti = &now
		} else {
			if hewan.WaktuMulaiKuliti == nil || hewan.WaktuSelesaiKuliti != nil {
				sendError(c, "Status kulit tidak valid untuk diselesaikan")
				return
			}
			hewan.WaktuSelesaiKuliti = &now
		}

	case "cacah_daging":
		if input.Status == "mulai" {
			if hewan.WaktuMulaiCacahDaging != nil {
				sendError(c, "Proses cacah daging sudah dimulai")
				return
			}
			hewan.WaktuMulaiCacahDaging = &now
		} else {
			if hewan.WaktuMulaiCacahDaging == nil || hewan.WaktuSelesaiCacahDaging != nil {
				sendError(c, "Status cacah daging tidak valid untuk diselesaikan")
				return
			}
			hewan.WaktuSelesaiCacahDaging = &now
		}

	case "cacah_tulang":
		if input.Status == "mulai" {
			if hewan.WaktuMulaiCacahTulang != nil {
				sendError(c, "Proses cacah tulang sudah dimulai")
				return
			}
			hewan.WaktuMulaiCacahTulang = &now
		} else {
			if hewan.WaktuMulaiCacahTulang == nil || hewan.WaktuSelesaiCacahTulang != nil {
				sendError(c, "Status cacah tulang tidak valid untuk diselesaikan")
				return
			}
			hewan.WaktuSelesaiCacahTulang = &now
		}

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Pos operasional tidak dikenali"})
		return
	}

	if err := config.DB.Save(&hewan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan progress"})
		return
	}

	broadcastHewanUpdate(hewan)

	c.JSON(http.StatusOK, gin.H{
		"message": "Progress " + pos + " berhasil diperbarui",
		"data":    hewan,
	})
}

func UpdateTimbangHewan(c *gin.Context) {
	id := c.Param("id")
	userRole := c.GetString("userRole")
	userID := c.GetUint("userID")

	var hewan models.Hewan
	if err := config.DB.First(&hewan, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Hewan tidak ditemukan"})
		return
	}

	if userRole == string(models.RolePengawas) && hewan.PengawasID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak berhak mengelola hewan ini"})
		return
	}

	if hewan.WaktuSelesaiJagal == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Proses jagal belum selesai"})
		return
	}

	var input struct {
		BeratDaging *float64 `json:"berat_daging" binding:"required"`
		BeratTulang *float64 `json:"berat_tulang" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kirimkan berat_daging dan berat_tulang berupa angka"})
		return
	}

	hewan.BeratDaging = *input.BeratDaging
	hewan.BeratTulang = *input.BeratTulang

	if err := config.DB.Save(&hewan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data timbangan"})
		return
	}

	broadcastHewanUpdate(hewan)

	c.JSON(http.StatusOK, gin.H{"message": "Data timbangan berhasil disimpan", "data": hewan})
}

func UpdatePackingHewan(c *gin.Context) {
	id := c.Param("id")

	var hewan models.Hewan
	if err := config.DB.First(&hewan, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Hewan tidak ditemukan"})
		return
	}

	if hewan.WaktuSelesaiJagal == nil {
		sendError(c, "Proses jagal belum selesai")
		return
	}

	var input struct {
		Status       string `json:"status" binding:"required,oneof=mulai selesai"`
		TotalKantong *int   `json:"total_kantong"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		sendError(c, "Payload tidak valid")
		return
	}

	now := time.Now()

	switch input.Status {
	case "mulai":
		if hewan.WaktuMulaiPacking != nil {
			sendError(c, "Proses packing sudah dimulai")
			return
		}
		hewan.WaktuMulaiPacking = &now

	case "selesai":
		if hewan.WaktuMulaiPacking == nil {
			sendError(c, "Proses packing belum dimulai")
			return
		}
		if input.TotalKantong == nil {
			sendError(c, "Total kantong wajib diisi saat menyelesaikan packing")
			return
		}
		if hewan.WaktuSelesaiPacking == nil {
			hewan.WaktuSelesaiPacking = &now
		}
		hewan.KantongPacking = input.TotalKantong
	}

	if err := config.DB.Save(&hewan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data packing"})
		return
	}

	broadcastHewanUpdate(hewan)

	c.JSON(http.StatusOK, gin.H{"message": "Data packing berhasil diperbarui", "data": hewan})
}

func UpdateKelengkapanHewan(c *gin.Context) {
	id := c.Param("id")
	userRole := c.GetString("userRole")
	userID := c.GetUint("userID")

	var hewan models.Hewan
	if err := config.DB.First(&hewan, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Hewan tidak ditemukan"})
		return
	}

	if userRole == "pengawas" && hewan.PengawasID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak berhak mengelola hewan ini"})
		return
	}

	var input struct {
		CekKepala     *bool `json:"cek_kepala"`
		CekKaki       *bool `json:"cek_kaki"`
		CekKulit      *bool `json:"cek_kulit"`
		CekEkor       *bool `json:"cek_ekor"`
		CekDistribusi *bool `json:"cek_distribusi"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payload tidak valid"})
		return
	}

	if input.CekKepala != nil {
		hewan.CekKepala = *input.CekKepala
	}
	if input.CekKaki != nil {
		hewan.CekKaki = *input.CekKaki
	}
	if input.CekKulit != nil {
		hewan.CekKulit = *input.CekKulit
	}
	if input.CekEkor != nil {
		hewan.CekEkor = *input.CekEkor
	}
	if input.CekDistribusi != nil {
		hewan.CekDistribusi = *input.CekDistribusi
	}

	if err := config.DB.Save(&hewan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan kelengkapan"})
		return
	}

	broadcastHewanUpdate(hewan)

	c.JSON(http.StatusOK, gin.H{"message": "Data kelengkapan berhasil diperbarui", "data": hewan})
}

func sendError(c *gin.Context, msg string) {
	c.JSON(http.StatusBadRequest, gin.H{"error": msg})
}
