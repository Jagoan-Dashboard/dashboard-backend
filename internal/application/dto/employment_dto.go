
package dto

type EmploymentOverviewResponse struct {
    Tahun                     int         `json:"tahun"`
    TPT                       float64     `json:"tpt"`
    PerubahanTPT              *float64    `json:"perubahan_tpt"`
    TPAK                      float64     `json:"tpak"`
    PerubahanTPAK             *float64    `json:"perubahan_tpak"`
    TPAKPerempuan             float64     `json:"tpak_perempuan"`
    PerubahanTPAKPerempuan    *float64    `json:"perubahan_tpak_perempuan"`
    UpahMinimum               float64     `json:"upah_minimum"`
    PerubahanUpahMinimum      *float64    `json:"perubahan_upah_minimum"`
    TrendTPT                  []TrendData `json:"trend_tpt"`
    TrendTPAK                 []TrendData `json:"trend_tpak"`
}