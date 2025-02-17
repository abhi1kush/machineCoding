package main

import (
	"log"

	"ecom.com/cache"
	"ecom.com/config"
	"ecom.com/db"
	"ecom.com/logger"
	"ecom.com/queue"
	"ecom.com/routes"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	logger.InitLogger("app.log", 10, 5, 30, true)
	logger.Logger.Println("Logger initialized")
	config.LoadConfig("config/config.yaml")

	db.InitDB(config.AppConfig.Database.Driver, config.AppConfig.Database.DSN)
	defer db.CloseDB()

	db.InitMetricsDB(config.AppConfig.Metrics.Driver, config.AppConfig.Metrics.DSN)
	defer db.CloseMetricsDB()

	cache.InitRedis(config.AppConfig.Redis.Addr, config.AppConfig.Redis.Password, config.AppConfig.Redis.DB)

	queue.StartOrderProcessor(config.AppConfig.Queue.WorkerPool, config.AppConfig.Queue.QueueCapacity)

	r := gin.Default()
	routes.SetupRoutes(r)

	queue.StartOrderProcessor(config.AppConfig.Queue.WorkerPool, config.AppConfig.Queue.QueueCapacity)

	log.Println("Server running on port", config.AppConfig.Server.Port)
	r.Run(":" + config.AppConfig.Server.Port)
}
