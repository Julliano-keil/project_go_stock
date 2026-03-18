package setup

import (
	"log"

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

	_ = stockMovementUseCase

	// ### ENDPOINT MODULES ###

	categoryModule := category.NewCategoryModule(categoryUseCase)
	subcategoryModule := subcategory.NewSubcategoryModule(subcategoryUseCase)
	equipmentModule := equipment.NewEquipmentModule(equipmentUseCase)
	stockUnitModule := stock_unit.NewStockUnitModule(stockUnitUseCase)

	// ### Authentication Module ###

	authenticationModule := user.NewAuthenticationModule(cfg, userUseCase)

	appModules := []modules.AppModule{
		categoryModule,
		subcategoryModule,
		equipmentModule,
		stockUnitModule,
	}

	
	apiRouter := r.PathPrefix("/api").Subrouter()
	routerBase := authenticationModule.Setup(apiRouter)
	for _, am := range appModules {
		moduleSubRouter := routerBase.PathPrefix(am.Path()).Subrouter()
		_ = am.Setup(moduleSubRouter)
	}

}
