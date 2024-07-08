package initialization

import (
	"application-web/api/v1/system"
	"application-web/config"
	"application-web/database"
	"application-web/global"
	"application-web/logger"
	"application-web/pkg/utils/common"
	"os"
)

func InitBase() {
	// 加载配置
	config.LoadConfig()

	// 日志处理
	logger.InitLogging(global.System.Log.Path, "algo.log", global.System.Log.Level)

	// 初始化sqlite
	database.InitDB()

	if !common.FileIsExisted(global.System.Picture.Dir) {
		os.MkdirAll(global.System.Picture.Dir, os.ModePerm)
	}

	if !common.FileIsExisted(global.System.Face.Dir) {
		os.MkdirAll(global.System.Face.Dir, os.ModePerm)
	}

	system.FaceTaskInit()
	logger.Info("%v", global.System)

}
