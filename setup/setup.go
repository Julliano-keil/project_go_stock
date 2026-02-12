package setup

import (
	"context"
	"fmt"
	"log"
	"time"

	"lince/datastore"
	"lince/entities"
	"lince/migrations"
	"lince/modules/category"
	"lince/modules/equipment"
	"lince/modules/stock_movement"
	"lince/modules/stock_unit"
	"lince/modules/subcategory"

	"github.com/robfig/cron/v3"
)

func SetupModules(cfg entities.Config) func() error {
	log.Println("Start modules setup")

	database, err := datastore.NewSettingsRepository(cfg)
	if err != nil {
		log.Fatalf("database connection: %v", err)
	}

	db := database.Connection(entities.CompanyDatabaseConfig{})
	if err := migrations.Run(db); err != nil {
		log.Fatalf("migrations: %v", err)
	}

	// ### REPOSITORIES ###

	categoryRepository := category.NewCategoryRepository(database)
	subcategoryRepository := subcategory.NewSubcategoryRepository(database)
	equipmentRepository := equipment.NewEquipmentRepository(database)
	stockUnitRepository := stock_unit.NewStockUnitRepository(database)
	stockMovementRepository := stock_movement.NewStockMovementRepository(database)

	// ### USE CASES ###
	categoryUseCase := category.NewCategoryUseCase(categoryRepository, cfg)
	subcategoryUseCase := subcategory.NewSubcategoryUseCase(subcategoryRepository, cfg)
	equipmentUseCase := equipment.NewEquipmentUseCase(equipmentRepository, cfg)
	stockUnitUseCase := stock_unit.NewStockUnitUseCase(stockUnitRepository, cfg)
	stockMovementUseCase := stock_movement.NewStockMovementUseCase(stockMovementRepository, cfg)

	_ = categoryUseCase
	_ = subcategoryUseCase
	_ = equipmentUseCase
	_ = stockUnitUseCase
	_ = stockMovementUseCase

	cronHandler := cron.New(cron.WithLocation(time.Local))

	_, _ = cronHandler.AddFunc("@midnight", func() {
		log.Println("CRON @midnight =================")
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"))

		ctx := context.Background()
		company := entities.CompanyDatabaseConfig{}

		_ = ctx
		_ = company
	})

	return func() error {
		cronHandler.Start()

		for {
			ctx := context.Background()
			company := entities.CompanyDatabaseConfig{}

			log.Println("SYNC =================")

			_ = ctx
			_ = company

			log.Println("Finished successfully SYNC")

			log.Println("Sleep for 10 minutes")
			time.Sleep(10 * time.Minute)
		}
	}
}
