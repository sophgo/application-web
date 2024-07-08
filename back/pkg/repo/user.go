package repo

import (
	"application-web/global"
	"application-web/pkg/model"
)

func QueryUserWithName(username string) (*model.User, error) {
	var user model.User

	db := global.DB.Model(&model.User{}).Where("user_name = ?", &username).First(&user)
	if err := db.Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func QueryUserWithToken(token string) (*model.User, error) {
	var user model.User
	db := global.DB.Model(&model.User{}).Where("token = ?", &token).First(&user)
	if err := db.Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func UpdateUser(user *model.User) {
	values := map[string]interface{}{
		"token":       user.Token,
		"expire_time": user.ExpireTime,
	}
	global.DB.Model(&model.User{}).Where("id = ?", user.ID).Update(values)
}

func SaveUser(user *model.User) {
	global.DB.Model(&model.User{}).Save(user)
}
