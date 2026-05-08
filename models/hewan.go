package models

import (
	"time"

	"gorm.io/gorm"
)

type Hewan struct {
	gorm.Model
	KodeHewan   string   `gorm:"uniqueIndex;not null" json:"kode_hewan"`
	Tipe        string   `gorm:"type:enum('qurban','sedekah');default:'qurban'" json:"tipe"`
	JenisHewan  string   `gorm:"type:enum('sapi','kambing');not null" json:"jenis_hewan"`
	NamaSohibul []string `gorm:"serializer:json;not null" json:"nama_sohibul"`
	Catatan     string   `gorm:"type:text" json:"catatan"`

	WaktuMulaiJagal         *time.Time `json:"waktu_mulai_jagal"`
	WaktuSelesaiJagal       *time.Time `json:"waktu_selesai_jagal"`
	WaktuMulaiKuliti        *time.Time `json:"waktu_mulai_kuliti"`
	WaktuSelesaiKuliti      *time.Time `json:"waktu_selesai_kuliti"`
	WaktuMulaiCacahDaging   *time.Time `json:"waktu_mulai_cacah_daging"`
	WaktuSelesaiCacahDaging *time.Time `json:"waktu_selesai_cacah_daging"`
	WaktuMulaiCacahTulang   *time.Time `json:"waktu_mulai_cacah_tulang"`
	WaktuSelesaiCacahTulang *time.Time `json:"waktu_selesai_cacah_tulang"`
	WaktuMulaiPacking       *time.Time `json:"waktu_mulai_packing"`
	WaktuSelesaiPacking     *time.Time `json:"waktu_selesai_packing"`

	KantongPacking *int    `gorm:"default:null" json:"kantong_packing"`
	BeratDaging    float64 `gorm:"type:decimal(10,2);default:0" json:"berat_daging"`
	BeratTulang    float64 `gorm:"type:decimal(10,2);default:0" json:"berat_tulang"`

	CekKepala     bool `gorm:"default:false" json:"cek_kepala"`
	CekKaki       bool `gorm:"default:false" json:"cek_kaki"`
	CekKulit      bool `gorm:"default:false" json:"cek_kulit"`
	CekEkor       bool `gorm:"default:false" json:"cek_ekor"`
	CekDistribusi bool `gorm:"default:false" json:"cek_distribusi"`

	PengawasID uint `json:"pengawas_id"`
	Pengawas   User `gorm:"foreignKey:PengawasID" json:"pengawas"`
}
