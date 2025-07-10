package main

import (
	"fmt"
	"math"
	"net/http"

	"github.com/Knetic/govaluate"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Calculation struct {
	ID 					string `json:"id"`
	Expression 	string `json:"expression"`
	Result 			string `json:"result"`
}

type CalculationRequest struct {
	Expression 	string `json:"expression"`
}

// временно хранение в памяти slice
var calculations = []Calculation{}

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

	calculations = append(calculations, calc)

	return c.JSON(http.StatusCreated, calc)
}

func main() {
	fmt.Println("Hello, world")
	fmt.Println(math.Pi)

	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	
	e.GET("/calculations", getCalculations)
	
	e.POST("/calculations", postCalculation)
	
	e.Start("localhost:8080")
}