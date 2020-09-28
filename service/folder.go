package service

import (
	"context"
	"github.com/baiyecha/cloud_disk/model"
)

type folderService struct {
	model.FolderStore
}

func ListFolder(ctx context.Context, folderIds []int64, userId int64) (folder []*model.Folder, err error) {
	return FromContext(ctx).ListFolder(folderIds, userId)
}

func LoadFolder(ctx context.Context, id, userId int64, isLoadRelated bool) (folder *model.Folder, err error) {
	return FromContext(ctx).LoadFolder(id, userId, isLoadRelated)
}
func LoadSimpleFolder(ctx context.Context, id, userId int64) (folder *model.SimpleFolder, err error) {
	return FromContext(ctx).LoadSimpleFolder(id, userId)
}

func CreateFolder(ctx context.Context, folder *model.Folder) (err error) {
	return FromContext(ctx).CreateFolder(folder)
}

func ExistFolder(ctx context.Context, userId, parentId int64, folderName string) (isExist bool) {
	return FromContext(ctx).ExistFolder(userId, parentId, folderName)
}

func DeleteFolder(ctx context.Context, ids []int64, userId int64) (allowDelFileHashList []string, err error) {
	return FromContext(ctx).DeleteFolder(ids, userId)
}

func MoveFolder(ctx context.Context, to *model.Folder, ids []int64) (err error) {
	return FromContext(ctx).MoveFolder(to, ids)
}

func CopyFolder(ctx context.Context, to *model.Folder, foders []*model.Folder) (totalSize uint64, err error) {
	return FromContext(ctx).CopyFolder(to, foders)
}

func RenameFolder(ctx context.Context, id, currentFolderId int64, newName string) (err error) {
	return FromContext(ctx).RenameFolder(id, currentFolderId, newName)
}

func NewFolderService(ds model.FolderStore) model.FolderService {
	return &folderService{ds}
}
