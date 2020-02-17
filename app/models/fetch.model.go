package models

import (
	"time"
)

type BaseModel struct {
	ID        uint64    `json:"id",gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at"`
	//UpdatedAt time.Time  `json:"updated_at"`
	//DeletedAt *time.Time `json:"deleted_at",sql:"index"`
}

type FetchModel struct {
	BaseModel
	Url            string `json:"url"`
	Interval       int    `json:"interval"`
	LockedDownload bool   `json:"-"`
}

type FetchHistoryModel struct {
	BaseModel
	FetchID  uint64  `json:"fetch_id"`
	Response []byte  `json:"response",sql:"type:text"`
	Duration float32 `json:"duration"`
}
