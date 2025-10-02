package seeds

import (
	"building-report-backend/internal/domain/entity"
	"building-report-backend/pkg/utils"
	"gorm.io/gorm"
	"time"
)

func SeedWaterResources(db *gorm.DB) error {
	// Get first user as CreatedBy
	var firstUser entity.User
	if err := db.First(&firstUser).Error; err != nil {
		return err
	}

	reports := []entity.WaterResourcesReport{
		{
			ID:                    utils.GenerateULID(),
			ReporterName:          "Agus Setiawan",
			InstitutionUnit:       entity.InstitutionUnitDinas,
			PhoneNumber:           "081234567890",
			ReportDateTime:        time.Now().AddDate(0, 0, -7),
			IrrigationAreaName:    "DI Karangjati",
			IrrigationType:        entity.IrrigationSaluranSekunder,
			Latitude:              -7.4040,
			Longitude:             111.4464,
			DamageType:            entity.DamageRetakBocor,
			DamageLevel:           entity.DamageLevelBerat,
			EstimatedLength:       50.5,
			EstimatedWidth:        2.5,
			EstimatedVolume:       126.25,
			AffectedRiceFieldArea: 15.5,
			AffectedFarmersCount:  25,
			UrgencyCategory:       entity.UrgencyCategoryMendesak,
			Status:                entity.WaterResourceStatusPending,
			Notes:                 "Saluran sekunder bocor parah, perlu perbaikan segera untuk musim tanam",
			CreatedBy:             firstUser.ID,
			CreatedAt:             time.Now().AddDate(0, 0, -7),
			UpdatedAt:             time.Now().AddDate(0, 0, -7),
		},
		{
			ID:                    utils.GenerateULID(),
			ReporterName:          "Yuni Astuti",
			InstitutionUnit:       entity.InstitutionUnitDesa,
			PhoneNumber:           "082345678901",
			ReportDateTime:        time.Now().AddDate(0, 0, -12),
			IrrigationAreaName:    "Pintu Air Paron",
			IrrigationType:        entity.IrrigationPintuAir,
			Latitude:              -7.3950,
			Longitude:             111.4350,
			DamageType:            entity.DamagePintuAirMacet,
			DamageLevel:           entity.DamageLevelSedang,
			EstimatedLength:       0,
			EstimatedWidth:        0,
			EstimatedVolume:       0,
			AffectedRiceFieldArea: 8.3,
			AffectedFarmersCount:  12,
			UrgencyCategory:       entity.UrgencyCategoryRutin,
			Status:                entity.WaterResourceStatusVerified,
			Notes:                 "Pintu air perlu penggantian komponen mekanis yang sudah karatan",
			HandlingRecommendation: "Penggantian gear box dan pengecatan ulang struktur baja",
			EstimatedBudget:       15000000,
			CreatedBy:             firstUser.ID,
			CreatedAt:             time.Now().AddDate(0, 0, -12),
			UpdatedAt:             time.Now().AddDate(0, 0, -5),
		},
		{
			ID:                    utils.GenerateULID(),
			ReporterName:          "Budi Hartono",
			InstitutionUnit:       entity.InstitutionUnitKecamatan,
			PhoneNumber:           "083456789012",
			ReportDateTime:        time.Now().AddDate(0, 0, -4),
			IrrigationAreaName:    "Bendung Geneng",
			IrrigationType:        entity.IrrigationBendung,
			Latitude:              -7.3800,
			Longitude:             111.4200,
			DamageType:            entity.DamageStrukturBetonRusak,
			DamageLevel:           entity.DamageLevelBerat,
			EstimatedLength:       15.0,
			EstimatedWidth:        5.0,
			EstimatedVolume:       75.0,
			AffectedRiceFieldArea: 45.8,
			AffectedFarmersCount:  78,
			UrgencyCategory:       entity.UrgencyCategoryMendesak,
			Status:                entity.WaterResourceStatusInProgress,
			Notes:                 "Struktur beton bendung mengalami kerusakan serius, berisiko runtuh",
			HandlingRecommendation: "Perbaikan struktural dengan grouting dan penguatan balok",
			EstimatedBudget:       250000000,
			CreatedBy:             firstUser.ID,
			CreatedAt:             time.Now().AddDate(0, 0, -4),
			UpdatedAt:             time.Now().AddDate(0, 0, -2),
		},
		{
			ID:                    utils.GenerateULID(),
			ReporterName:          "Siti Nurjanah",
			InstitutionUnit:       entity.InstitutionUnitDinas,
			PhoneNumber:           "084567890123",
			ReportDateTime:        time.Now().AddDate(0, 0, -18),
			IrrigationAreaName:    "Embung Mantingan",
			IrrigationType:        entity.IrrigationEmbungDam,
			Latitude:              -7.4100,
			Longitude:             111.4600,
			DamageType:            entity.DamageSedimentasiTinggi,
			DamageLevel:           entity.DamageLevelRingan,
			EstimatedLength:       0,
			EstimatedWidth:        0,
			EstimatedVolume:       5000,
			AffectedRiceFieldArea: 5.2,
			AffectedFarmersCount:  8,
			UrgencyCategory:       entity.UrgencyCategoryRutin,
			Status:                entity.WaterResourceStatusCompleted,
			Notes:                 "Pengerukan embung sudah selesai dilaksanakan",
			HandlingRecommendation: "Pengerukan sedimen secara berkala setiap tahun",
			EstimatedBudget:       35000000,
			CreatedBy:             firstUser.ID,
			CreatedAt:             time.Now().AddDate(0, 0, -18),
			UpdatedAt:             time.Now().AddDate(0, 0, -1),
		},
		{
			ID:                    utils.GenerateULID(),
			ReporterName:          "Eko Wahyudi",
			InstitutionUnit:       entity.InstitutionUnitDesa,
			PhoneNumber:           "085678901234",
			ReportDateTime:        time.Now().AddDate(0, 0, -2),
			IrrigationAreaName:    "Saluran Ketanggi",
			IrrigationType:        entity.IrrigationSaluranSekunder,
			Latitude:              -7.4200,
			Longitude:             111.4500,
			DamageType:            entity.DamageTersumbatSampah,
			DamageLevel:           entity.DamageLevelSedang,
			EstimatedLength:       120.0,
			EstimatedWidth:        1.2,
			EstimatedVolume:       144.0,
			AffectedRiceFieldArea: 6.5,
			AffectedFarmersCount:  15,
			UrgencyCategory:       entity.UrgencyCategoryMendesak,
			Status:                entity.WaterResourceStatusPending,
			Notes:                 "Saluran tersumbat sampah dan sedimentasi tinggi",
			HandlingRecommendation: "Pembersihan manual dan normalisasi saluran",
			EstimatedBudget:       8000000,
			CreatedBy:             firstUser.ID,
			CreatedAt:             time.Now().AddDate(0, 0, -2),
			UpdatedAt:             time.Now().AddDate(0, 0, -2),
		},
	}

	for _, report := range reports {
		if err := db.Create(&report).Error; err != nil {
			return err
		}
	}

	return nil
}