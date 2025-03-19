package model

type Fact struct {
	PeriodStart         string
	PeriodEnd           string
	PeriodKey           string
	IndicatorToMoID     int
	IndicatorToMoFactID int
	Value               int
	FactTime            string
	IsPlan              int
	AuthUserID          int
	Comment             string
}
