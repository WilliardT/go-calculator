package calculationservice

type Calculation struct {
	ID 					string `gorm:"primaryKey" json:"id"` // тег на первичный ключ
	Expression 	string `json:"expression"`
	Result 			string `json:"result"`
}

type CalculationRequest struct {
	Expression 	string `json:"expression"`
}
