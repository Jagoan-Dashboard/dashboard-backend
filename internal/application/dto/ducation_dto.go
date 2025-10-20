
package dto

type EducationOverviewResponse struct {
    Tahun                        int         `json:"tahun"`
    RataRataLamaSekolah          float64     `json:"rata_rata_lama_sekolah"`
    PerubahanRataRataLamaSekolah *float64    `json:"perubahan_rata_rata_lama_sekolah"`
    HarapanLamaSekolah           float64     `json:"harapan_lama_sekolah"`
    PerubahanHarapanLamaSekolah  *float64    `json:"perubahan_harapan_lama_sekolah"`
    ProporsiPendidikanTinggi     float64     `json:"proporsi_pendidikan_tinggi"`
    PerubahanProporsiPendidikan  *float64    `json:"perubahan_proporsi_pendidikan"`
    TrendRataRataLamaSekolah     []TrendData `json:"trend_rata_rata_lama_sekolah"`
    TrendHarapanLamaSekolah      []TrendData `json:"trend_harapan_lama_sekolah"`
}