package seeds

import (
	"building-report-backend/internal/domain/entity"
	"building-report-backend/pkg/utils"
	"log"
	"time"

	"gorm.io/gorm"
)

func timePtr(t time.Time) *time.Time {
	return &t
}

func SeedAgriculture(db *gorm.DB) error {
	now := time.Now()

	// Get first user for CreatedBy
	var firstUser entity.User
	if err := db.First(&firstUser).Error; err != nil {
		return err
	}

	reports := []entity.AgricultureReport{
		// 1. Padi Sawah - Irigasi Teknis
		{
			ID:               utils.GenerateULID(),
			ExtensionOfficer: "Ir. Supardi, SP",
			VisitDate:        now.AddDate(0, 0, -5),
			FarmerName:       "Pak Sukirman",
			FarmerGroup:      "Poktan Tani Makmur",
			FarmerGroupType:  entity.FarmerGroupTypePoktan,
			Village:          "Sumbersari",
			District:         "Ngawi",
			Latitude:         -7.4040,
			Longitude:        111.4460,

			// Food Crops
			FoodCommodity:    entity.FoodCommodityPadiSawah,
			FoodLandStatus:   entity.LandStatusMilikSendiri,
			FoodLandArea:     2.5,
			FoodGrowthPhase:  entity.GrowthPhaseVegetatifAwal,
			FoodPlantAge:     45,
			FoodPlantingDate: timePtr(now.AddDate(0, 0, -45)),
			FoodHarvestDate:  timePtr(now.AddDate(0, 0, 75)),
			FoodTechnology:   entity.TechnologyMethodJajarLegowo,

			// Environmental
			WeatherCondition: entity.WeatherConditionCerah,
			WeatherImpact:    entity.WeatherImpactTidakAda,
			MainConstraint:   entity.MainConstraintIrigasiSulit,

			// Farmer Needs
			FarmerHope:     entity.FarmerHopeBantuanAlsintan,
			TrainingNeeded: entity.TrainingNeededPHT,
			UrgentNeeds:    entity.UrgentNeedsBibitPupukSegera,
			WaterAccess:    entity.WaterAccessMudah,
			Suggestions:    "Perlu bantuan pompa air untuk mengantisipasi musim kemarau",
			// CreatedBy:        firstUser.ID,
			CreatedAt: now.AddDate(0, 0, -5),
			UpdatedAt: now.AddDate(0, 0, -5),
		},

		// 2. Jagung - Lahan Kering
		{
			ID:               utils.GenerateULID(),
			ExtensionOfficer: "Dr. Ir. Bambang Hartono, MP",
			VisitDate:        now.AddDate(0, 0, -8),
			FarmerName:       "Pak Joko Widodo",
			FarmerGroup:      "Gapoktan Sumber Rezeki",
			FarmerGroupType:  entity.FarmerGroupTypeGapoktan,
			Village:          "Geneng",
			District:         "Ngawi",
			Latitude:         -7.3800,
			Longitude:        111.4200,

			// Food Crops
			FoodCommodity:    entity.FoodCommodityJagung,
			FoodLandStatus:   entity.LandStatusMilikSendiri,
			FoodLandArea:     1.8,
			FoodGrowthPhase:  entity.GrowthPhaseGeneratif2,
			FoodPlantAge:     70,
			FoodPlantingDate: timePtr(now.AddDate(0, 0, -70)),
			FoodHarvestDate:  timePtr(now.AddDate(0, 0, 20)),
			FoodTechnology:   entity.TechnologyMethodBibitUnggul,

			// Environmental
			WeatherCondition: entity.WeatherConditionHujan,
			WeatherImpact:    entity.WeatherImpactTidakAda,
			MainConstraint:   entity.MainConstraintModal,

			// Farmer Needs
			FarmerHope:     entity.FarmerHopeBantuanModal,
			TrainingNeeded: entity.TrainingNeededPascapanen,
			UrgentNeeds:    entity.UrgentNeedsModalDarurat,
			WaterAccess:    entity.WaterAccessTerbatas,
			Suggestions:    "Memerlukan bantuan modal dan akses pasar untuk hasil panen jagung",
			// CreatedBy:      firstUser.ID,
			CreatedAt: now.AddDate(0, 0, -8),
			UpdatedAt: now.AddDate(0, 0, -8),
		},

		// 3. Kedelai - Organik
		{
			ID:               utils.GenerateULID(),
			ExtensionOfficer: "Ir. Siti Nurhaliza, SP",
			VisitDate:        now.AddDate(0, 0, -3),
			FarmerName:       "Bu Wahyuni",
			FarmerGroup:      "Kelompok Tani Organik Lestari",
			FarmerGroupType:  entity.FarmerGroupTypePoktan,
			Village:          "Karangjati",
			District:         "Karangjati",
			Latitude:         -7.4200,
			Longitude:        111.3800,

			// Food Crops
			FoodCommodity:    entity.FoodCommodityKedelai,
			FoodLandStatus:   entity.LandStatusSewa,
			FoodLandArea:     1.2,
			FoodGrowthPhase:  entity.GrowthPhaseVegetatifAwal,
			FoodPlantAge:     30,
			FoodPlantingDate: timePtr(now.AddDate(0, 0, -30)),
			FoodHarvestDate:  timePtr(now.AddDate(0, 0, 60)),
			FoodTechnology:   entity.TechnologyMethodPupukOrganik,

			// Environmental
			WeatherCondition: entity.WeatherConditionMendung,
			WeatherImpact:    entity.WeatherImpactTidakAda,
			MainConstraint:   entity.MainConstraintHama,

			// Farmer Needs
			FarmerHope:     entity.FarmerHopePelatihan,
			TrainingNeeded: entity.TrainingNeededSertifikasi,
			UrgentNeeds:    entity.UrgentNeedsObatHama,
			WaterAccess:    entity.WaterAccessTerbatas,
			Suggestions:    "Pelatihan pengendalian hama organik dan sertifikasi produk organik",
			// CreatedBy:      firstUser.ID,
			CreatedAt: now.AddDate(0, 0, -3),
			UpdatedAt: now.AddDate(0, 0, -3),
		},

		// 4. Padi Ladang - Lahan Tadah Hujan
		{
			ID:               utils.GenerateULID(),
			ExtensionOfficer: "Ir. Ahmad Dahlan, SP, MP",
			VisitDate:        now.AddDate(0, 0, -10),
			FarmerName:       "Pak Sutrisno",
			FarmerGroup:      "Poktan Harapan Jaya",
			FarmerGroupType:  entity.FarmerGroupTypePoktan,
			Village:          "Paron",
			District:         "Paron",
			Latitude:         -7.5100,
			Longitude:        111.3400,

			// Food Crops
			FoodCommodity:    entity.FoodCommodityPadiLadang,
			FoodLandStatus:   entity.LandStatusBagiHasil,
			FoodLandArea:     0.8,
			FoodGrowthPhase:  entity.GrowthPhaseGeneratif1,
			FoodPlantAge:     85,
			FoodPlantingDate: timePtr(now.AddDate(0, 0, -85)),
			FoodHarvestDate:  timePtr(now.AddDate(0, 0, 25)),
			FoodTechnology:   entity.TechnologyMethodTidakAda,

			// Environmental
			WeatherCondition: entity.WeatherConditionHujan,
			WeatherImpact:    entity.WeatherImpactTanamanRebah,
			MainConstraint:   entity.MainConstraintIklim,

			// Farmer Needs
			FarmerHope:     entity.FarmerHopeAsuransi,
			TrainingNeeded: entity.TrainingNeededBudidayaModern,
			UrgentNeeds:    entity.UrgentNeedsBibitPupukSegera,
			WaterAccess:    entity.WaterAccessJauh,
			Suggestions:    "Asuransi pertanian dan sistem early warning cuaca ekstrem",
			// CreatedBy:      firstUser.ID,
			CreatedAt: now.AddDate(0, 0, -10),
			UpdatedAt: now.AddDate(0, 0, -10),
		},

		// 5. Sayuran (Cabai) - Hortikultura
		{
			ID:               utils.GenerateULID(),
			ExtensionOfficer: "Ir. Dewi Sartika, SP",
			VisitDate:        now.AddDate(0, 0, -2),
			FarmerName:       "Pak Hadi Suprapto",
			FarmerGroup:      "Kelompok Hortikultura Sejahtera",
			FarmerGroupType:  entity.FarmerGroupTypePoktan,
			Village:          "Kedunggalar",
			District:         "Kedunggalar",
			Latitude:         -7.4500,
			Longitude:        111.5200,

			// Horticulture Crops
			HortiCommodity:    entity.HortiCommoditySayuran,
			HortiSubCommodity: "Cabai Merah Keriting",
			HortiLandStatus:   entity.LandStatusMilikSendiri,
			HortiLandArea:     0.5,
			HortiGrowthPhase:  entity.HortiGrowthPhasePembuahan,
			HortiPlantAge:     60,
			HortiPlantingDate: timePtr(now.AddDate(0, 0, -60)),
			HortiHarvestDate:  timePtr(now.AddDate(0, 0, 30)),
			HortiTechnology:   entity.HortiTechnologyGreenhouse,

			// Environmental
			WeatherCondition: entity.WeatherConditionCerah,
			WeatherImpact:    entity.WeatherImpactTidakAda,
			MainConstraint:   entity.MainConstraintHama,

			// Farmer Needs
			FarmerHope:     entity.FarmerHopeAksesPasar,
			TrainingNeeded: entity.TrainingNeededPHT,
			UrgentNeeds:    entity.UrgentNeedsObatHama,
			WaterAccess:    entity.WaterAccessMudah,
			Suggestions:    "Akses ke pasar modern dan cold storage untuk menjaga kualitas cabai",
			// CreatedBy:      firstUser.ID,
			CreatedAt: now.AddDate(0, 0, -2),
			UpdatedAt: now.AddDate(0, 0, -2),
		},

		// 6. Tomat - Teknologi Hidroponik
		{
			ID:               utils.GenerateULID(),
			ExtensionOfficer: "Dr. Ir. Rina Kartika, MP",
			VisitDate:        now.AddDate(0, 0, -1),
			FarmerName:       "Bu Lastri Wulandari",
			FarmerGroup:      "Agrobisnis Modern Ngawi",
			FarmerGroupType:  entity.FarmerGroupTypeGapoktan,
			Village:          "Margomulyo",
			District:         "Ngawi",
			Latitude:         -7.4010,
			Longitude:        111.4380,

			// Horticulture Crops
			HortiCommodity:    entity.HortiCommoditySayuran,
			HortiSubCommodity: "Tomat",
			HortiLandStatus:   entity.LandStatusMilikSendiri,
			HortiLandArea:     0.3,
			HortiGrowthPhase:  entity.HortiGrowthPhaseVegetatif,
			HortiPlantAge:     35,
			HortiPlantingDate: timePtr(now.AddDate(0, 0, -35)),
			HortiHarvestDate:  timePtr(now.AddDate(0, 0, 55)),
			HortiTechnology:   entity.HortiTechnologyHydroponik,

			// Environmental
			WeatherCondition: entity.WeatherConditionCerah,
			WeatherImpact:    entity.WeatherImpactTidakAda,
			MainConstraint:   entity.MainConstraintTeknologi,

			// Farmer Needs
			FarmerHope:     entity.FarmerHopePelatihan,
			TrainingNeeded: entity.TrainingNeededGreenhouseIoT,
			UrgentNeeds:    entity.UrgentNeedsAlsintan,
			WaterAccess:    entity.WaterAccessMudah,
			Suggestions:    "Bantuan peralatan hidroponik dan pendampingan intensif teknologi modern",
			// CreatedBy:      firstUser.ID,
			CreatedAt: now.AddDate(0, 0, -1),
			UpdatedAt: now.AddDate(0, 0, -1),
		},

		// 7. Bawang Merah - Musim Panen
		{
			ID:               utils.GenerateULID(),
			ExtensionOfficer: "Ir. Budi Santoso, SP",
			VisitDate:        now.AddDate(0, 0, -4),
			FarmerName:       "Pak Sarimin",
			FarmerGroup:      "Poktan Bawang Ngawi",
			FarmerGroupType:  entity.FarmerGroupTypePoktan,
			Village:          "Mantingan",
			District:         "Mantingan",
			Latitude:         -7.4300,
			Longitude:        111.3600,

			// Horticulture Crops
			HortiCommodity:    entity.HortiCommoditySayuran,
			HortiSubCommodity: "Bawang Merah",
			HortiLandStatus:   entity.LandStatusMilikSendiri,
			HortiLandArea:     1.5,
			HortiGrowthPhase:  entity.HortiGrowthPhasePanen,
			HortiPlantAge:     65,
			HortiPlantingDate: timePtr(now.AddDate(0, 0, -65)),
			HortiHarvestDate:  timePtr(now.AddDate(0, 0, 5)),
			HortiTechnology:   entity.HortiTechnologyBibitUnggul,

			// Environmental
			WeatherCondition: entity.WeatherConditionCerah,
			WeatherImpact:    entity.WeatherImpactTidakAda,
			MainConstraint:   entity.MainConstraintAksesPasar,

			// Farmer Needs
			FarmerHope:     entity.FarmerHopeAksesPasar,
			TrainingNeeded: entity.TrainingNeededPascapanen,
			UrgentNeeds:    entity.UrgentNeedsPasarDarurat,
			WaterAccess:    entity.WaterAccessMudah,
			Suggestions:    "Perlu gudang penyimpanan dengan sirkulasi udara baik dan akses ke pedagang besar",
			// CreatedBy:      firstUser.ID,
			CreatedAt: now.AddDate(0, 0, -4),
			UpdatedAt: now.AddDate(0, 0, -4),
		},

		// 8. Tebu - Tanaman Perkebunan
		{
			ID:               utils.GenerateULID(),
			ExtensionOfficer: "Ir. Heru Prasetyo, MP",
			VisitDate:        now.AddDate(0, 0, -15),
			FarmerName:       "Pak Agus Salim",
			FarmerGroup:      "Gapoktan Tebu Manis",
			FarmerGroupType:  entity.FarmerGroupTypeGapoktan,
			Village:          "Jogorogo",
			District:         "Jogorogo",
			Latitude:         -7.5300,
			Longitude:        111.2800,

			// Plantation Crops
			PlantationCommodity:    entity.PlantationCommodityTebu,
			PlantationLandStatus:   entity.LandStatusMilikSendiri,
			PlantationLandArea:     3.5,
			PlantationGrowthPhase:  entity.PlantationGrowthPhaseTanamanMenghasilkan,
			PlantationPlantAge:     180,
			PlantationPlantingDate: timePtr(now.AddDate(0, 0, -180)),
			PlantationHarvestDate:  timePtr(now.AddDate(0, 0, 185)),
			PlantationTechnology:   entity.PlantationTechnologyTidakAda,

			// Environmental
			WeatherCondition: entity.WeatherConditionCerah,
			WeatherImpact:    entity.WeatherImpactTidakAda,
			MainConstraint:   entity.MainConstraintIrigasiSulit,

			// Farmer Needs
			FarmerHope:     entity.FarmerHopeBantuanAlsintan,
			TrainingNeeded: entity.TrainingNeededBudidayaModern,
			UrgentNeeds:    entity.UrgentNeedsBibitPupukSegera,
			WaterAccess:    entity.WaterAccessTerbatas,
			Suggestions:    "Diperlukan sistem irigasi tetes dan alat tebang tebu mekanis",
			// CreatedBy:      firstUser.ID,
			CreatedAt: now.AddDate(0, 0, -15),
			UpdatedAt: now.AddDate(0, 0, -15),
		},

		// 9. Kopi - Perkebunan
		{
			ID:               utils.GenerateULID(),
			ExtensionOfficer: "Dr. Ir. Fatimah Zahra, SP, MP",
			VisitDate:        now.AddDate(0, 0, -7),
			FarmerName:       "Pak Teguh Santoso",
			FarmerGroup:      "Kelompok Tani Kopi Nusantara",
			FarmerGroupType:  entity.FarmerGroupTypePoktan,
			Village:          "Sine",
			District:         "Sine",
			Latitude:         -7.4800,
			Longitude:        111.5800,

			// Plantation Crops
			PlantationCommodity:    entity.PlantationCommodityKopi,
			PlantationLandStatus:   entity.LandStatusMilikSendiri,
			PlantationLandArea:     2.0,
			PlantationGrowthPhase:  entity.PlantationGrowthPhaseTanamanMenghasilkan,
			PlantationPlantAge:     730, // 2 tahun
			PlantationPlantingDate: timePtr(now.AddDate(-2, 0, 0)),
			PlantationHarvestDate:  timePtr(now.AddDate(0, 2, 0)),
			PlantationTechnology:   entity.PlantationTechnologyPupukOrganik,

			// Environmental
			WeatherCondition: entity.WeatherConditionMendung,
			WeatherImpact:    entity.WeatherImpactTidakAda,
			MainConstraint:   entity.MainConstraintAksesPasar,

			// Farmer Needs
			FarmerHope:     entity.FarmerHopePelatihan,
			TrainingNeeded: entity.TrainingNeededPascapanen,
			UrgentNeeds:    entity.UrgentNeedsAlsintan,
			WaterAccess:    entity.WaterAccessJauh,
			Suggestions:    "Pelatihan roasting dan packaging untuk meningkatkan nilai jual kopi",
			// CreatedBy:      firstUser.ID,
			CreatedAt: now.AddDate(0, 0, -7),
			UpdatedAt: now.AddDate(0, 0, -7),
		},

		// 10. Ubi Kayu - Pangan Alternatif
		{
			ID:               utils.GenerateULID(),
			ExtensionOfficer: "Ir. Sulistyowati, SP",
			VisitDate:        now.AddDate(0, 0, -12),
			FarmerName:       "Pak Paijo",
			FarmerGroup:      "Poktan Sumber Pangan",
			FarmerGroupType:  entity.FarmerGroupTypePoktan,
			Village:          "Gerih",
			District:         "Gerih",
			Latitude:         -7.5500,
			Longitude:        111.3200,

			// Food Crops
			FoodCommodity:    entity.FoodCommodityUbiKayu,
			FoodLandStatus:   entity.LandStatusMilikSendiri,
			FoodLandArea:     1.0,
			FoodGrowthPhase:  entity.GrowthPhaseVegetatifAkhir,
			FoodPlantAge:     240, // 8 bulan
			FoodPlantingDate: timePtr(now.AddDate(0, -8, 0)),
			FoodHarvestDate:  timePtr(now.AddDate(0, 2, 0)),
			FoodTechnology:   entity.TechnologyMethodTidakAda,

			// Environmental
			WeatherCondition: entity.WeatherConditionCerah,
			WeatherImpact:    entity.WeatherImpactTidakAda,
			MainConstraint:   entity.MainConstraintAksesPasar,

			// Farmer Needs
			FarmerHope:     entity.FarmerHopeAksesPasar,
			TrainingNeeded: entity.TrainingNeededPascapanen,
			UrgentNeeds:    entity.UrgentNeedsAlsintan,
			WaterAccess:    entity.WaterAccessJauh,
			Suggestions:    "Akses ke pabrik tepung tapioka dan pelatihan diversifikasi olahan singkong",
			// CreatedBy:      firstUser.ID,
			CreatedAt: now.AddDate(0, 0, -12),
			UpdatedAt: now.AddDate(0, 0, -12),
		},

		// 11. Pisang - Buah-buahan
		{
			ID:               utils.GenerateULID(),
			ExtensionOfficer: "Ir. Nugroho Wijaya, SP",
			VisitDate:        now.AddDate(0, 0, -5),
			FarmerName:       "Pak Supardi",
			FarmerGroup:      "Kelompok Tani Buah Tropis",
			FarmerGroupType:  entity.FarmerGroupTypePoktan,
			Village:          "Padas",
			District:         "Padas",
			Latitude:         -7.6200,
			Longitude:        111.4700,

			// Horticulture Crops
			HortiCommodity:    entity.HortiCommodityBuah,
			HortiSubCommodity: "Pisang Raja",
			HortiLandStatus:   entity.LandStatusMilikSendiri,
			HortiLandArea:     0.7,
			HortiGrowthPhase:  entity.HortiGrowthPhasePembuahan,
			HortiPlantAge:     270, // 9 bulan
			HortiPlantingDate: timePtr(now.AddDate(0, -9, 0)),
			HortiHarvestDate:  timePtr(now.AddDate(0, 1, 0)),
			HortiTechnology:   entity.HortiTechnologyBibitUnggul,

			// Environmental
			WeatherCondition: entity.WeatherConditionMendung,
			WeatherImpact:    entity.WeatherImpactTidakAda,
			MainConstraint:   entity.MainConstraintHama,

			// Farmer Needs
			FarmerHope:     entity.FarmerHopePelatihan,
			TrainingNeeded: entity.TrainingNeededPHT,
			UrgentNeeds:    entity.UrgentNeedsObatHama,
			WaterAccess:    entity.WaterAccessTerbatas,
			Suggestions:    "Pengendalian hama penggerek batang dan akses pasar buah segar",
			// CreatedBy:      firstUser.ID,
			CreatedAt: now.AddDate(0, 0, -5),
			UpdatedAt: now.AddDate(0, 0, -5),
		},

		// 12. Padi Sawah - Serangan Hama
		{
			ID:               utils.GenerateULID(),
			ExtensionOfficer: "Ir. Tri Wahyuni, SP",
			VisitDate:        now.AddDate(0, 0, -2),
			FarmerName:       "Pak Kusno",
			FarmerGroup:      "Poktan Sido Makmur",
			FarmerGroupType:  entity.FarmerGroupTypePoktan,
			Village:          "Ketanggi",
			District:         "Kedunggalar",
			Latitude:         -7.4600,
			Longitude:        111.5100,

			// Food Crops
			FoodCommodity:    entity.FoodCommodityPadiSawah,
			FoodLandStatus:   entity.LandStatusMilikSendiri,
			FoodLandArea:     1.5,
			FoodGrowthPhase:  entity.GrowthPhaseGeneratif1,
			FoodPlantAge:     65,
			FoodPlantingDate: timePtr(now.AddDate(0, 0, -65)),
			FoodHarvestDate:  timePtr(now.AddDate(0, 0, 45)),
			FoodTechnology:   entity.TechnologyMethodPengendalianHama,

			// Pest Disease
			HasPestDisease:       true,
			PestDiseaseType:      entity.PestDiseaseWerengCoklat,
			PestDiseaseCommodity: "Padi Sawah",
			AffectedArea:         entity.AffectedAreaLevel25Sampai50,
			ControlAction:        entity.ControlActionPHT,

			// Environmental
			WeatherCondition: entity.WeatherConditionHujan,
			WeatherImpact:    entity.WeatherImpactDaunMenguning,
			MainConstraint:   entity.MainConstraintHama,

			// Farmer Needs
			FarmerHope:     entity.FarmerHopePelatihan,
			TrainingNeeded: entity.TrainingNeededPHT,
			UrgentNeeds:    entity.UrgentNeedsObatHama,
			WaterAccess:    entity.WaterAccessMudah,
			Suggestions:    "Serangan wereng coklat cukup parah, perlu pengendalian segera dan monitoring intensif",
			// CreatedBy:      firstUser.ID,
			CreatedAt: now.AddDate(0, 0, -2),
			UpdatedAt: now.AddDate(0, 0, -2),
		},

		// 13. Kacang Tanah - Program Diversifikasi
		{
			ID:               utils.GenerateULID(),
			ExtensionOfficer: "Ir. Yanto Suryanto, SP, MP",
			VisitDate:        now.AddDate(0, 0, -14),
			FarmerName:       "Bu Mariyem",
			FarmerGroup:      "Gapoktan Tani Maju",
			FarmerGroupType:  entity.FarmerGroupTypeGapoktan,
			Village:          "Ngale",
			District:         "Ngawi",
			Latitude:         -7.3950,
			Longitude:        111.4550,

			// Food Crops
			FoodCommodity:    entity.FoodCommodityKacangTanah,
			FoodLandStatus:   entity.LandStatusSewa,
			FoodLandArea:     0.8,
			FoodGrowthPhase:  entity.GrowthPhaseGeneratif3,
			FoodPlantAge:     80,
			FoodPlantingDate: timePtr(now.AddDate(0, 0, -80)),
			FoodHarvestDate:  timePtr(now.AddDate(0, 0, 10)),
			FoodTechnology:   entity.TechnologyMethodBibitUnggul,

			// Environmental
			WeatherCondition: entity.WeatherConditionCerah,
			WeatherImpact:    entity.WeatherImpactTidakAda,
			MainConstraint:   entity.MainConstraintModal,

			// Farmer Needs
			FarmerHope:     entity.FarmerHopeBantuanModal,
			TrainingNeeded: entity.TrainingNeededPascapanen,
			UrgentNeeds:    entity.UrgentNeedsAlsintan,
			WaterAccess:    entity.WaterAccessJauh,
			Suggestions:    "Memerlukan mesin pengering kacang dan akses ke industri pengolahan",
			// CreatedBy:      firstUser.ID,
			CreatedAt: now.AddDate(0, 0, -14),
			UpdatedAt: now.AddDate(0, 0, -14),
		},

		// 14. Kelapa - Tanaman Perkebunan
		{
			ID:               utils.GenerateULID(),
			ExtensionOfficer: "Ir. Retno Wulandari, SP",
			VisitDate:        now.AddDate(0, 0, -20),
			FarmerName:       "Pak Karno",
			FarmerGroup:      "Kelompok Tani Kelapa Sejahtera",
			FarmerGroupType:  entity.FarmerGroupTypePoktan,
			Village:          "Paron",
			District:         "Paron",
			Latitude:         -7.5150,
			Longitude:        111.3350,

			// Plantation Crops
			PlantationCommodity:    entity.PlantationCommodityKelapa,
			PlantationLandStatus:   entity.LandStatusMilikSendiri,
			PlantationLandArea:     1.5,
			PlantationGrowthPhase:  entity.PlantationGrowthPhaseTanamanMenghasilkan,
			PlantationPlantAge:     2190, // 6 tahun
			PlantationPlantingDate: timePtr(now.AddDate(-6, 0, 0)),
			PlantationHarvestDate:  timePtr(now.AddDate(0, 1, 0)),
			PlantationTechnology:   entity.PlantationTechnologyTidakAda,

			// Environmental
			WeatherCondition: entity.WeatherConditionCerah,
			WeatherImpact:    entity.WeatherImpactTidakAda,
			MainConstraint:   entity.MainConstraintAksesPasar,

			// Farmer Needs
			FarmerHope:     entity.FarmerHopeAksesPasar,
			TrainingNeeded: entity.TrainingNeededPascapanen,
			UrgentNeeds:    entity.UrgentNeedsAlsintan,
			WaterAccess:    entity.WaterAccessJauh,
			Suggestions:    "Akses ke pabrik kopra dan pelatihan diversifikasi olahan kelapa (VCO, serat)",
			// CreatedBy:      firstUser.ID,
			CreatedAt: now.AddDate(0, 0, -20),
			UpdatedAt: now.AddDate(0, 0, -20),
		},

		// 15. Ubi Jalar - Pangan Alternatif
		{
			ID:               utils.GenerateULID(),
			ExtensionOfficer: "Ir. Dwi Handayani, SP",
			VisitDate:        now.AddDate(0, 0, -6),
			FarmerName:       "Bu Suminah",
			FarmerGroup:      "Poktan Makmur Jaya",
			FarmerGroupType:  entity.FarmerGroupTypePoktan,
			Village:          "Ngrambe",
			District:         "Ngrambe",
			Latitude:         -7.3800,
			Longitude:        111.3400,

			// Food Crops
			FoodCommodity:    entity.FoodCommodityUbiJalar,
			FoodLandStatus:   entity.LandStatusMilikSendiri,
			FoodLandArea:     0.6,
			FoodGrowthPhase:  entity.GrowthPhaseGeneratif1,
			FoodPlantAge:     90,
			FoodPlantingDate: timePtr(now.AddDate(0, 0, -90)),
			FoodHarvestDate:  timePtr(now.AddDate(0, 0, 30)),
			FoodTechnology:   entity.TechnologyMethodPupukOrganik,

			// Environmental
			WeatherCondition: entity.WeatherConditionCerah,
			WeatherImpact:    entity.WeatherImpactTidakAda,
			MainConstraint:   entity.MainConstraintAksesPasar,

			// Farmer Needs
			FarmerHope:     entity.FarmerHopeAksesPasar,
			TrainingNeeded: entity.TrainingNeededPascapanen,
			UrgentNeeds:    entity.UrgentNeedsPasarDarurat,
			WaterAccess:    entity.WaterAccessTerbatas,
			Suggestions:    "Pengolahan ubi jalar menjadi produk olahan bernilai tambah tinggi",
			// CreatedBy:      firstUser.ID,
			CreatedAt: now.AddDate(0, 0, -6),
			UpdatedAt: now.AddDate(0, 0, -6),
		},
	}

	for _, report := range reports {
		if err := db.Create(&report).Error; err != nil {
			return err
		}
	}

	log.Printf("Successfully seeded %d agriculture reports", len(reports))
	return nil
}
