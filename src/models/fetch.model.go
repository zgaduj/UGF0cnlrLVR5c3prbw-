package models

import (
	"time"
)

type BaseModel struct {
	ID        uint      `gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

type FetchModel struct {
	BaseModel
	Url      string `json:"url"`
	Interval int    `json:"interval"`
}

type FetchHistoryModel struct {
	BaseModel
	Response string  `json:"response"`
	Duration float32 `json:"duration"`
}
