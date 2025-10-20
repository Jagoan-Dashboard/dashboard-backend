
package dto

type TrendData struct {
    Tahun int     `json:"tahun"`
    Nilai float64 `json:"nilai"`
}

type EkonomiOverviewResponse struct {
    Tahun                    int         `json:"tahun"`
    LajuPertumbuhanEkonomi   float64     `json:"laju_pertumbuhan_ekonomi"`
    PerubahanLPE             *float64    `json:"perubahan_lpe"` 
    PertanianPDRB            float64     `json:"pertanian_pdrb"`
    PerubahanPertanian       *float64    `json:"perubahan_pertanian"`
    PengolahanPDRB           float64     `json:"pengolahan_pdrb"`
    PerubahanPengolahan      *float64    `json:"perubahan_pengolahan"`
    ICOR                     float64     `json:"icor"`
    PerubahanICOR            *float64    `json:"perubahan_icor"`
    ILOR                     float64     `json:"ilor"`
    PerubahanILOR            *float64    `json:"perubahan_ilor"`
    Inflasi                  float64     `json:"inflasi"`
    PerubahanInflasi         *float64    `json:"perubahan_inflasi"`
    TrendLajuPertumbuhan     []TrendData `json:"trend_laju_pertumbuhan"`
    TrendInflasi             []TrendData `json:"trend_inflasi"`
}