package database

import (
	"app/src/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
	"os"
)

type Store struct {
	DB *gorm.DB
}

var db *gorm.DB

func NewConnection(path string) (*Store, error) {
	db, err := gorm.Open("sqlite3", path)
	db.LogMode(true)
	db.SetLogger(log.New(os.Stdout, "\r\n", 0))

	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&models.FetchModel{}).Error
	if err != nil {
		log.Panic(err)
	}

	return &Store{
		DB: db,
	}, nil
}

func (s *Store) Close() {
	s.DB.Close()
}

func (s *Store) GetDB() *gorm.DB {
	return s.DB
}
