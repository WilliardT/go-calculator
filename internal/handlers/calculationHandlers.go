package handlers

import (
	calculationservice "GO-Calc/internal/calculationService"
	"net/http"

	"github.com/labstack/echo/v4"
)

// принимвет

type CalculationHandler struct {
	service calculationservice.CalculationService
}

func NewCalculationHandler(s calculationservice.CalculationService) *CalculationHandler {
	return &CalculationHandler{service: s}
}

// история расчетов
func (h *CalculationHandler) GetCalculations(c echo.Context) error {
	calculations, err := h.service.GetAllCalculations()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not get calculations"})
	}

	return c.JSON(http.StatusOK, calculations)
}

func (h *CalculationHandler) PostCalculation (c echo.Context) error {
	var req calculationservice.CalculationRequest

	if err := c.Bind(&req); err != nil {
		// проблема с данными, не можем декодировать
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	calc, err := h.service.CreateCalculate(req.Expression)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Could not create calculation"})
	}

	return c.JSON(http.StatusCreated, calc)
}

func (h *CalculationHandler) PatchCalculation(c echo.Context) error {
	id := c.Param("id")

	var req calculationservice.CalculationRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	updatedCalc, err := h.service.UpdateCalculation(id, req.Expression)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Could not update calculation"})
	}

	return c.JSON(http.StatusOK, updatedCalc)
}

func (h *CalculationHandler) DeleteCalculation (c echo.Context) error {
	id := c.Param("id")

	if err := h.service.DeleteCalculation(id); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Could not delete calculation"})
	}

	return c.NoContent(http.StatusNoContent)
}
