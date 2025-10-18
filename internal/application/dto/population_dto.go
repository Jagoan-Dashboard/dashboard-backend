package dto

type PopulationOverviewResponse struct {
    Tahun                     int         `json:"tahun"`
    KepadatanPenduduk         float64     `json:"kepadatan_penduduk"`
    PerubahanKepadatan        *float64    `json:"perubahan_kepadatan"`
    RasioKetergantungan       float64     `json:"rasio_ketergantungan"`
    PerubahanRasio            *float64    `json:"perubahan_rasio"`
    PendudukProduktif         float64     `json:"penduduk_produktif"`
    PerubahanProduktif        *float64    `json:"perubahan_produktif"`
    PendudukNonProduktif      float64     `json:"penduduk_non_produktif"`
    PerubahanNonProduktif     *float64    `json:"perubahan_non_produktif"`
    TrendKepadatanPenduduk    []TrendData `json:"trend_kepadatan_penduduk"`
    TrendRasioKetergantungan  []TrendData `json:"trend_rasio_ketergantungan"`
}