package system

import (
	"application-web/config"
	"application-web/database"
	"application-web/global"
	"application-web/logger"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"application-web/pkg/dto"
	"application-web/pkg/handle"
	"application-web/pkg/model"
	"strconv"
	"sync"

	"strings"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

const (
	MB = 1 * 1024 * 1024    // 1MB
	GB = 1024 * 1024 * 1024 // 1GB
)

type QueryApi struct{}

func init() {
	c := cron.New(cron.WithSeconds())
	_, err := c.AddFunc("0 0/10 * * * ?", func() {
		clearDisk()
	})
	if err != nil {
		fmt.Println("err:", err)
	}

	c.Start()
}

func (b *QueryApi) GetRecord(c *gin.Context) {
	var algoQuery dto.AlgoQuery
	body, _ := io.ReadAll(c.Request.Body)
	json.Unmarshal(body, &algoQuery)

	records, total := queryRecord(algoQuery)
	var items []dto.AlarmReqItem

	for _, record := range records {

		smallImagePath := fmt.Sprintf("/api/alarm/image?taskId=%s&type=%d&fileName=%s", record.TaskId, record.Type, record.SamllPictureFilename)
		bigImagePath := fmt.Sprintf("/api/alarm/image?taskId=%s&fileName=%s", record.TaskId, record.BigPictureFilename)
		item := dto.AlarmReqItem{
			Id:         int(record.ID),
			TaskId:     record.TaskId,
			SrcID:      record.SrcID,
			FrameIndex: record.FrameIndex,
			Type:       record.Type,
			BigImage:   bigImagePath,
			SmallImage: smallImagePath,
			Time:       record.Date,
			Box: dto.Box{
				LeftTopY:  record.LeftTopY,
				LeftTopX:  record.LeftTopX,
				RightBtmY: record.RightBtmY,
				RightBtmX: record.RightBtmX,
			},
			Extend: record.Extend,
		}
		items = append(items, item)
	}

	alarmReq := dto.AlarmReq{
		Total:     total,
		PageSize:  algoQuery.PageSize,
		PageNo:    algoQuery.PageNo,
		PageCount: total/algoQuery.PageSize + 1,
		Items:     items,
	}
	usize, msize := getSize()
	alarmReq.UsedSize = formatSize(usize)
	alarmReq.MaxSize = strconv.Itoa(int(msize) / 1024 / 1024)
	c.JSON(http.StatusOK, handle.Success(alarmReq))

}

func (b *QueryApi) GetImage(c *gin.Context) {
	taskId := c.Query("taskId")
	event := c.Query("type")
	fileName := c.Query("fileName")

	var file *os.File
	var err error
	// 打开图片文件
	if event == "" {
		file, err = os.Open(global.System.Picture.Dir + "/" + taskId + "/" + fileName)
	} else {
		file, err = os.Open(global.System.Picture.Dir + "/" + taskId + "/" + event + "/" + fileName)
	}

	if err != nil {
		c.String(http.StatusInternalServerError, "Unable to open image")
		return
	}
	defer file.Close()

	// 设置响应头部信息
	c.Header("Content-Type", "image/jpeg")
	_, err = io.Copy(c.Writer, file)
	if err != nil {
		c.String(http.StatusInternalServerError, "Unable to send image")
		return
	}
}

func (b *QueryApi) GetQueryInfo(c *gin.Context) {
	var taskIds, srcIds []string
	var types []int
	database.DB.Model(&model.Record{}).Pluck("distinct(task_id)", &taskIds)
	database.DB.Model(&model.Record{}).Pluck("distinct(src_id)", &srcIds)
	database.DB.Model(&model.Record{}).Pluck("distinct(type)", &types)

	queryList := dto.QueryInfo{
		Types:   types,
		TaskIds: taskIds,
		SrcIds:  srcIds,
	}

	c.JSON(http.StatusOK, handle.Success(queryList))

}

func (b *QueryApi) ModSize(c *gin.Context) {
	var maxSize struct {
		MaxSize int64 `json:"maxSize"`
	}
	if err := c.ShouldBindJSON(&maxSize); err != nil {
		c.JSON(http.StatusOK, handle.FailWithMsg(-1, "参数错误"))
		return
	}

	global.System.Picture.MaxSize = maxSize.MaxSize

	err := config.SaveConfig()
	if err != nil {
		c.JSON(http.StatusOK, handle.Fail(-1, "set error"))
		return
	}
	go clearDisk()
	c.JSON(http.StatusOK, handle.Ok())
}

func (b *QueryApi) DeleteAlarms(c *gin.Context) {
	var deleteAlarm struct {
		Number int64 `json:"number"`
	}
	if err := c.ShouldBindJSON(&deleteAlarm); err != nil || deleteAlarm.Number < 1 {
		c.JSON(http.StatusOK, handle.FailWithMsg(-1, "参数错误"))
		return
	}
	go func() {
		res := removeAlarmRecode(deleteAlarm.Number) // 执行删除操作
		logger.Info("删除告警记录，%d \n", res)
	}()

	c.JSON(http.StatusOK, handle.OkWithMsg("后台删除中"))
}

func queryRecord(algoQuery dto.AlgoQuery) ([]model.Record, int) {
	db := database.DB.Model(&model.Record{})
	var records []model.Record
	// 构建基本的查询语句
	query := ""
	// 构建条件语句
	var vars []interface{}

	// 添加 BeginTime 和 EndTime 条件
	if algoQuery.BeginTime != 0 {
		query = query + " and date >= ?  "
		vars = append(vars, &algoQuery.BeginTime)
	}
	if algoQuery.EndTime != 0 {
		query = query + " and date <= ?  "
		vars = append(vars, &algoQuery.EndTime)
	}

	// 添加 task_id 条件
	if algoQuery.TaskID != "" {
		query = query + " and task_id = ? "
		vars = append(vars, algoQuery.TaskID)
	}

	// 添加 task_id 条件
	if algoQuery.SrcID != "" {
		query = query + " and src_id = ? "
		vars = append(vars, algoQuery.SrcID)
	}

	// 添加 Alarms 条件
	if len(algoQuery.Types) > 0 {
		alarmConditions := []string{}
		for _, alarm := range algoQuery.Types {
			alarmConditions = append(alarmConditions, "type = ?")
			vars = append(vars, alarm)
		}
		query = query + " and (" + strings.Join(alarmConditions, " or ") + ")"
	}

	var total int
	if len(query) > 0 {
		db.Where(query[4:], vars...).Count(&total)
		db.Where(query[4:], vars...).Offset((algoQuery.PageNo - 1) * algoQuery.PageSize).Limit(algoQuery.PageSize).Order("date desc").Find(&records)
	} else {
		db.Count(&total)
		db.Where(query, vars...).Offset((algoQuery.PageNo - 1) * algoQuery.PageSize).Limit(algoQuery.PageSize).Order("date desc").Find(&records)
	}
	return records, total
}

func getSize() (int64, int64) {
	size, err := calculateSize(global.System.Picture.Dir)
	if err != nil {
		logger.Error("文件大小计算错误：%v", err)
		return 0, 0
	}

	return size, global.System.Picture.MaxSize * 1024 * 1024
}

func formatSize(size int64) string {
	const (
		B  = 1
		KB = 1024 * B
		MB = 1024 * KB
		GB = 1024 * MB
		TB = 1024 * GB
	)

	switch {
	case size >= TB:
		return fmt.Sprintf("%.2f TB", float64(size)/float64(TB))
	case size >= GB:
		return fmt.Sprintf("%.2f GB", float64(size)/float64(GB))
	case size >= MB:
		return fmt.Sprintf("%.2f MB", float64(size)/float64(MB))
	case size >= KB:
		return fmt.Sprintf("%.2f KB", float64(size)/float64(KB))
	default:
		return fmt.Sprintf("%d B", size)
	}
}

func calculateSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})

	return size, err
}

