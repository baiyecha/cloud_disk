package service

import (
	"context"
	"github.com/baiyecha/cloud_disk/model"
)

type folderFileService struct {
	model.FolderFileStore
}

func LoadFolderFilesByFolderIds(ctx context.Context, folderIds []int64, userId int64) (folderFiles []*model.WrapFolderFile, err error) {
	return FromContext(ctx).LoadFolderFilesByFolderIds(folderIds, userId)
}

func LoadFolderFilesByFolderIdAndFileIds(ctx context.Context, folderId int64, fileIds []int64, userId int64) (folderFiles []*model.WrapFolderFile, err error) {
	return FromContext(ctx).LoadFolderFilesByFolderIdAndFileIds(folderId, fileIds, userId)
}

func ExistFile(ctx context.Context, filename string, folderId, userId int64) (isExist bool, err error) {
	return FromContext(ctx).ExistFile(filename, folderId, userId)
}

func NewFolderFileService(fs model.FolderFileStore) model.FolderFileService {
	return &folderFileService{fs}
}
