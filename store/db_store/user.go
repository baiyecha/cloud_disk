package db_store

import (
	"github.com/baiyecha/cloud_disk/errors"
	"github.com/baiyecha/cloud_disk/model"
	"github.com/jinzhu/gorm"
)

type dbUser struct {
	db *gorm.DB
}

func (u *dbUser) UserUpdateUsedStorage(userId int64, storage uint64, operator string) (err error) {
	if userId <= 0 {
		return model.ErrUserNotExist
	}
	switch operator {
	case "+":
	case "-":
	default:
		return model.ErrorOperatorNotValid
	}
	return u.db.Model(model.User{Id: userId}).
		UpdateColumn("used_storage", gorm.Expr("used_storage "+operator+" ?", storage)).Error
}

func (u *dbUser) UserExist(id int64) (bool, error) {
	var count uint8
	err := u.db.Model(model.User{}).Where(model.User{Id: id}).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (u *dbUser) UserIsNotExistErr(err error) bool {
	return model.UserIsNotExistErr(err)
}

func (u *dbUser) UserLoad(id int64) (user *model.User, err error) {
	if id <= 0 {
		return nil, model.ErrUserNotExist
	}
	user = &model.User{}
	err = u.db.Where(model.User{Id: id}).First(user).Error
	if gorm.IsRecordNotFoundError(err) {
		err = model.ErrUserNotExist
	}
	return
}

func (u *dbUser) UserLoadAndRelated(userId int64) (user *model.User, err error) {
	user, err = u.UserLoad(userId)
	if err != nil {
		return
	}
	group := &model.Group{}
	err = u.db.Where("id = ?", user.GroupId).First(&group).Error
	if gorm.IsRecordNotFoundError(err) {
		err = errors.RecordNotFound("用户组不存在")
	}
	user.Group = group
	return
}

func (u *dbUser) UserUpdate(userId int64, data map[string]interface{}) error {
	if userId <= 0 {
		return model.ErrUserNotExist
	}
	return u.db.Model(model.User{Id: userId}).
		Select(
			"name", "gender", "password", "is_ban", "group_id",
			"is_admin", "nickname", "email", "avatar_hash", "profile",
		).
		Updates(data).Error
}

func (u *dbUser) UserCreate(user *model.User) (err error) {
	err = u.db.Create(&user).Error
	return
}

func (u *dbUser) UserListByUserIds(userIds []interface{}) (users []*model.User, err error) {
	if len(userIds) == 0 {
		return
	}
	users = make([]*model.User, 0, len(userIds))
	err = u.db.Where("id in (?)", userIds).
		Set("gorm:auto_preload", true).
		Find(&users).
		Error
	return
}

func (u *dbUser) UserList(offset, limit int64) (users []*model.User, count int64, err error) {
	users = make([]*model.User, 0, 10)
	err = u.db.Preload("Group").
		Offset(offset).
		Limit(limit).
		Find(&users).
		Count(&count).
		Error
	return
}

func NewDBUser(db *gorm.DB) model.UserStore {
	return &dbUser{db: db}
}
