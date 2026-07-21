package model

type LogEventRequest struct {
	EventText string `json:"event_text"`
}

type ConfirmEventRequest struct {
	Title              string  `json:"title"`
	EventText          string  `json:"event_text"`
	Category           string  `json:"category"`
	AIImpactAnalysis   string  `json:"ai_impact_analysis"`
	IncomeImpact       float64 `json:"income_impact"`
	ExpenseImpact      float64 `json:"expense_impact"`
	NewMaritalStatus   *string `json:"new_marital_status"`
	NewRiskScore       *int    `json:"new_risk_score"`
	NewRiskTolerance   *string `json:"new_risk_tolerance"`
	IncomeStreamsToAdd []struct {
		Name      string  `json:"name"`
		Type      string  `json:"type"`
		Amount    float64 `json:"amount"`
		IsPassive bool    `json:"is_passive"`
	} `json:"income_streams_to_add"`
	DependentsToAdd []struct {
		Name         string  `json:"name"`
		Relationship string  `json:"relationship"`
		MonthlyCost  float64 `json:"monthly_cost"`
	} `json:"dependents_to_add"`
}
