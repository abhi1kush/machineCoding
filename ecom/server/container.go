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

	OrderRepo  repository.OrderRepositoryI
	MetricRepo repository.MetricRepositoryI

	OrderService  *services.Order
	MetricService *services.Metric

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
	orderRepo := repository.NewSQLiteOrderRepository(db)
	itemRepo := repository.NewSQLiteItemRepository(db)
	metricRepo := repository.NewSQLiteMetricRepository(metricDb)

	// Initialize service
	orderService := services.NewOrderService(appConfig, orderRepo, itemRepo, metricRepo, cache)
	metricService := services.NewMetricService(metricRepo)

	// Initialize handlers
	orderHandler := handlers.NewOrderHandler(orderService)
	metricHandler := handlers.NewMetricHandler(metricService)

	return &Container{
		Cache: cache,

		DB:       db,
		MetricDB: metricDb,

		OrderRepo:  orderRepo,
		MetricRepo: metricRepo,

		OrderService: orderService,

		OrderHandler:  orderHandler,
		MetricHandler: metricHandler,

		RoutesCfg: &routes.RouterConfig{
			OrderHandler:  orderHandler,
			MetricHandler: metricHandler,
		},
	}
}
