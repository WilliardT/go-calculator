package main

import (
	"log"

	calculationservice "GO-Calc/internal/calculationService"
	"GO-Calc/internal/db"
	"GO-Calc/internal/handlers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)


func main() {
	database, err := db.InitDB()

	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	e := echo.New()

	calcRepo := calculationservice.NewCalculationRepository(database)
	calcService := calculationservice.NewCalculationService(calcRepo)
	calcHandlers := handlers.NewCalculationHandler(calcService)

	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	
	e.GET("/calculations", calcHandlers.GetCalculations)

	e.POST("/calculations", calcHandlers.PostCalculation)

	e.PATCH("/calculations/:id", calcHandlers.PatchCalculation)

	e.DELETE("/calculations/:id", calcHandlers.DeleteCalculation)
	
	e.Start("localhost:8080")
}