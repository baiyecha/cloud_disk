package middleware

import (
	"github.com/baiyecha/cloud_disk/queue"
	"github.com/gin-gonic/gin"
)

func Pub(pub queue.PubQueue) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request = c.Request.WithContext(queue.NewContext(c.Request.Context(), pub))
		c.Next()
	}
}
