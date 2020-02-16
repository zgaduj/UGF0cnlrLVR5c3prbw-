package main

import (
	"app/config"
	"app/src/database"
	"app/src/models"
	"fmt"
	"github.com/jinzhu/gorm"
	"log"
	"strconv"
	"time"
)

type Env struct {
	db     *gorm.DB
	config config.Config
}

type WorkerActionType string

const (
	WActionRequest  WorkerActionType = "request"
	WActionResponse                  = "response"
	WActionError                     = "error"
)

type WorkerRequest struct {
	Action  WorkerActionType
	Payload []byte
}

var WorkerChannel = make(chan WorkerRequest)

func main() {
	log.Print("[WORKER] START")
	_config := config.GetConfig()

	env := &Env{
		config: *_config,
	}

	dbStore, error := database.NewConnection(env.config.DB)
	if error != nil {
		log.Panic(error)
	} else {

		defer dbStore.Close()

		env.db = dbStore.DB

		ticker := time.NewTicker(5000 * time.Millisecond)
		done := make(chan bool)
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				fmt.Println("Tick at", t)

				env.getListUrls()

				//WorkerChannel <- WorkerRequest{
				//	Action:  WActionResponse,
				//	Payload: []byte(t.Format("s")),
				//}

			}
		}

		//time.Sleep(1000 * time.Millisecond)
		//ticker.Stop()
		//done <- true
		//fmt.Println("Ticker stopped")
	}
}

func (env *Env) getListUrls() {
	//var listFetchers []models.FetchModel
	var fetcher models.FetchModel
	where := models.FetchModel{LockedDownload: false}
	rows, err := env.db.Model(fetcher).Where(where).Rows()

	defer rows.Close()
	if err != nil {
		log.Print(err)
		return
	}

	i := 0
	for rows.Next() {
		env.db.ScanRows(rows, &fetcher)

		env.CheckAndRun(&fetcher)

		log.Print("fetcher ", i, " ", fetcher)
		i += 1
	}

	log.Print("-------------------------")
}

func (env *Env) CheckAndRun(fetcher *models.FetchModel) {
	var fetcherHistory models.FetchHistoryModel
	lastRun := env.db.Model(fetcherHistory).Last(fetcherHistory)
	run := false
	if lastRun.RecordNotFound() {
		run = true
	} else {
		if fetcher.LockedDownload == false {
			//lastTime := fetcherHistory.CreatedAt
			//nowTime := time.Now()
			//diff := lastTime.Sub(nowTime)
			//fetcherHistory.CreatedAt.
			//log.Print("@@@@@@@@@@@@ diff ", diff.Seconds())
			log.Print("@@@@@@@@@@@@ diff ")
		} else {
			log.Print("[Fetcher ID:" + strconv.FormatUint(fetcher.ID, 10) + "] is locked ")

		}
	}

	if run {
		log.Print("[Fetcher ID:" + strconv.FormatUint(fetcher.ID, 10) + "] RUN ")
		fetcher.LockedDownload = true
		env.db.Save(fetcher)
		time.Sleep(2000 * time.Millisecond)
		saveHistory := models.FetchHistoryModel{
			FetchID:  fetcher.ID,
			Response: "adsdfsdfsf",
			Duration: 0,
		}
		//saveHistory.CreatedAt = time.Now()
		historyError := env.db.Save(&saveHistory)
		if historyError.Error != nil {
			log.Print(historyError.Error.Error())
		}
		fetcher.LockedDownload = false
		env.db.Save(fetcher)

	}

}
