package main

import (
	"log"

	"ecom.com/config"
	"ecom.com/database"
	"ecom.com/logger"
	"ecom.com/routes"
	"ecom.com/server"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

/*
Initialize the database connection.
Create the repository layer.
Initialize the service layer with the repository.
Initialize the handler with the service.
Pass the handler to the router for route registration.
Start the Gin server.
*/

func main() {
	logger.InitLogger("app.log", 10, 5, 30, true)
	logger.Logger.Println("Logger initialized")
	config.LoadConfig("config/config.yaml")

	container := server.NewContainer(config.AppConfig)
	defer database.CloseDB(container.DB)
	defer database.CloseDB(container.MetricDB)

	container.OrderService.GetOrderCreationQueue().StartOrderProcessor()
	container.OrderService.GetOrderProcessQueue().StartOrderProcessor()

	r := gin.Default()
	routes.RegisterRoutes(r, container.RoutesCfg)

	log.Println("Server running on port", config.AppConfig.Server.Port)
	r.Run(":" + config.AppConfig.Server.Port)
}
