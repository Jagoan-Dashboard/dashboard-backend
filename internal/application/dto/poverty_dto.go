
package dto

type PovertyOverviewResponse struct {
    Tahun                            int         `json:"tahun"`
    AngkaKemiskinan                  float64     `json:"angka_kemiskinan"`
    PerubahanAngkaKemiskinan         *float64    `json:"perubahan_angka_kemiskinan"`
    IndeksKedalamanKemiskinan        float64     `json:"indeks_kedalaman_kemiskinan"`
    PerubahanIndeksKedalaman         *float64    `json:"perubahan_indeks_kedalaman"`
    IndeksKeparahanKemiskinan        float64     `json:"indeks_keparahan_kemiskinan"`
    PerubahanIndeksKeparahan         *float64    `json:"perubahan_indeks_keparahan"`
    IPM                              float64     `json:"ipm"`
    PerubahanIPM                     *float64    `json:"perubahan_ipm"`
    IndeksGini                       float64     `json:"indeks_gini"`
    PerubahanIndeksGini              *float64    `json:"perubahan_indeks_gini"`
    PengeluaranPerKapita             float64     `json:"pengeluaran_per_kapita"`
    PerubahanPengeluaran             *float64    `json:"perubahan_pengeluaran"`
    UmurHarapanHidup                 float64     `json:"umur_harapan_hidup"`
    PerubahanUHH                     *float64    `json:"perubahan_uhh"`
    GarisKemiskinan                  float64     `json:"garis_kemiskinan"`
    PerubahanGarisKemiskinan         *float64    `json:"perubahan_garis_kemiskinan"`
    TrendIndeksKedalamanKemiskinan   []TrendData `json:"trend_indeks_kedalaman_kemiskinan"`
    TrendIndeksKeparahanKemiskinan   []TrendData `json:"trend_indeks_keparahan_kemiskinan"`
}