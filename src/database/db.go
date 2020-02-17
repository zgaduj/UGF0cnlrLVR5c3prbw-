package database

import (
	"app/config"
	"app/src/models"
	"github.com/jinzhu/gorm"
	"log"
	"os"
	"strconv"
)
import _ "github.com/go-sql-driver/mysql"

type DBConfig struct {
	//FilePath string
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

type Store struct {
	DB *gorm.DB
}

var db *gorm.DB

func NewConnection(dbConf *config.DBConfig) (*Store, error) {
	db, err := gorm.Open("mysql", dbConf.User+":"+dbConf.Password+"@tcp("+dbConf.Host+":"+strconv.Itoa(dbConf.Port)+")/"+dbConf.Database+"?charset=utf8&parseTime=True&loc=Local")
	db.LogMode(true)
	db.SetLogger(log.New(os.Stdout, "\r\n", 0))

	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&models.FetchModel{}, &models.FetchHistoryModel{}).Error

	db.Model(&models.FetchHistoryModel{}).ModifyColumn("response", "text") // force change column type  varbinary(255) -> text

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
