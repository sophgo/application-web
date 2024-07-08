package service

import (
	"net/http"
	"strings"
	"time"

	"application-web/pkg/model"
	"application-web/pkg/repo"

	"github.com/patrickmn/go-cache"
)

const (
	acceptLanguage = "Accept-Language"
	authorization  = "Token"
	contentType    = "Content-Type"
	multipart      = "multipart/form-data"
	Pattern        = "2006-01-02 15:04:05"
)

var tokenCache *cache.Cache

func init() {
	tokenCache = cache.New(2*time.Hour, 5*time.Minute)
}

func Token(request *http.Request) string {
	return request.Header.Get(authorization)
}

func GetUser(token string) *model.User {
	user, found := tokenCache.Get(token)
	if found {
		return user.(*model.User)
	} else {
		user, _ := repo.QueryUserWithToken(token)
		return user
	}
}

func SetUser(token string, user *model.User) {
	tokenCache.Set(token, user, 2*time.Hour)
}

func RemoveUser(token string) {
	tokenCache.Delete(token)
}

func IsMultiPartRequest(request *http.Request) bool {
	return strings.Contains(request.Header.Get(contentType), multipart)
}
