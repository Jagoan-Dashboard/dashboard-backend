package seeds

import (
	"building-report-backend/internal/domain/entity"
	"building-report-backend/pkg/utils"
	"gorm.io/gorm"
	"time"
)

func SeedReports(db *gorm.DB) error {
	// Get first user as CreatedBy
	var firstUser entity.User
	if err := db.First(&firstUser).Error; err != nil {
		return err
	}

	workTypeRehab := entity.WorkTypePerbaikanAtap
	workTypeBaru := entity.WorkTypePerbaikanDinding
	conditionBaik := entity.ConditionBaik
	conditionPerbaikan := entity.ConditionButuhPerbaikan

	reports := []entity.Report{
		{
			ID:                   utils.GenerateULID(),
			ReporterName:         "Budi Santoso",
			ReporterRole:         entity.RolePerangkatDesa,
			Village:              "Karangjati",
			District:             "Ngawi",
			BuildingName:         "SDN 1 Karangjati",
			BuildingType:         entity.BuildingSekolah,
			ReportStatus:         entity.StatusRehabilitasi,
			FundingSource:        entity.FundingAPBDKab,
			LastYearConstruction: 2015,
			FullAddress:          "Jl. Raya Karangjati No. 45, Ngawi",
			Latitude:             -7.4040,
			Longitude:            111.4464,
			FloorArea:            450.5,
			FloorCount:           2,
			WorkType:             &workTypeRehab,
			ConditionAfterRehab:  &conditionBaik,
			CreatedBy:            firstUser.ID,
			CreatedAt:            time.Now().AddDate(0, 0, -10),
			UpdatedAt:            time.Now().AddDate(0, 0, -10),
		},
		{
			ID:                   utils.GenerateULID(),
			ReporterName:         "Siti Aminah",
			ReporterRole:         entity.RoleOPD,
			Village:              "Paron",
			District:             "Ngawi",
			BuildingName:         "Puskesmas Paron",
			BuildingType:         entity.BuildingPuskesmas,
			ReportStatus:         entity.StatusRehabilitasi,
			FundingSource:        entity.FundingAPBDProv,
			LastYearConstruction: 2018,
			FullAddress:          "Jl. Raya Paron No. 123, Ngawi",
			Latitude:             -7.3950,
			Longitude:            111.4350,
			FloorArea:            650.0,
			FloorCount:           2,
			WorkType:             &workTypeRehab,
			ConditionAfterRehab:  &conditionBaik,
			CreatedBy:            firstUser.ID,
			CreatedAt:            time.Now().AddDate(0, 0, -8),
			UpdatedAt:            time.Now().AddDate(0, 0, -8),
		},
		{
			ID:                   utils.GenerateULID(),
			ReporterName:         "Ahmad Fauzi",
			ReporterRole:         entity.RolePerangkatDesa,
			Village:              "Geneng",
			District:             "Ngawi",
			BuildingName:         "Pasar Geneng",
			BuildingType:         entity.BuildingPasar,
			ReportStatus:         entity.StatusRehabilitasi,
			FundingSource:        entity.FundingAPBDKab,
			LastYearConstruction: 2010,
			FullAddress:          "Jl. Pasar Geneng No. 78, Ngawi",
			Latitude:             -7.3800,
			Longitude:            111.4200,
			FloorArea:            850.0,
			FloorCount:           1,
			WorkType:             &workTypeRehab,
			ConditionAfterRehab:  &conditionPerbaikan,
			CreatedBy:            firstUser.ID,
			CreatedAt:            time.Now().AddDate(0, 0, -15),
			UpdatedAt:            time.Now().AddDate(0, 0, -15),
		},
		{
			ID:                   utils.GenerateULID(),
			ReporterName:         "Dewi Lestari",
			ReporterRole:         entity.RoleOPD,
			Village:              "Mantingan",
			District:             "Ngawi",
			BuildingName:         "Balai Desa Mantingan",
			BuildingType:         entity.BuildingKantorPemerintah,
			ReportStatus:         entity.StatusPembangunanBaru,
			FundingSource:        entity.FundingDanaDesa,
			LastYearConstruction: 2020,
			FullAddress:          "Jl. Raya Mantingan No. 12, Ngawi",
			Latitude:             -7.4100,
			Longitude:            111.4600,
			FloorArea:            280.0,
			FloorCount:           1,
			WorkType:             &workTypeBaru,
			CreatedBy:            firstUser.ID,
			CreatedAt:            time.Now().AddDate(0, 0, -5),
			UpdatedAt:            time.Now().AddDate(0, 0, -5),
		},
		{
			ID:                   utils.GenerateULID(),
			ReporterName:         "Eko Prasetyo",
			ReporterRole:         entity.RoleKelompokMasyarakat,
			Village:              "Ketanggi",
			District:             "Ngawi",
			BuildingName:         "Lapangan Olahraga Ketanggi",
			BuildingType:         entity.BuildingSaranaOlahraga,
			ReportStatus:         entity.StatusRehabilitasi,
			FundingSource:        entity.FundingSwadaya,
			LastYearConstruction: 2016,
			FullAddress:          "Desa Ketanggi, Kec. Ngawi",
			Latitude:             -7.4200,
			Longitude:            111.4500,
			FloorArea:            1200.0,
			FloorCount:           1,
			WorkType:             &workTypeRehab,
			ConditionAfterRehab:  &conditionBaik,
			CreatedBy:            firstUser.ID,
			CreatedAt:            time.Now().AddDate(0, 0, -3),
			UpdatedAt:            time.Now().AddDate(0, 0, -3),
		},
	}

	for _, report := range reports {
		if err := db.Create(&report).Error; err != nil {
			return err
		}
	}

	return nil
}

// func stringPtr(s string) *string {
// 	return &s
// }