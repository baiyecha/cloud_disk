package handler

import (
	"net/http"

	"github.com/baiyecha/cloud_disk/errors"
	"github.com/baiyecha/cloud_disk/handler/middleware"
	"github.com/baiyecha/cloud_disk/service"
	"github.com/gin-gonic/gin"
)

type meHandler struct {
	//imageUrl image_url.URL
}

func (m *meHandler) Show(c *gin.Context) {
	uid := middleware.UserId(c)
	user, err := service.UserLoadAndRelated(c.Request.Context(), uid)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, convert2UserResp(user))
}

func (m *meHandler) UpdateInfo(c *gin.Context) {
	var authId = middleware.UserId(c)
	l := struct {
		Email      string `json:"email" form:"email"`
		Profile    string `json:"profile" form:"profile"`
		Nickname   string `json:"nickname" form:"nickname"`
		AvatarHash string `json:"avatar_hash" form:"avatar_hash"`
		Gender     int8   `json:"gender" form:"gender"`
	}{}
	if err := c.ShouldBind(&l); err != nil {
		_ = c.Error(errors.BindError(err))
		return
	}
	err := service.UserUpdate(c.Request.Context(), authId, map[string]interface{}{
		"nickname":    l.Nickname,
		"avatar_hash": l.AvatarHash,
		"profile":     l.Profile,
		"email":       l.Email,
		"gender":      l.Gender,
	})
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(204, nil)
}

func NewMeHandler() *meHandler {
	return &meHandler{}
}
