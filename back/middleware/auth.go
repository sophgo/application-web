package middleware

import (
	"application-web/pkg/repo"
	"application-web/pkg/service"
	"net/http"

	"time"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		_token := service.Token(c.Request)
		if _token != "" {
			user := service.GetUser(_token)
			if user != nil {
				now := time.Now()
				if now.Before(user.ExpireTime) {
					if user.ExpireTime.After(user.ExpireTime.Add(time.Minute * 10)) {
						user.ExpireTime = now.Add(time.Hour * 2)
						repo.UpdateUser(user)
					}
				}
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusUnauthorized, "无效的token")
	}
}
