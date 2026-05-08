package models

import "time"

type DashboardSummary struct {
	TotalHewan             int64      `json:"total_hewan"`
	TotalHewanSelesai      int64      `json:"total_hewan_selesai"`
	TotalKantongPacking    int        `json:"total_kantong_packing"`
	TotalKantongDistribusi int        `json:"total_kantong_distribusi"`
	WaktuMulai             *time.Time `json:"waktu_mulai"`
	WaktuSelesai           *time.Time `json:"waktu_selesai"`
}
