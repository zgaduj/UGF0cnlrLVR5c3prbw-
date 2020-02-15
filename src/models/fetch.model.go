package models

import (
	"time"
)

type BaseModel struct {
	ID        uint       `json:"id",gorm:"primary_key"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at",sql:"index"`
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
