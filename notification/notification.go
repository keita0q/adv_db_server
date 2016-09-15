package notification

import (
	"net/url"
	"path"
	"net/http"
	"encoding/json"
	"bytes"
	"github.com/keita0q/adv_db_server/model"
)

type Notification struct {
	Urls []string
}

func New(aUrls []string) *Notification {
	return &Notification{Urls:aUrls}
}

func (aNoti *Notification) NotifyUpdate(aAdvs *map[string]model.Advertiser) error {
	for _, tUrl := range aNoti.Urls {
		go func(aUrl string) error {
			tURL, tError := url.Parse(aUrl)
			if tError != nil {
				return tError
			}
			tURL.Path = path.Join(tURL.Path, "advs")
			tBytes, tError := json.Marshal(aAdvs)
			if tError != nil {
				return tError
			}

			tRequest, tError := http.NewRequest("PUT", tURL.String(), bytes.NewBuffer(tBytes))
			if tError != nil {
				return tError
			}

			tResponse, tError := http.DefaultClient.Do(tRequest)
			if tError != nil {
				return tError
			}
			defer tResponse.Body.Close()
			return nil
		}(tUrl)
	}
	return nil
}

func (aNoti *Notification) NotifyParamG(aID string, aG *model.Param) error {
	for _, tUrl := range aNoti.Urls {
		go func(aUrl string) error {
			tURL, tError := url.Parse(aUrl)
			if tError != nil {
				return tError
			}
			tURL.Path = path.Join(tURL.Path, "advs", aID, "g")

			tBytes, tError := json.Marshal(aG)
			if tError != nil {
				return tError
			}

			tRequest, tError := http.NewRequest("PUT", tURL.String(), bytes.NewBuffer(tBytes))
			if tError != nil {
				return tError
			}

			tResponse, tError := http.DefaultClient.Do(tRequest)
			if tError != nil {
				return tError
			}
			defer tResponse.Body.Close()
			return nil
		}(tUrl)
	}
	return nil
}

func (aNoti *Notification) NotifyParamA(aID string, aA *model.Param) error {
	for _, tUrl := range aNoti.Urls {
		go func(aUrl string) error {
			tURL, tError := url.Parse(aUrl)
			if tError != nil {
				return tError
			}
			tURL.Path = path.Join(tURL.Path, "advs", aID, "a")
			tBytes, tError := json.Marshal(aA)
			if tError != nil {
				return tError
			}

			tRequest, tError := http.NewRequest("PUT", tURL.String(), bytes.NewBuffer(tBytes))
			if tError != nil {
				return tError
			}

			tResponse, tError := http.DefaultClient.Do(tRequest)
			if tError != nil {
				return tError
			}
			defer tResponse.Body.Close()
			return nil
		}(tUrl)
	}
	return nil
}
