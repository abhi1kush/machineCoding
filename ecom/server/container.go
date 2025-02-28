package server

import (
	"database/sql"

	"ecom.com/cache"
	"ecom.com/config"
	"ecom.com/database"
	"ecom.com/handlers"
	"ecom.com/repository"
	"ecom.com/routes"
	"ecom.com/services"
)

// Container holds all dependencies
type Container struct {
	Cache    cache.CacheI
	DB       *sql.DB
	MetricDB *sql.DB

	UserRepo   repository.UserRepository
	OrderRepo  repository.OrderRepositoryI
	MetricRepo repository.MetricRepositoryI

	UserService   *services.UserService
	OrderService  *services.Order
	MetricService *services.Metric

	UserHandler   *handlers.UserHandler
	MetricHandler *handlers.MetricHandler
	OrderHandler  *handlers.OrderHandler

	RoutesCfg *routes.RouterConfig
}

// NewContainer initializes dependencies
func NewContainer(appConfig config.Config) *Container {
	//setup cache
	cache := cache.NewRedis(appConfig.Redis.Addr, appConfig.Redis.Password, 1)

	// Initialize database
	db := database.ConnectDB(appConfig.Database.Driver, appConfig.Database.DSN)
	metricDb := database.ConnectMetricsDB(appConfig.Metrics.Driver, appConfig.Metrics.DSN)

	// Initialize repository
	userRepo := repository.NewSQLiteUserRepository(db)
	orderRepo := repository.NewSQLiteOrderRepository(db)
	metricRepo := repository.NewSQLiteMetricRepository(metricDb)

	// Initialize service
	userService := services.NewUserService(userRepo)
	orderService := services.NewOrderService(appConfig, orderRepo, metricRepo, cache)
	metricService := services.NewMetricService(metricRepo)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)
	orderHandler := handlers.NewOrderHandler(orderService)
	metricHandler := handlers.NewMetricHandler(metricService)

	return &Container{
		Cache: cache,

		DB:       db,
		MetricDB: metricDb,

		UserRepo:   userRepo,
		OrderRepo:  orderRepo,
		MetricRepo: metricRepo,

		UserService:  userService,
		OrderService: orderService,

		UserHandler:   userHandler,
		OrderHandler:  orderHandler,
		MetricHandler: metricHandler,

		RoutesCfg: &routes.RouterConfig{
			UserHandler:   userHandler,
			OrderHandler:  orderHandler,
			MetricHandler: metricHandler,
		},
	}
}
