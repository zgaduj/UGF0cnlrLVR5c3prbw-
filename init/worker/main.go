package main

import (
	"app/app/database"
	"app/app/models"
	"app/config"
	"context"
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

type WorkerService struct {
	db     *gorm.DB
	config *WorkerConf
}

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

	env := &WorkerService{ // @todo: move to single config?
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
		cancelCtx()
	}

	defer dbStore.Close()

	env.db = dbStore.DB

	ticker := time.NewTicker(time.Duration(env.config.Worker.IntervalCheckURL) * time.Second) // @todo: move to env (type time.Duration)
	for {
		select {
		case <-ticker.C:
			env.getListUrls()
		}
	}

	go func() {
		<-signals
		ticker.Stop()
		cancelCtx()
	}()

	<-appCtx.Done()
}

func (worker *WorkerService) getListUrls() {
	var fetcher models.FetchModel
	where := models.FetchModel{LockedDownload: false}
	rows, err := worker.db.Model(fetcher).Where(where).Rows()

	defer rows.Close()
	if err != nil {
		log.Print(err)
		return
	}

	i := 0
	log.Print("[getListUrls] ------------------------- START LOOP")

	for rows.Next() {
		worker.db.ScanRows(rows, &fetcher)

		go worker.CheckAndRun(&fetcher)

		//log.Print("fetcher ", i, " ", fetcher)
		i += 1
	}

	log.Print("[getListUrls] ------------------------- END LOOP")
}

func (worker *WorkerService) CheckAndRun(fetcher *models.FetchModel) {
	log.Print("[CheckAndRun] -------------------------")
	var fetcherHistory models.FetchHistoryModel
	lastRun := worker.db.Model(fetcherHistory).Last(&fetcherHistory)
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

		go worker.DownloadContent(fetcher) // run async sub process

	}

}

func (worker *WorkerService) DownloadContent(fetcher *models.FetchModel) {

	defer func() {
		log.Print("*** UNLOCKING FETCHER ", fetcher.Url)
		//fetcher.LockedDownload = false
		//worker.db.Model(&fetcher).Updates(models.FetchModel{LockedDownload: false}) // not working with false value
		worker.db.Model(&fetcher).Update("LockedDownload", false)
		time.Sleep(2 * time.Second)
	}()

	//fetcher.LockedDownload = true
	log.Print("*** LOCKING FETCHER ", fetcher.Url)
	worker.db.Model(&fetcher).Updates(models.FetchModel{LockedDownload: true})

	startRequest := time.Now()
	client := http.Client{Timeout: time.Duration(worker.config.Worker.Timeout) * time.Second} // @todo: move to worker (type time.Duration)
	resp, err := client.Get(fetcher.Url)
	diff := time.Since(startRequest)
	log.Print("DOWNLOAD DURATION ", diff.Milliseconds(), diff.Seconds())

	if err != nil {
		log.Print("Error downloading site content: ", fetcher.Url, " ", err)
		log.Print(time.Duration(worker.config.Worker.Timeout) * time.Second)
		worker.db.Model(&fetcher).Updates(models.FetchModel{LockedDownload: false})
		return
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
	historyError := worker.db.Save(&saveHistory)
	if historyError.Error != nil {
		log.Print(historyError.Error.Error())
	}
	//fetcher.LockedDownload = false
	//worker.db.Save(&fetcher)

}
