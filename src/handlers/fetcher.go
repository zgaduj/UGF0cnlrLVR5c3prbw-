package handlers

import (
	"app/src/models"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
)

func FetcherSave(db *gorm.DB, w http.ResponseWriter, r *http.Request) {

	var fetcherBody models.FetchModel
	var fetcher models.FetchModel

	err := json.NewDecoder(r.Body).Decode(&fetcherBody)

	//log.Print("fetcherBody.ID ", fetcherBody.ID)

	if err != nil {
		EncodeErrorMessage(w, err, 400)
		return
	}

	fetcher.Url = fetcherBody.Url
	fetcher.Interval = fetcherBody.Interval

	if fetcherBody.ID == 0 {
		result := db.Create(&fetcher)

		EncodeOrError(EncodeOrErrorInterface{
			Write:     w,
			Error:     result.Error,
			ErrorCode: 400,
			Encode: struct {
				ID uint64 `json:"id"`
			}{
				ID: fetcher.ID,
			},
		})
		return

	}

	find := db.Where("id = ?", fetcherBody.ID).First(&fetcher)
	findError := find.Error
	if findError != nil {
		EncodeErrorMessage(w, findError, 400)
		return
	}
	if find.RecordNotFound() {
		EncodeErrorMessage(w, errors.New("Not found"), 404)
		return
	}

	EncodeOrError(EncodeOrErrorInterface{
		Write:     w,
		Error:     find.Save(&fetcherBody).Error,
		ErrorCode: 400,
		Encode: struct {
			ID uint64 `json:"id"`
		}{
			ID: fetcherBody.ID,
		},
	})

	return
}

func FetcherDelete(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		EncodeErrorMessage(w, err, 400)
		return
	}

	id := chi.URLParam(r, "id")

	var fetcher models.FetchModel

	find := tx.Where("id = ?", id).First(&fetcher)
	if find.RecordNotFound() {
		EncodeErrorMessage(w, errors.New("Not found"), 404)
		return
	}

	del := find.Delete(&fetcher)
	if del.Error != nil {
		EncodeErrorMessage(w, del.Error, 400)
		return
	}

	var fetcherHistory models.FetchHistoryModel

	delHistory := tx.Where("fetch_id = ?", id).Delete(&fetcherHistory)
	if delHistory.Error != nil {
		EncodeErrorMessage(w, delHistory.Error, 400)
		return
	}

	commit := tx.Commit()

	EncodeOrError(EncodeOrErrorInterface{
		Write:     w,
		Error:     commit.Error,
		ErrorCode: 400,
		Encode: struct {
			ID int `json:"id"`
		}{
			ID: int(fetcher.ID),
		},
	})
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
