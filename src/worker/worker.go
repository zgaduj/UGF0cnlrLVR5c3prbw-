package worker

import (
	"app/config"
	"app/src/database"
	"app/src/models"
	"fmt"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

type Env struct {
	db     *gorm.DB
	config config.Config
}

func RunWorker() {

	_config := config.GetConfig()

	env := &Env{
		config: *_config,
	}

	dbStore, error := database.NewConnection(env.config.DB.FilePath)
	if error != nil {
		log.Panic(error)
	} else {

		defer dbStore.Close()

		env.db = dbStore.DB

		ticker := time.NewTicker(2000 * time.Millisecond)
		done := make(chan bool)
		go func() {
			for {
				select {
				case <-done:
					return
				case t := <-ticker.C:
					fmt.Println("Tick at", t)
					env.getListUrls()
				}
			}
		}()

		//time.Sleep(1000 * time.Millisecond)
		//ticker.Stop()
		//done <- true
		//fmt.Println("Ticker stopped")

	}
}

func (env *Env) getListUrls() {
	var listFetchers []models.FetchModel
	find := env.db.Where("locked_download", false).Find(&listFetchers)
	if find.Error != nil {
		log.Panic(find.Error)
		return
	}
	log.Print("listFetchers", listFetchers)
}
