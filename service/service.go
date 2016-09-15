package service

import (
	"net/http"
	"net/url"
	"encoding/json"
	"log"
	"io/ioutil"
	"strconv"
	"strings"
	"fmt"
	"path/filepath"
	"github.com/keita0q/adv_db_server/database"
	"github.com/keita0q/adv_db_server/manager"
	"github.com/keita0q/adv_db_server/model"
)

type Service struct {
	manager      *manager.Manager
	contextPath  string
	resourcePath string
}

type Config struct {
	Manager      *manager.Manager
	ContextPath  string
	ResourcePath string
}

func New(aConfig *Config) *Service {
	return &Service{
		manager: aConfig.Manager,
		contextPath:  aConfig.ContextPath,
		resourcePath: aConfig.ResourcePath,
	}
}

func (aService *Service) GetFile(aWriter http.ResponseWriter, aRequest *http.Request) {
	tPath := strings.TrimPrefix(aRequest.RequestURI, aService.contextPath)
	fmt.Println(tPath)
	if i := strings.Index(tPath, "?"); i > 0 {
		tPath = tPath[:i]
	}
	http.ServeFile(aWriter, aRequest, filepath.Join(aService.resourcePath, tPath))
}

func (aService *Service) GetAllAdvs(aWriter http.ResponseWriter, aRequest *http.Request) {
	aService.handler(func(aQueries url.Values, aRequestBody []byte) (int, interface{}, error) {
		tAdvs := aService.manager.GetAllAdvs()
		return http.StatusOK, tAdvs, nil
	})(aWriter, aRequest)
}

func (aService *Service) GetAdv(aWriter http.ResponseWriter, aRequest *http.Request) {
	aService.handler(func(aQueries url.Values, aRequestBody []byte) (int, interface{}, error) {
		tID := aQueries.Get(":id")
		tAdv := aService.manager.GetAdv(tID)
		if tAdv == nil {
			return http.StatusBadRequest, nil, nil
		}
		return http.StatusOK, tAdv, nil
	})(aWriter, aRequest)
}

//type cost struct {
//	Value  float64`json:"value"`
//}

func (aService *Service) Win(aWriter http.ResponseWriter, aRequest *http.Request) {
	aService.handler(func(aQueries url.Values, aRequestBody []byte) (int, interface{}, error) {
		//tCost := &cost{}
		tID := aQueries.Get(":id")
		//if tError := json.Unmarshal(aRequestBody, tCost); tError != nil {
		//	return http.StatusBadRequest, nil, tError
		//}

		//tAdv, tError := aService.manager.UpdateAdv(tID, tCost.Value)
		tAdv, tError := aService.manager.Click(tID)
		if tError != nil {
			return http.StatusInternalServerError, nil, tError
		}
		return http.StatusOK, tAdv, nil

	})(aWriter, aRequest)
}

func (aService *Service) ChangeG(aWriter http.ResponseWriter, aRequest *http.Request) {
	aService.handler(func(aQueries url.Values, aRequestBody []byte) (int, interface{}, error) {
		tParam := &model.Param{}
		tID := aQueries.Get(":id")
		if tError := json.Unmarshal(aRequestBody, tParam); tError != nil {
			return http.StatusBadRequest, nil, tError
		}

		if tError := aService.manager.ChangeG(tID, tParam); tError != nil {
			return http.StatusInternalServerError, nil, tError
		}
		return http.StatusOK, nil, nil
	})(aWriter, aRequest)
}

func (aService *Service) ChangeA(aWriter http.ResponseWriter, aRequest *http.Request) {
	aService.handler(func(aQueries url.Values, aRequestBody []byte) (int, interface{}, error) {
		tParam := &model.Param{}
		tID := aQueries.Get(":id")
		if tError := json.Unmarshal(aRequestBody, tParam); tError != nil {
			return http.StatusBadRequest, nil, tError
		}

		if tError := aService.manager.ChangeA(tID, tParam); tError != nil {
			return http.StatusInternalServerError, nil, tError
		}
		return http.StatusOK, nil, nil
	})(aWriter, aRequest)
}

func handleError(aError error) int {
	if _, ok := aError.(*database.NotFoundError); ok {
		return http.StatusNotFound
	}
	return http.StatusInternalServerError
}

func (aService *Service) handler(aAPI func(url.Values, []byte) (int, interface{}, error)) func(http.ResponseWriter, *http.Request) {
	return func(aWriter http.ResponseWriter, aRequest *http.Request) {
		log.Printf("[INFO] access:%s", aRequest.RequestURI)
		defer aRequest.Body.Close()

		tResponseBody, tError := ioutil.ReadAll(aRequest.Body)
		if tError != nil {
			http.Error(aWriter, tError.Error(), http.StatusBadRequest)
		}
		tStatusCode, tResult, tError := aAPI(aRequest.URL.Query(), tResponseBody)
		if tError != nil {
			http.Error(aWriter, tError.Error(), tStatusCode)
			return
		}

		if tStatusCode == http.StatusNoContent {
			aWriter.WriteHeader(http.StatusNoContent)
			return
		}

		tBytes, tError := json.MarshalIndent(tResult, "", "  ")
		if tError != nil {
			http.Error(aWriter, tError.Error(), tStatusCode)
			return
		}

		aWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
		aWriter.Header().Set("Content-Length", strconv.Itoa(len(tBytes)))
		aWriter.Header().Set("Access-Control-Allow-Origin", "*")
		aWriter.WriteHeader(tStatusCode)
		aWriter.Write(tBytes)
	}
}
