package global

import (
	"application-web/pkg/model"

	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/gorm"
)

var (
	System           model.SystemConf
	BlockAllRequests bool

	Arch  string
	VALID *validator.Validate

	DB *gorm.DB
)
