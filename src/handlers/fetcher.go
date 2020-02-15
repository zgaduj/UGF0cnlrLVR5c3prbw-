package handlers

import (
	"app/src/models"
	"encoding/json"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
)

func FetcherSave(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	newFetcher := models.FetchModel{
		Url:      string("https://httpbin.org/range/15"),
		Interval: 20,
	}

	//exists := db.Find(&newFetcher)

	query := db.Save(&newFetcher)
	//if err != nil {
	//
	//}

	if query.Error != nil {
		log.Print(query.Error)
	}

	w.Write([]byte("..."))

}

func FetcherDelete(w http.ResponseWriter, r *http.Request) {

}

func FetcherGet(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var fetchers []models.FetchModel

	list, err := db.Find(&fetchers).Rows()
	if err != nil {
		log.Print(err)
		w.Write([]byte("nope"))
	}

	encoded, err2 := json.Marshal(list)
	if err2 != nil {
		log.Print(err2)
		w.Write([]byte("nope2"))
	}

	w.Write(encoded)
}

func FetcherHistory(w http.ResponseWriter, r *http.Request) {

}
