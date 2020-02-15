package handlers

import (
	"app/src/models"
	"encoding/json"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
)

func FetcherSave(db *gorm.DB, w http.ResponseWriter, r *http.Request) {

	var fetcherBody models.FetchModel
	var fetcher models.FetchModel

	//body, err := ioutil.ReadAll(r.Body)
	//if err != nil {
	//	panic("1")
	//}
	//
	//log.Print(string(body))
	//
	//err = json.Unmarshal(body, &fetcherBody)
	//if err != nil {
	//	panic("2")
	//}
	//log.Print(fetcherBody)

	err := json.NewDecoder(r.Body).Decode(&fetcherBody)

	log.Print("fetcherBody.ID ", fetcherBody.ID)

	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(""))
		return
	} else {

		fetcher.Url = fetcherBody.Url
		fetcher.Interval = fetcherBody.Interval

		//newFetcher := models.FetchModel{
		//	Url:      string("https://httpbin.org/range/15"),
		//	Interval: 20,
		//}

		if fetcherBody.ID == 0 { // create new fetcher
			result := db.Create(&fetcher)
			if result.Error != nil {
				log.Print("Cannot create new record")
			} else {
				json.NewEncoder(w).Encode(struct {
					id uint
				}{id: fetcher.ID})
			}
			return
		} else {
			//var foundFetcher models.FetchModel

			log.Print("ID is not null")

			find := db.Where("id = ?", fetcherBody.ID).First(&fetcher)
			findError := find.Error
			if findError != nil {
				w.WriteHeader(400)
				log.Print(findError.Error())
				w.Write([]byte(findError.Error()))

				return
			}
			if find.RecordNotFound() {
				log.Print("record not found")
				w.WriteHeader(404)
			} else {
				log.Print("fetcherBody")
				log.Print(fetcherBody)
				saveError := find.Save(&fetcherBody).Error

				if saveError != nil {
					w.WriteHeader(404)
					w.Write([]byte(saveError.Error()))
					return
				} else {
					enc, err := json.Marshal(struct {
						id uint
					}{
						id: fetcherBody.ID,
					})
					if err != nil {
						w.Write([]byte("dupa"))
						return
					}
					w.Write(enc)
				}
				//log.Print("save state: ", saveState, saveState.Error)
				//db.Save(fetcher)
			}
		}

		return

		//query := db.Save(&fetcher)
		//if err != nil {
		//
		//}

		//if query.Error != nil {
		//	log.Print(query.Error)
		//}

		//log.Print(fetcher.Url)
		//log.Print(fetcher.Interval)
		//log.Print(fetcher)

		//w.Write([]byte("asd"))
	}
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
