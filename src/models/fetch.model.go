package models

import "github.com/jinzhu/gorm"

type FetchModel struct {
	gorm.Model
	Url      string `gorm:"unique_index:idx_url"`
	Interval int    `gorm:"unique_index:idx_url"`
}

type FetchHistoryModel struct {
	gorm.Model
	Response string
	Duration float32
	//created_at time.Time
}
