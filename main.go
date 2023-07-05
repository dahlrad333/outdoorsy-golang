package main

import (
	"interview-challenge-backend/database"
	"interview-challenge-backend/handlers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

func main() {
	// Connect to the database
	database.ConnectDB()
	defer database.DB.Close()

	// Initialize Echo
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Init routes
	initPublicRoutes(e)

	// Start the server
	e.Logger.Fatal(e.Start(":8000"))
}

func initPublicRoutes(e *echo.Echo) {
	// Routes
	e.GET("/rentals/:id", handlers.GetRental)
	e.GET("/rentals", handlers.ListRentals)
}








