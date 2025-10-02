package seeds

import (
	"gorm.io/gorm"
	
)

func SeedAll(db *gorm.DB) error {
	seeders := []struct {
		name string
		fn   func(*gorm.DB) error
	}{
		{"Users", SeedUsers},
		{"Reports", SeedReports},
		{"Spatial Planning", SeedSpatialPlanning},
		{"Water Resources", SeedWaterResources},
		{"Bina Marga", SeedBinaMarga},
		{"Agriculture", SeedAgriculture},
	}

	for _, seeder := range seeders {
		if err := seeder.fn(db); err != nil {
			return err
		}
	}

	return nil
}