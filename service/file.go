package service

import (
	"context"
	"github.com/baiyecha/cloud_disk/model"
)

type fileService struct {
	model.FileStore
}

func SaveFileToFolder(ctx context.Context, file *model.File, folderId int64) (err error) {
	return FromContext(ctx).SaveFileToFolder(file, folderId)
}

func DeleteFile(ctx context.Context, ids []int64, folderId int64) (allowDelFileHashList []string, err error) {
	return FromContext(ctx).DeleteFile(ids, folderId)
}

func MoveFile(ctx context.Context, fromId, toId int64, fileIds []int64) (err error) {
	return FromContext(ctx).MoveFile(fromId, toId, fileIds)
}

func CopyFile(ctx context.Context, fromId, toId int64, fileIds []int64) (totalSize uint64, err error) {
	return FromContext(ctx).CopyFile(fromId, toId, fileIds)
}

func RenameFile(ctx context.Context, folderId, fileId int64, newName string) (err error) {
	return FromContext(ctx).RenameFile(folderId, fileId, newName)
}

func LoadFile(ctx context.Context, folderId, fileId, userId int64) (file *model.File, err error) {
	return FromContext(ctx).LoadFile(folderId, fileId, userId)
}

func NewFileService(fs model.FileStore) model.FileService {
	return &fileService{fs}
}
