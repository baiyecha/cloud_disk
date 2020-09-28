package middleware

import (
	"github.com/baiyecha/cloud_disk/service"
	"github.com/gin-gonic/gin"
)

func Service(svc service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request = c.Request.WithContext(service.NewContext(c.Request.Context(), svc))
		c.Next()
	}
}
