package calculationservice

import (
	"fmt"

	"github.com/Knetic/govaluate"
	"github.com/google/uuid"
)

// бизнес логика
// получаем выражение иp ручек, считаем и передаем в repository

type CalculationService interface {
	CreateCalculate(expression string) (Calculation, error)
	GetAllCalculations() ([]Calculation, error)
	GetCalculationByID(id string) (Calculation, error)
	UpdateCalculation(id string, expression string) (Calculation, error)
	DeleteCalculation(id string) error
}

type calcService struct {
	repo CalculationRepository
}


func NewCalculationService(repo CalculationRepository) CalculationService {
	return &calcService{repo: repo}
}

func (s *calcService) calculateExpression(expression string) (string, error) {
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

func (s *calcService) CreateCalculate(expression string) (Calculation, error) {
	result, err := s.calculateExpression(expression)

	if err != nil {
		return Calculation{}, err
	}

	calc := Calculation{
		ID:					uuid.NewString(),
		Expression: expression,
		Result:		  result,
	}

	if err := s.repo.CreateCalculation(calc); err != nil {
		return Calculation{}, err
	}

	return calc, nil
}

func (s *calcService) GetAllCalculations() ([]Calculation, error) {
	return s.repo.GetAllCalculations()
}

func (s *calcService) GetCalculationByID(id string) (Calculation, error) {
	// возможно валидация на id (слишком короткий или другое...)
	return s.repo.GetCalculationByID(id)
}

func (s *calcService) UpdateCalculation(id string, expression string) (Calculation, error) {
	calc, err := s.GetCalculationByID(id)

	if err != nil {
		return Calculation{}, err
	}

	result, err := s.calculateExpression(expression)

	if err != nil {
		return Calculation{}, err
	}

	calc.Expression = expression
	calc.Result = result

	if err := s.repo.UpdateCalculation(calc); err != nil {
		return Calculation{}, err
	}

	return calc, nil
}

func (s *calcService) DeleteCalculation(id string) error {
	return s.repo.DeleteCalculation(id)
}


