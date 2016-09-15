package database

import "github.com/keita0q/adv_db_server/model"

type Database interface {
	LoadAllAdvertiser() (map[string]model.Advertiser, error)
	LoadAdvertiser(string) (*model.Advertiser, error)
	SaveAdvertiser(*model.Advertiser) error
	SaveAllAdvertisers(*map[string]model.Advertiser) error
}

type NotFoundError struct {
	message string
}

func NewNotFoundError(aMessage string) *NotFoundError {
	return &NotFoundError{message: aMessage}
}

func (aNotFoundError *NotFoundError) Error() string {
	return aNotFoundError.message
}

type  Budgets struct {
	AdvBudget map[string]*BudgetInfo
}

type BudgetInfo struct {
	Budget int `json:"budget"`
	Cpc    int `json:"cpc"`
}

type NgDomains struct {
	AdvNgs map[string][]string
}