package setup

import (
	"context"
	"fmt"
	"log"
	"time"

	"lince/datastore"
	"lince/entities"
	"lince/migrations"
	"lince/modules"
	"lince/modules/category"
	"lince/modules/equipment"
	"lince/modules/stock_movement"
	"lince/modules/stock_unit"
	"lince/modules/subcategory"
	"lince/modules/user"

	"github.com/gorilla/mux"
	"github.com/robfig/cron/v3"
)

func SetupModules(r *mux.Router, cfg entities.Config) {
	log.Println("Setup modules")

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
	userRepository := user.NewUserRepository(database)
	equipmentRepository := equipment.NewEquipmentRepository(database)
	stockUnitRepository := stock_unit.NewStockUnitRepository(database)
	stockMovementRepository := stock_movement.NewStockMovementRepository(database)

	// ### USE CASES ###

	categoryUseCase := category.NewCategoryUseCase(categoryRepository, cfg)
	subcategoryUseCase := subcategory.NewSubcategoryUseCase(subcategoryRepository, cfg)
	userUseCase := user.NewUserUseCase(userRepository, cfg)
	equipmentUseCase := equipment.NewEquipmentUseCase(equipmentRepository, cfg)
	stockUnitUseCase := stock_unit.NewStockUnitUseCase(stockUnitRepository, cfg)
	stockMovementUseCase := stock_movement.NewStockMovementUseCase(stockMovementRepository, cfg)

	_ = equipmentUseCase
	_ = stockUnitUseCase
	_ = stockMovementUseCase

	// ### ENDPOINT MODULES ###

	categoryModule := category.NewCategoryModule(categoryUseCase)
	subcategoryModule := subcategory.NewSubcategoryModule(subcategoryUseCase)

	// ### Authentication Module ###

	authenticationModule := user.NewAuthenticationModule(cfg, userUseCase)

	appModules := []modules.AppModule{
		categoryModule,
		subcategoryModule,
	}

	routerBase := authenticationModule.Setup(r)
	for _, am := range appModules {
		moduleSubRouter := routerBase.PathPrefix(am.Path()).Subrouter()
		_ = am.Setup(moduleSubRouter)
	}

	// Cron em background
	cronHandler := cron.New(cron.WithLocation(time.Local))
	_, _ = cronHandler.AddFunc("@midnight", func() {
		log.Println("CRON @midnight =================")
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
		_ = context.Background()
		_ = entities.CompanyDatabaseConfig{}
	})
	cronHandler.Start()

	go func() {
		for {
			log.Println("SYNC =================")
			time.Sleep(10 * time.Minute)
		}
	}()
}
