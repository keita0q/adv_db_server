package manager

import (
	"github.com/keita0q/adv_db_server/database"
	"github.com/keita0q/adv_db_server/model"
	"sync"
	"github.com/keita0q/adv_db_server/notification"
	"fmt"
)

type Manager struct {
	sync.Mutex
	advertisers  map[string]model.Advertiser
	notification *notification.Notification
	database     database.Database
}

type Config struct {
	Notification *notification.Notification
	Database     database.Database
}

func New(aConfig *Config) (*Manager, error) {
	tAdvertisers, tError := aConfig.Database.LoadAllAdvertiser()
	if tError != nil {
		return nil, tError
	}
	return &Manager{advertisers: tAdvertisers, notification:aConfig.Notification, database:aConfig.Database}, nil
}

func (aManager *Manager) GetAllAdvs() map[string]model.Advertiser {
	return aManager.advertisers
}

func (aManager *Manager) GetAdv(aID string) *model.Advertiser {
	aManager.Lock()
	defer aManager.Unlock()
	tAdv := aManager.advertisers[aID]
	return &tAdv
}

//func (aManager *Manager) UpdateAdv(aID string, aCost float64) (*model.Advertiser, error) {
//	aManager.Lock()
//	defer aManager.Unlock()
//	tAdv, tOK := aManager.advertisers[aID]
//	if !tOK {
//		return nil, database.NewNotFoundError(aID + "は存在しない")
//	}
//
//	tAdv.Budget = tAdv.Budget - int(aCost)
//
//	aManager.advertisers[aID] = tAdv
//
//	return &tAdv, nil
//}

func (aManager *Manager) Click(aID string) (*model.Advertiser, error) {
	aManager.Lock()
	defer aManager.Unlock()
	tAdv, tOK := aManager.advertisers[aID]
	if !tOK {
		return nil, database.NewNotFoundError(aID + "は存在しない")
	}

	tAdv.Budget = tAdv.Budget - tAdv.Cpc

	aManager.advertisers[aID] = tAdv

	fmt.Println(tAdv)
	go aManager.database.SaveAllAdvertisers(&aManager.advertisers)

	if tAdv.Budget < tAdv.Cpc * 16 {
		delete(aManager.advertisers, aID)
		go aManager.notification.NotifyUpdate(&aManager.advertisers)
	}

	return &tAdv, nil
}

func (aManager *Manager) ChangeG(aID string, aParam *model.Param) error {
	return aManager.notification.NotifyParamG(aID, aParam)
}

func (aManager *Manager) ChangeA(aID string, aParam *model.Param) error {
	return aManager.notification.NotifyParamA(aID, aParam)
}