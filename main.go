package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Knetic/govaluate"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func initDB() {
	// TODO вынести в env
	// data source name
	dsn := "host=localhost user=postgres password=yourpassword dbname=postgres port=5432 sslmode=disable" // sslmode безопасное соединение

	var err error

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{}) // если менять в настройках базы данных

	if err != nil {
		log.Fatalf("Could not connect to databse: %v",  err)
	}

	// нет SQL . автомиграция
	if err := db.AutoMigrate(&Calculation{}); err != nil {
		log.Fatalf("Could not migrate database: %v", err)
	}
}

type Calculation struct {
	ID 					string `gorm:"primaryKey" json:"id"` // тег на первичный ключ
	Expression 	string `json:"expression"`
	Result 			string `json:"result"`
}

type CalculationRequest struct {
	Expression 	string `json:"expression"`
}


func calculateExpression(expression string) (string, error) {
	expr, err := govaluate.NewEvaluableExpression(expression)
	
	if err != nil {
		return "", err
	}

	result, err := expr.Evaluate(nil)

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%v", result), nil
}

// история расчетов
func getCalculations(c echo.Context) error {
	var calculations []Calculation

	// Find находит все записи, соответствующие заданным условиям conds
	if err := db.Find(&calculations).Error; err !=nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not get calculations"})
	}

	return c.JSON(http.StatusOK, calculations)
}

func postCalculation (c echo.Context) error {
	var req CalculationRequest

	if err := c.Bind(&req); err != nil {
		// проблема с данными, не можем декодировать
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	result, err := calculateExpression(req.Expression)

	if err != nil {
		//проблема с самим выражением
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid expression"})
	}

	calc := Calculation{
		ID: 				uuid.NewString(),
		Expression: req.Expression,
		Result: 		result,
	}

	if err := db.Create(&calc).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not create calculation"})
	}

	return c.JSON(http.StatusCreated, calc)
}

func patchCalculation(c echo.Context) error {
	id := c.Param("id")

	var req CalculationRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	result, err := calculateExpression(req.Expression)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid expression"})
	}

	var calc Calculation

	// Сначала находит первую запись, упорядоченную по первичному ключу, соответствующую заданным условиям.
	if err := db.First(&calc, "id = ?", id).Error; err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Count not found expression"})
	}

	calc.Expression = req.Expression
	calc.Result = result

	if err := db.Save(&calc).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not update calculation"})
	}

	// for i, calculation := range calculations {
	// 	if calculation.ID == id {
	// 		calculations[i].Expression = req.Expression
	// 		calculations[i].Result = result

	// 		return c.JSON(http.StatusOK, calculations[i])
	// 	}
	// }

	return c.JSON(http.StatusOK, calc)
}

func deleteCalculation (c echo.Context) error {
	id := c.Param("id")

	if err := db.Delete(&Calculation{}, "id = ?", id).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not delete calculation"})
	}

	// временно, пока данные в храняться в slice calculations
	// for i, calculation := range calculations {
	// 	if calculation.ID == id {

	// 		calculations = append(calculations[:i], calculations[i+1:]...)
			
	// 		return c.JSON(http.StatusOK, map[string]string{"message": "Calculation deleted"})
	// 		return c.NoContent(http.StatusNoContent)
	// 	}
	//	return c.JSON(http.StatusBadRequest, map[string]string{"error": "Calculation not found"})
	// }

	return c.NoContent(http.StatusNoContent)
}

func main() {
	initDB()

	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	
	e.GET("/calculations", getCalculations)
	
	e.POST("/calculations", postCalculation)

	e.PATCH("/calculations/:id", patchCalculation)

	e.DELETE("/calculations/:id", deleteCalculation)
	
	e.Start("localhost:8080")
}