package server

import (
	"github.com/baiyecha/cloud_disk/config"
	"github.com/baiyecha/cloud_disk/pkg/pubsub"
	"github.com/baiyecha/cloud_disk/service"
	go_file_uploader "github.com/baiyecha/go-file-uploader"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"github.com/spf13/afero"
	"go.uber.org/zap"
)

type Server struct {
	Debug      bool
	BucketName string
	AppEnv     string
	DB         *gorm.DB
	Logger     *zap.Logger
	//ImageUrl      image_url.URL
	RedisClient  *redis.Client
	Conf         *config.Config
	Service      service.Service
	Pub          pubsub.PubQueue
	FileUploader go_file_uploader.Uploader
	BaseFs       afero.Fs
}
