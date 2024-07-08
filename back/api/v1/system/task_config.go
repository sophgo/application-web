package system

import (
	"application-web/database"
	"application-web/logger"
	"application-web/pkg/dto"
	"application-web/pkg/handle"
	"application-web/pkg/model"
	"application-web/pkg/utils/common"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func (b *TaskApi) GetTaskConfig(c *gin.Context) {
	var req dto.TaskReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, handle.FailWithMsg(-1, "参数错误"))
		return
	}

	var task model.AlgoTaskSql
	db := database.DB.Where("task_id = ?", req.TaskId).First(&task)

	if db.Error == gorm.ErrRecordNotFound {
		logger.Error("任务不存在")
		c.JSON(http.StatusOK, handle.FailWithMsg(-1, "任务不存在"))
		return
	}

	var algorithms []dto.Algorithm
	common.ParseJSONString(task.Abilities, &algorithms)

	res := dto.AlgorithmConf{
		Device: dto.Device{
			Url:        task.Url,
			DeviceName: task.DeviceName,
			CodeName:   task.CodeName,
			Resolution: strconv.Itoa(task.VideoWidth) + "*" + strconv.Itoa(task.VideoHeight),
			Height:     task.VideoHeight,
			Width:      task.VideoWidth,
		},
		Algorithms: algorithms,
	}

	c.JSON(http.StatusOK, handle.Success(res))
}

// 修改算法任务参数配置
func (b *TaskApi) ModTaskConfig(c *gin.Context) {
	// reqBody, _ := io.ReadAll(c.Request.Body)
	// logger.Info("%s", string(reqBody))
	var req dto.TaskConf
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, handle.FailWithMsg(-1, "参数错误"))
		return
	}

	var task model.AlgoTaskSql
	db := database.DB.Where("task_id = ?", req.TaskId).First(&task)

	if db.Error == gorm.ErrRecordNotFound {
		logger.Error("任务不存在")
		c.JSON(http.StatusOK, handle.FailWithMsg(-1, "任务不存在"))
		return
	}

	var algorithms []dto.Algorithm
	common.ParseJSONString(task.Abilities, &algorithms)

	for i, algorithm := range algorithms {
		if algorithm.Type == req.Algorithm.Type {
			// algorithms[i] = req.Algorithm
			CopyAlgorithmValues(&algorithms[i], &req.Algorithm)
		}
	}

	str, _ := common.StructToString(algorithms)
	task.Abilities = str

	// 保存到数据库
	_ = database.DB.Save(&task)
	logger.Info("任务参数修改成功%s", common.StructPrint(task))

	c.JSON(http.StatusOK, handle.OkWithMsg("修改成功"))
}

func CopyAlgorithmValues(dst, src *dto.Algorithm) {

	dst.TrackInterval = src.TrackInterval
	dst.TargetSize = src.TargetSize
	dst.DetectInterval = src.DetectInterval
	dst.Threshold = src.Threshold
	dst.AlarmInterval = src.AlarmInterval

	if len(src.DetectInfos) > 0 {
		dst.DetectInfos = src.DetectInfos
	}
	if src.Extend != nil {
		dst.Extend = src.Extend
	}
}
