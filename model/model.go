package model

import "sync"

type Advertiser struct {
	sync.Mutex
	ID        string `json:"id"`
	Budget    int `json:"budget"`
	Cpc       int `json:"cpc"`
	NgDomains []string `json:"ngdomains"`
}

type Param struct {
	Value float64 `json:"value"`
}
