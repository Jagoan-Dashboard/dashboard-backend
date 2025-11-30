package main

import (
	"building-report-backend/pkg/config"
	"building-report-backend/pkg/database"
	"building-report-backend/seeds"
	"fmt"
	"log"
	"os"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Connect to database
	db, err := database.NewPostgresDB(cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Get command line arguments
	args := os.Args[1:]
	if len(args) == 0 {
		log.Println("Usage: go run cmd/seeder/main.go [table_name|all]")
		log.Println("Available tables: users, reports, spatial_planning, water_resources, bina_marga, agriculture, all")
		return
	}

	command := args[0]

	switch command {
	case "users":
		if err := seeds.SeedUsers(db); err != nil {
			log.Fatal("Failed to seed users:", err)
		}
		fmt.Println("✓ Users seeded successfully")

	case "reports":
		if err := seeds.SeedReports(db); err != nil {
			log.Fatal("Failed to seed reports:", err)
		}
		fmt.Println("✓ Reports seeded successfully")

	case "spatial_planning":
		if err := seeds.SeedSpatialPlanning(db); err != nil {
			log.Fatal("Failed to seed spatial planning:", err)
		}
		fmt.Println("✓ Spatial Planning seeded successfully")

	case "water_resources":
		if err := seeds.SeedWaterResources(db); err != nil {
			log.Fatal("Failed to seed water resources:", err)
		}
		fmt.Println("✓ Water Resources seeded successfully")

	case "bina_marga":
		if err := seeds.SeedBinaMarga(db); err != nil {
			log.Fatal("Failed to seed bina marga:", err)
		}
		fmt.Println("✓ Bina Marga seeded successfully")

	case "agriculture":
		if err := seeds.SeedAgriculture(db); err != nil {
			log.Fatal("Failed to seed agriculture:", err)
		}
		fmt.Println("✓ Agriculture seeded successfully")
	case "rice_fields":
		if err := seeds.SeedRiceFields(db); err != nil {
			log.Fatal("Failed to seed rice fields:", err)
		}
		fmt.Println("✓ Rice Fields seeded successfully")

	case "all":
		log.Println("Seeding all tables...")

		if err := seeds.SeedUsers(db); err != nil {
			log.Fatal("Failed to seed users:", err)
		}
		fmt.Println("✓ Users seeded successfully")

		if err := seeds.SeedReports(db); err != nil {
			log.Fatal("Failed to seed reports:", err)
		}
		fmt.Println("✓ Reports seeded successfully")

		if err := seeds.SeedSpatialPlanning(db); err != nil {
			log.Fatal("Failed to seed spatial planning:", err)
		}
		fmt.Println("✓ Spatial Planning seeded successfully")

		if err := seeds.SeedWaterResources(db); err != nil {
			log.Fatal("Failed to seed water resources:", err)
		}
		fmt.Println("✓ Water Resources seeded successfully")

		if err := seeds.SeedBinaMarga(db); err != nil {
			log.Fatal("Failed to seed bina marga:", err)
		}
		fmt.Println("✓ Bina Marga seeded successfully")

		if err := seeds.SeedAgriculture(db); err != nil {
			log.Fatal("Failed to seed agriculture:", err)
		}
		fmt.Println("✓ Agriculture seeded successfully")

		if err := seeds.SeedRiceFields(db); err != nil {
			log.Fatal("Failed to seed rice fields:", err)
		}
		fmt.Println("✓ Rice Fields seeded successfully")

		fmt.Println("\n✓ All tables seeded successfully!")

	default:
		log.Printf("Unknown command: %s\n", command)
		log.Println("Available commands: users, reports, spatial_planning, water_resources, bina_marga, agriculture, all")
	}
}
