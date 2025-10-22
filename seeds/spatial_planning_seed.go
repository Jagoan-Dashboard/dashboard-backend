package seeds

import (
	"building-report-backend/internal/domain/entity"
	"building-report-backend/pkg/utils"
	"time"

	"gorm.io/gorm"
)

func SeedSpatialPlanning(db *gorm.DB) error {
	// Get first user as CreatedBy
	var firstUser entity.User
	if err := db.First(&firstUser).Error; err != nil {
		return err
	}

	reports := []entity.SpatialPlanningReport{
		{
			ID:                  utils.GenerateULID(),
			ReporterName:        "Bambang Suryono",
			Institution:         entity.InstitutionDinas,
			PhoneNumber:         "081234567890",
			ReportDateTime:      time.Now().AddDate(0, 0, -10),
			AreaDescription:     "Pembangunan ruko tanpa izin di kawasan perumahan yang melanggar tata ruang",
			AreaCategory:        entity.AreaPermukiman,
			ViolationType:       entity.ViolationTanpaIzin,
			ViolationLevel:      entity.ViolationBerat,
			EnvironmentalImpact: entity.ImpactGangguanAktivitas,
			UrgencyLevel:        entity.UrgencyMendesak,
			Latitude:            -7.4040,
			Longitude:           111.4464,
			Address:             "Jl. Karangjati Gang 5 No. 23, Ngawi",
			Status:              entity.SpatialStatusPending,
			Notes:               "Perlu tindakan segera karena mengganggu akses jalan warga",
			// CreatedBy:           firstUser.ID,
			CreatedAt: time.Now().AddDate(0, 0, -10),
			UpdatedAt: time.Now().AddDate(0, 0, -10),
		},
		{
			ID:                  utils.GenerateULID(),
			ReporterName:        "Sri Wahyuni",
			Institution:         entity.InstitutionKecamatan,
			PhoneNumber:         "082345678901",
			ReportDateTime:      time.Now().AddDate(0, 0, -5),
			AreaDescription:     "Alih fungsi lahan pertanian menjadi kawasan industri tanpa izin",
			AreaCategory:        entity.AreaTanamanPangan,
			ViolationType:       entity.ViolationAlihFungsiPertanian,
			ViolationLevel:      entity.ViolationSedang,
			EnvironmentalImpact: entity.ImpactKualitasRuang,
			UrgencyLevel:        entity.UrgencyBiasa,
			Latitude:            -7.3950,
			Longitude:           111.4350,
			Address:             "Kawasan Pertanian Paron Blok C, Ngawi",
			Status:              entity.SpatialStatusReviewing,
			Notes:               "Sedang dalam tahap verifikasi dokumen kepemilikan",
			CreatedAt:           time.Now().AddDate(0, 0, -5),
			UpdatedAt:           time.Now().AddDate(0, 0, -5),
		},
		{
			ID:                  utils.GenerateULID(),
			ReporterName:        "Andi Wijaya",
			Institution:         entity.InstitutionDesa,
			PhoneNumber:         "083456789012",
			ReportDateTime:      time.Now().AddDate(0, 0, -15),
			AreaDescription:     "Bangunan di sempadan sungai yang berpotensi menimbulkan banjir",
			AreaCategory:        entity.AreaPermukiman,
			ViolationType:       entity.ViolationSempadanSungai,
			ViolationLevel:      entity.ViolationBerat,
			EnvironmentalImpact: entity.ImpactBanjirLongsor,
			UrgencyLevel:        entity.UrgencyMendesak,
			Latitude:            -7.3800,
			Longitude:           111.4200,
			Address:             "Sempadan Sungai Bengawan Solo, Desa Geneng",
			Status:              entity.SpatialStatusProcessing,
			Notes:               "Koordinasi dengan Dinas Lingkungan Hidup dan BPBD",
			// CreatedBy:           firstUser.ID,
			CreatedAt: time.Now().AddDate(0, 0, -15),
			UpdatedAt: time.Now().AddDate(0, 0, -15),
		},
		{
			ID:                  utils.GenerateULID(),
			ReporterName:        "Fitri Handayani",
			Institution:         entity.InstitutionDinas,
			PhoneNumber:         "084567890123",
			ReportDateTime:      time.Now().AddDate(0, 0, -3),
			AreaDescription:     "Pembangunan villa di kawasan resapan air tanpa IMB",
			AreaCategory:        entity.AreaPariwisata,
			ViolationType:       entity.ViolationTanpaIzin,
			ViolationLevel:      entity.ViolationSedang,
			EnvironmentalImpact: entity.ImpactKualitasRuang,
			UrgencyLevel:        entity.UrgencyBiasa,
			Latitude:            -7.4100,
			Longitude:           111.4600,
			Address:             "Kawasan Wisata Mantingan, Ngawi",
			Status:              entity.SpatialStatusPending,
			Notes:               "Menunggu hasil survey lapangan dari tim teknis",
			// CreatedBy:           firstUser.ID,
			CreatedAt: time.Now().AddDate(0, 0, -3),
			UpdatedAt: time.Now().AddDate(0, 0, -3),
		},
		{
			ID:                  utils.GenerateULID(),
			ReporterName:        "Hendra Saputra",
			Institution:         entity.InstitutionKecamatan,
			PhoneNumber:         "085678901234",
			ReportDateTime:      time.Now().AddDate(0, 0, -20),
			AreaDescription:     "Alih fungsi ruang terbuka hijau menjadi bangunan komersial",
			AreaCategory:        entity.AreaPermukiman,
			ViolationType:       entity.ViolationAlihFungsiRTH,
			ViolationLevel:      entity.ViolationRingan,
			EnvironmentalImpact: entity.ImpactKualitasRuang,
			UrgencyLevel:        entity.UrgencyBiasa,
			Latitude:            -7.4200,
			Longitude:           111.4500,
			Address:             "RTH Ketanggi, Ngawi",
			Status:              entity.SpatialStatusResolved,
			Notes:               "Sudah ditindaklanjuti dengan pemberian teguran tertulis",
			// CreatedBy:           firstUser.ID,
			CreatedAt: time.Now().AddDate(0, 0, -20),
			UpdatedAt: time.Now().AddDate(0, 0, -1),
		},
	}

	for _, report := range reports {
		if err := db.Create(&report).Error; err != nil {
			return err
		}
	}

	return nil
}
