package providers

import (
	"app/src"
	"fmt"
	"github.com/go-xorm/xorm"
	"log"
	"time"
	"xorm.io/core"
)

type FetchModel struct {
	id       int // should be uint
	url      string
	interval int
}

type FetchHistoryModel struct {
	response   string
	duration   float32
	created_at time.Time
}

type DB struct {
	connection *xorm.Engine
}

func Connect(a *app.App, path string) bool {
	log.Print("Before DB Connect", " ", path)
	engine, err := xorm.NewEngine("sqlite3", path)
	log.Print("After connection")
	if err != nil {
		fmt.Print(err)
		return false
	}

	//_db := &DB{}
	//_db.connection = engine
	a.DB = engine
	//a.DB = &_db

	errSync := engine.Sync(&FetchModel{}, &FetchHistoryModel{})

	if errSync != nil {
		log.Fatal(err)
	}

	engine.ShowSQL(true)
	engine.Logger().SetLevel(core.LOG_DEBUG)
	return true
}
