package system

import (
	"net/http"
	"strings"
	"time"

	"application-web/pkg/buserr"
	"application-web/pkg/dto"
	"application-web/pkg/handle"
	"application-web/pkg/repo"
	"application-web/pkg/service"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

type BaseApi struct{}

// Login
func (b *BaseApi) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, handle.FailWithMsg(-1, "参数错误"))
		return
	}

	userName := req.UserName
	password := req.Password
	if userName == "" || password == "" {
		c.JSON(http.StatusOK, handle.Fail(buserr.InvalidUsernameOrPassword, "无效的用户名或密码"))
		return
	}

	user, _ := repo.QueryUserWithName(userName)
	var token string
	if user == nil || user.Password != password { // 验证密码
		c.JSON(http.StatusOK, handle.Fail(buserr.InvalidUsernameOrPassword, "无效的用户名或密码"))
		return
	} else {
		now := time.Now()
		if now.After(user.ExpireTime) {
			token = strings.ReplaceAll(uuid.New().String(), "-", "")
			user.Token = token
			user.LoginTime = now
			user.ExpireTime = now.Add(time.Hour * 2)
			repo.UpdateUser(user)
		} else {
			token = user.Token
		}
	}
	service.SetUser(token, user)

	c.JSON(http.StatusOK, handle.Success(dto.LoginResponse{
		Token: token,
	}))
}

func (b *BaseApi) Logout(c *gin.Context) {
	var req dto.LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, handle.FailWithMsg(-1, "参数错误"))
		return
	}

	user, _ := repo.QueryUserWithName("admin")
	user.Token = ""
	user.ExpireTime = time.Now()
	repo.UpdateUser(user)
	service.RemoveUser(req.Token)
	c.JSON(http.StatusOK, handle.Success(nil))

	// if req.Token != "" {
	// 	user, err := repo.QueryUserWithToken(req.Token)
	// 	if err == nil && user != nil {
	// 		user.Token = ""
	// 		user.ExpireTime = time.Now()
	// 		repo.UpdateUser(user)
	// 		service.RemoveUser(req.Token)
	// 		c.JSON(http.StatusOK, handle.Success(nil))

	// 	} else {
	// 		c.JSON(http.StatusOK, handle.Fail(buserr.InvalidUsernameOrPassword, "无效的token"))
	// 		return
	// 	}
	// } else {
	// 	c.JSON(http.StatusOK, handle.Fail(buserr.InvalidUsernameOrPassword, "无效的token"))
	// }
}
