package main

import (
	"encoding/json"
	"log"
	"os"
	"path"
	"net/http"
	"strconv"
	"github.com/codegangsta/cli"
	"path/filepath"
	"io/ioutil"
	"github.com/keita0q/adv_db_server/database/local"
	"github.com/keita0q/adv_db_server/manager"
	"github.com/keita0q/adv_db_server/service"
	"github.com/drone/routes"
	"github.com/keita0q/adv_db_server/notification"
)

func main() {
	tApp := cli.NewApp()
	tApp.Flags = []cli.Flag{
		cli.StringFlag{
			Name: "config, c",
		},
	}

	tApp.Action = func(aContext *cli.Context) {
		var tConfJSONPath string
		if aContext.String("config") != "" {
			tConfJSONPath = aContext.String("config")
		} else {
			tRunningProgramDirectory, tError := filepath.Abs(filepath.Dir(os.Args[0]))
			if tError != nil {
				log.Println("プログラムの走っているディレクトリの取得に失敗")
				os.Exit(1)
			}
			tConfJSONPath = path.Join(tRunningProgramDirectory, "config.json")
		}

		tJSONBytes, tError := ioutil.ReadFile(tConfJSONPath)
		if tError != nil {
			log.Println(tConfJSONPath + "の読み取りに失敗")
			os.Exit(1)
		}

		tConfig := &config{}
		if tError := json.Unmarshal(tJSONBytes, tConfig); tError != nil {
			log.Println(tError)
			log.Println(tConfJSONPath + "が不正なフォーマットです。")
			os.Exit(1)
		}

		tContextPath := "/" + tConfig.ContextPath + "/"

		tDB := local.NewDatabase(tConfig.SavePath)

		tNoti := notification.New(tConfig.Urls)
		tManager, tError := manager.New(&manager.Config{
			Notification: tNoti,
			Database: tDB,
		})
		if tError != nil {
			log.Println(tError)
			os.Exit(1)
		}

		tService := service.New(&service.Config{
			Manager: tManager,
			ContextPath:  tContextPath,
			ResourcePath: tConfig.ResourcePath,
		})
		if tError != nil {
			log.Println(tError)
			os.Exit(1)
		}

		tRouter := routes.New()

		tRouter.Get(path.Join(tContextPath, "/rest/v1/advs"), tService.GetAllAdvs)
		tRouter.Get(path.Join(tContextPath, "/rest/v1/advs/:id"), tService.GetAdv)
		tRouter.Put(path.Join(tContextPath, "/rest/v1/advs/:id"), tService.Win)

		tRouter.Post(path.Join(tContextPath, "/rest/v1/advs/:id/param/g"), tService.ChangeG)
		tRouter.Post(path.Join(tContextPath, "/rest/v1/advs/:id/param/a"), tService.ChangeA)

		tRouter.Get(path.Join(tContextPath, "/.*"), tService.GetFile)

		http.Handle(tContextPath, tRouter)

		http.ListenAndServe(":" + strconv.Itoa(tConfig.Port), nil)
	}

	tApp.Run(os.Args)
	os.Exit(0)
}

type config struct {
	ContextPath  string `json:"context_path"`
	Port         int    `json:"port"`
	SavePath     string `json:"save_path"`
	ResourcePath string `json:"resource_path"`
	Urls         []string `json:"urls"`
}

