package db_store

import (
	"github.com/baiyecha/cloud_disk/model"
	"github.com/jinzhu/gorm"
)

type dbShare struct {
	db *gorm.DB
}

func NewDBShare(db *gorm.DB) model.ShareStore {
	return &dbShare{db}
}
