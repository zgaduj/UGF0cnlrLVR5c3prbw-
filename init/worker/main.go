package main

import (
	"app/app/database"
	"app/app/models"
	"app/config"
	"context"
	"fmt"
	"github.com/caarlos0/env"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"
)

type WorkerConf struct {
	Worker *config.WorkerConfig
	DB     *config.DBConfig
}

type Env struct {
	db     *gorm.DB
	config *WorkerConf
}

//type WorkerActionType string

//const (
//	WActionRequest  WorkerActionType = "request"
//	WActionResponse                  = "response"
//	WActionError                     = "error"
//)

//type WorkerRequest struct {
//	Action  WorkerActionType
//	Payload []byte
//}

//var WorkerChannel = make(chan WorkerRequest)

func main() {
	log.Print("[WORKER] START")
	workerConfig := &config.WorkerConfig{}
	if err := env.Parse(workerConfig); err != nil {
		log.Fatal("[ENV-WORKER] Error parse", err)
	}

	dbConfig := &config.DBConfig{}
	if err := env.Parse(dbConfig); err != nil {
		log.Fatal("[ENV-DB] Error parse", err)
	}

	signals := make(chan os.Signal)
	signal.Notify(signals, os.Interrupt)

	env := &Env{ // @todo: move to single config?
		config: &WorkerConf{
			Worker: workerConfig,
			DB:     dbConfig,
		},
	}

	appCtx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	dbStore, error := database.NewConnection(env.config.DB)
	if error != nil {
		log.Panic(error)
	} else {

		defer dbStore.Close()

		env.db = dbStore.DB

		ticker := time.NewTicker(time.Duration(env.config.Worker.IntervalCheckURL) * time.Second) // @todo: move to env (type time.Duration)
		done := make(chan bool)
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				fmt.Println("Tick at", t)
				env.getListUrls()
			}
		}
	}

	go func() {
		<-signals
		cancelCtx()
	}()

	<-appCtx.Done()
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
	log.Print("[getListUrls] ------------------------- START LOOP")

	for rows.Next() {
		env.db.ScanRows(rows, &fetcher)

		go env.CheckAndRun(&fetcher)

		log.Print("fetcher ", i, " ", fetcher)
		i += 1
	}

	log.Print("[getListUrls] ------------------------- END LOOP")
}

func (env *Env) CheckAndRun(fetcher *models.FetchModel) {
	log.Print("[CheckAndRun] -------------------------")
	var fetcherHistory models.FetchHistoryModel
	lastRun := env.db.Model(fetcherHistory).Last(&fetcherHistory)
	run := false
	if lastRun.RecordNotFound() {
		log.Print("[CheckAndRun] RecordNotFound")
		run = true
	} else {
		log.Print("[CheckAndRun] RecordFound")
		if fetcher.LockedDownload == false {
			lastTime := fetcherHistory.CreatedAt
			diff := time.Since(lastTime)
			log.Print("[CheckAndRun] diff ", diff.Seconds())

			if int(diff.Seconds()) >= fetcher.Interval {
				log.Print("[CheckAndRun] set flag run = true")
				run = true
			}
		} else {
			//log.Print("[CheckAndRun] url locked")
			log.Print("[CheckAndRun][Fetcher ID:" + strconv.FormatUint(fetcher.ID, 10) + "] is locked ")
		}
	}

	if run {
		log.Print("[Fetcher ID:" + strconv.FormatUint(fetcher.ID, 10) + "] RUN ")

		go env.DownloadContent(fetcher) // run async sub process

	}

}

func (env *Env) DownloadContent(fetcher *models.FetchModel) {
	fetcher.LockedDownload = true
	env.db.Save(fetcher)

	defer func(fetcher *models.FetchModel) {
		log.Print("*** UNLOCKING FETCHER ", fetcher.Url)
		fetcher.LockedDownload = false
		env.db.Save(&fetcher)
	}(fetcher)

	startRequest := time.Now()
	client := http.Client{Timeout: time.Duration(env.config.Worker.Timeout) * time.Second} // @todo: move to env (type time.Duration)
	resp, err := client.Get(fetcher.Url)
	diff := time.Since(startRequest)
	log.Print("DOWNLOAD DURATION ", diff.Milliseconds(), diff.Seconds())

	if err != nil {
		log.Print("Error downloading site content: ", fetcher.Url)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print("Error getting body: ", fetcher.Url)
	}
	saveHistory := models.FetchHistoryModel{
		FetchID:  fetcher.ID,
		Response: body,
		Duration: float32(diff.Microseconds()) / float32(time.Millisecond),
	}
	//saveHistory.CreatedAt = time.Now()
	historyError := env.db.Save(&saveHistory)
	if historyError.Error != nil {
		log.Print(historyError.Error.Error())
	}
	//fetcher.LockedDownload = false
	//env.db.Save(&fetcher)

}
