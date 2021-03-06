package server

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/baiyecha/cloud_disk/config"
	"github.com/baiyecha/cloud_disk/model"
	"github.com/baiyecha/cloud_disk/pkg/pubsub"
	"github.com/baiyecha/cloud_disk/service"
	go_file_uploader "github.com/baiyecha/go-file-uploader"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/joho/godotenv/autoload"
	"github.com/minio/minio-go"
	"github.com/spf13/afero"
	"go.uber.org/zap"
)

func setupGorm(debug bool, databaseConfig *config.DatabaseConfig) *gorm.DB {
	var dataSourceName string
	switch databaseConfig.Driver {
	case "sqlite3":
		dataSourceName = databaseConfig.DBName
	case "mysql":
		dataSourceName = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			databaseConfig.User,
			databaseConfig.Password,
			databaseConfig.Host+":"+databaseConfig.Port,
			databaseConfig.DBName,
		)
	}
	var (
		db  *gorm.DB
		err error
	)
	for i := 0; i < 10; i++ {
		db, err = gorm.Open(databaseConfig.Driver, dataSourceName)
		if err == nil {
			db.LogMode(debug)
			// group by 问题
			db.Exec("set session sql_mode='ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION'")
			if debug {
				autoMigrate(db)
			}
			return db
		}
		log.Println(err)
		time.Sleep(2 * time.Second)
	}
	log.Fatalf("数据库链接失败！ error: %+v", err)
	return nil
}

func autoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&model.User{},
		&model.Certificate{},
		&model.File{},
		&model.Group{},
		&model.Share{},
		&model.Folder{},
		&model.FolderFile{},
	).Error
	if err != nil {
		log.Fatalf("AutoMigrate 失败！ error: %+v", err)
	}
}

func setupRedis(redisConfig *config.RedisConfig) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: redisConfig.Address + ":" + redisConfig.Port,
	})
}

func setupFilesystem(fsConfig *config.FilesystemConfig) afero.Fs {
	switch fsConfig.Driver {
	case "os":
		return afero.NewBasePathFs(afero.NewOsFs(), fsConfig.Root)
	case "memory":
		return afero.NewBasePathFs(afero.NewMemMapFs(), fsConfig.Root)
	default:
		return afero.NewBasePathFs(afero.NewOsFs(), fsConfig.Root)
	}
}

func setupFileStore(s *Server) go_file_uploader.Store {
	return go_file_uploader.NewDBStore(s.DB)
}

// func setupFileUploader(s *Server) go_file_uploader.Uploader {
// 	return fileUploaderNos.NewNosUploader(
// 		go_file_uploader.HashFunc(go_file_uploader.MD5HashFunc),
// 		setupNos(s),
// 		setupFileStore(s),
// 		s.Conf.Nos.BucketName,
// 		go_file_uploader.Hash2StorageNameFunc(go_file_uploader.TwoCharsPrefixHash2StorageNameFunc),
// 		s.Conf.Nos.Endpoint,
// 		s.Conf.Nos.ExternalEndpoint,
// 	)
// }

func setupMinio(s *Server) *minio.Client {
	SslEnable := s.Conf.Minio.SSL == "true"
	minioClient, err := minio.New(
		s.Conf.Minio.Host,
		s.Conf.Minio.AccessKey,
		s.Conf.Minio.SecretKey,
		SslEnable,
	)
	if err != nil {
		log.Fatalf("minio client 创建失败! error: %+v", err)
	}
	return minioClient
}

// func setupNos(s *Server) *nosclient.NosClient {
// 	nosClient, err := nosclient.New(&nosConfig.Config{
// 		Endpoint:  s.Conf.Nos.Endpoint,
// 		AccessKey: s.Conf.Nos.AccessKey,
// 		SecretKey: s.Conf.Nos.SecretKey,
// 	})
// 	if err != nil {
// 		log.Fatalf("nos client 创建失败! error: %+v", err)
// 	}
// 	return nosClient
// }

// func setupImageUploader(s *Server) image_uploader.Uploader {
// 	nosClient := setupNos(s)
// 	return imageUploaderNos.NewNosUploader(
// 		image_uploader.HashFunc(image_uploader.MD5HashFunc),
// 		image_uploader.NewDBStore(s.DB),
// 		nosClient,
// 		s.Conf.Nos.BucketName,
// 		image_uploader.Hash2StorageNameFunc(image_uploader.TwoCharsPrefixHash2StorageNameFunc),
// 	)
// }

// func setupImageURL(s *Server) image_url.URL {
// 	return image_url.NewNosImageProxyURL(
// 		s.Conf.ImageProxy.Host,
// 		s.Conf.Nos.ExternalEndpoint,
// 		s.Conf.Nos.BucketName,
// 		s.Conf.ImageProxy.OmitBaseUrl == "true",
// 		image_uploader.Hash2StorageNameFunc(image_uploader.TwoCharsPrefixHash2StorageNameFunc),
// 	)
// }

func loadEnv(appEnv string) string {
	if appEnv == "" {
		appEnv = "production"
	}
	return appEnv
}

func setupLogger(serv *Server) *zap.Logger {
	var err error
	var logger *zap.Logger
	if serv.Debug {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}
	if err != nil {
		log.Fatal(err)
	}
	return logger
}

func SetupServer(configPath string) *Server {
	s := &Server{}
	s.AppEnv = loadEnv(os.Getenv("APP_ENV"))
	s.Debug = os.Getenv("DEBUG") == "true"
	s.Logger = setupLogger(s)
	s.Logger.Debug("load config...")
	s.Conf = config.LoadConfig(configPath)
	s.Logger.Debug("load filesystem...")
	s.BaseFs = setupFilesystem(&s.Conf.Fs)
	s.Logger.Debug("load redis...")
	s.RedisClient = setupRedis(&s.Conf.Redis)
	s.Logger.Debug("load database...")
	s.DB = setupGorm(s.Debug, &s.Conf.DB)
	s.Logger.Debug("load service...")
	s.Pub = pubsub.NewPub(s.RedisClient, s.Logger)
	s.Service = service.NewService(s.DB, s.RedisClient, s.BaseFs, s.Conf, s.Pub)
	s.Logger.Debug("load uploader service...")
	//s.FileUploader = setupFileUploader(s)
	//s.BucketName = s.Conf.Nos.BucketName
	return s
}