var mu1 sync.RWMutex

func removeAlarmRecode(num int64) int {
	mu1.Lock()         // 加锁
	defer mu1.Unlock() // 在函数退出时解锁

	db := database.DB.Model(&model.Record{})
	var records []model.Record

	// 查询最旧的num条记录
	db.Order("date").Limit(num).Find(&records)
	if len(records) < 1 {
		return 0
	}
	// 删除记录
	for _, record := range records {
		filePath := global.System.Picture.Dir + "/" + record.TaskId + "/" + strconv.Itoa(record.Type) + "/" + record.SamllPictureFilename
		err := os.Remove(filePath)
		if err != nil {
			logger.Error("无法删除文件：%v\n", err)
		}
		// 删除记录
		if err := db.Delete(&record).Error; err != nil {
			logger.Error("无法删除记录：%v\n", err)
		}
		if isNotInRecodes(record.BigPictureFilename) {
			filePath := global.System.Picture.Dir + "/" + record.TaskId + "/" + record.BigPictureFilename
			err := os.Remove(filePath)
			if err != nil {
				logger.Error("无法删除文件：%v\n", err)
			}
		}
	}
	return len(records)
}

func isNotInRecodes(name string) bool {
	var record model.Record
	notFound := database.DB.Where("big_picture_filename = ?", name).First(&record).RecordNotFound()
	return notFound
}

var mu sync.RWMutex

func clearDisk() {
	mu.Lock()         // 加锁
	defer mu.Unlock() // 在函数退出时解锁
	num := 0

	for i := 0; i < 100; i++ {
		usedSize, maxSize := getSize()
		needToFree := usedSize - maxSize
		if needToFree < 0 {
			logger.Info("递归次数：%d，删除告警记录，%d \n", i+1, num)
			return
		}
		numToDelete := int64(needToFree/(0.5*MB)) / 2
		if numToDelete == 0 {
			numToDelete = 1 // 至少删除一个记录
		}

		res := removeAlarmRecode(numToDelete) // 执行删除操作
		num += res
	}
	logger.Info("删除告警记录，%d \n", num)
}
