package store

import (
	"github.com/baiyecha/cloud_disk/model"
	"github.com/baiyecha/cloud_disk/store/db_store"
	"github.com/baiyecha/cloud_disk/store/redis_store"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
)

type Store interface {
	model.TicketStore
	model.UserStore
	model.CertificateStore
	model.FileStore
	model.ShareStore
	model.FolderStore
	model.GroupStore
	model.FolderFileStore
}

type store struct {
	model.TicketStore
	model.UserStore
	model.CertificateStore
	model.FileStore
	model.ShareStore
	model.FolderStore
	model.GroupStore
	model.FolderFileStore
}

func NewStore(db *gorm.DB, redisClient *redis.Client) Store {
	return &store{
		redis_store.NewRedisTicket(redisClient),
		db_store.NewDBUser(db),
		db_store.NewDBCertificate(db),
		db_store.NewDBFile(db),
		db_store.NewDBShare(db),
		db_store.NewDBFolder(db),
		db_store.NewDBGroup(db),
		db_store.NewDBFolderFile(db),
	}
}
