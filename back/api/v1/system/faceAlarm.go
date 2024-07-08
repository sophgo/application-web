package system

import (
	"application-web/config"
	"application-web/database"
	"application-web/global"
	"application-web/logger"
	"application-web/pkg/dto"
	"application-web/pkg/handle"
	"application-web/pkg/model"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

type FaceAlarmApi struct{}

var rm_mu sync.RWMutex
var clear_mu sync.RWMutex

func (b *FaceAlarmApi) Init() {
	c := cron.New(cron.WithSeconds())
	_, err := c.AddFunc("0 0/10 * * * ?", func() {
		b.clearDisk()
	})
	if err != nil {
		fmt.Println("err:", err)
	}

	c.Start()
}

func (b *FaceAlarmApi) List(c *gin.Context) {
	var faceQuery dto.FaceQuery
	body, _ := io.ReadAll(c.Request.Body)
	json.Unmarshal(body, &faceQuery)

	records, total := b.queryRecord(faceQuery)
	var items []dto.FaceQueryRspItem

	for _, record := range records {

		smallImagePath := fmt.Sprintf("/api/face/alarm/image?taskId=%s&fileName=%s", record.TaskId, record.SmallPictureFilename)
		bigImagePath := fmt.Sprintf("/api/face/alarm/image?taskId=%s&fileName=%s", record.TaskId, record.BigPictureFilename)
		item := dto.FaceQueryRspItem{
			Id:         int(record.ID),
			TaskId:     record.TaskId,
			SrcID:      record.SrcID,
			FrameIndex: record.FrameIndex,
			BigImage:   bigImagePath,
			SmallImage: smallImagePath,
			Time:       record.Date,
			Box: dto.FaceBox{
				LeftTopY:  record.LeftTopY,
				LeftTopX:  record.LeftTopX,
				RightBtmY: record.RightBtmY,
				RightBtmX: record.RightBtmX,
			},
			GenderCode:  record.GenderCode,
			GlassExtend: record.GlassExtend,
		}

		items = append(items, item)
	}

	faceReq := dto.FaceQueryRsp{
		Total:     total,
		PageSize:  faceQuery.PageSize,
		PageNo:    faceQuery.PageNo,
		PageCount: total/faceQuery.PageSize + 1,
		Items:     items,
	}
	usize, msize := b.getSize()
	faceReq.UsedSize = formatSize(usize)
	faceReq.MaxSize = strconv.Itoa(int(msize) / 1024 / 1024)
	c.JSON(http.StatusOK, handle.Success(faceReq))
}

func (b *FaceAlarmApi) ListSearchResult(c *gin.Context) {
	var searchQuery dto.SearchQuery
	body, _ := io.ReadAll(c.Request.Body)
	json.Unmarshal(body, &searchQuery)
	if len(comparisonTasks) == 0 {
		logger.Info("no comparisonTask now")
		c.JSON(http.StatusOK, handle.Success(dto.SearchResultQueryRsp{}))
		return
	}
	searchQuery.ComparisonTaskID = comparisonTasks[len(comparisonTasks)-1]
	records, total := b.querySearchRecord(searchQuery)
	var items []dto.SearchResultQueryRspItem

	for _, record := range records {

		smallImagePath := fmt.Sprintf("/api/face/alarm/image?taskId=%s&fileName=%s", record.TaskId, record.SmallPictureFilename)
		bigImagePath := fmt.Sprintf("/api/face/alarm/image?taskId=%s&fileName=%s", record.TaskId, record.BigPictureFilename)
		item := dto.SearchResultQueryRspItem{
			FaceQueryRspItem: dto.FaceQueryRspItem{
				Id:         int(record.ID),
				TaskId:     record.TaskId,
				SrcID:      record.SrcID,
				FrameIndex: record.FrameIndex,
				BigImage:   bigImagePath,
				SmallImage: smallImagePath,
				Time:       record.Date,
				Box: dto.FaceBox{
					LeftTopY:  record.LeftTopY,
					LeftTopX:  record.LeftTopX,
					RightBtmY: record.RightBtmY,
					RightBtmX: record.RightBtmX,
				},
				GenderCode:  record.GenderCode,
				GlassExtend: record.GlassExtend,
			},
			Score: record.Score,
		}

		items = append(items, item)
	}

	faceReq := dto.SearchResultQueryRsp{
		Total:     total,
		PageSize:  searchQuery.PageSize,
		PageNo:    searchQuery.PageNo,
		PageCount: total/searchQuery.PageSize + 1,
		Items:     items,
	}
	usize, msize := b.getSize()
	faceReq.UsedSize = formatSize(usize)
	faceReq.MaxSize = strconv.Itoa(int(msize) / 1024 / 1024)
	c.JSON(http.StatusOK, handle.Success(faceReq))
}

func (b *FaceAlarmApi) clearDisk() {
	clear_mu.Lock()
	defer clear_mu.Unlock()
	num := 0

	for i := 0; i < 100; i++ {
		usedSize, maxSize := b.getSize()
		needToFree := usedSize - maxSize
		if needToFree < 0 {
			logger.Info("递归次数：%d，删除人脸记录，%d \n", i+1, num)
			return
		}
		numToDelete := int64(needToFree/(0.5*MB)) / 2
		if numToDelete == 0 {
			numToDelete = 1
		}

		res := b.removeAlarmRecord(numToDelete)
		num += res
	}
	logger.Info("删除人脸记录，%d \n", num)
}

func (b *FaceAlarmApi) ModSize(c *gin.Context) {
	var maxSize struct {
		MaxSize int64 `json:"maxSize"`
	}
	if err := c.ShouldBindJSON(&maxSize); err != nil {
		c.JSON(http.StatusOK, handle.FailWithMsg(-1, "参数错误"))
		return
	}

	global.System.Face.MaxSize = maxSize.MaxSize

	err := config.SaveConfig()
	if err != nil {
		c.JSON(http.StatusOK, handle.Fail(-1, "set error"))
		return
	}
	go clearDisk()
	c.JSON(http.StatusOK, handle.Ok())
}

func (b *FaceAlarmApi) DeleteAlarms(c *gin.Context) {
	var deleteAlarm struct {
		Number int64 `json:"number"`
	}
	if err := c.ShouldBindJSON(&deleteAlarm); err != nil || deleteAlarm.Number < 1 {
		c.JSON(http.StatusOK, handle.FailWithMsg(-1, "参数错误"))
		return
	}
	go func() {
		res := b.removeAlarmRecord(deleteAlarm.Number)
		logger.Info("删除告警记录，%d \n", res)
	}()

	c.JSON(http.StatusOK, handle.OkWithMsg("后台删除中"))
}

func (b *FaceAlarmApi) GetQueryInfo(c *gin.Context) {
	var taskIds, srcIds []string
	database.DB.Model(&model.FaceRecord{}).Pluck("distinct(task_id)", &taskIds)
	database.DB.Model(&model.FaceRecord{}).Pluck("distinct(src_id)", &srcIds)

	queryList := dto.QueryInfo{
		TaskIds: taskIds,
		SrcIds:  srcIds,
	}

	c.JSON(http.StatusOK, handle.Success(queryList))

}

func (b *FaceAlarmApi) removeAlarmRecord(num int64) int {
	rm_mu.Lock()
	defer rm_mu.Unlock()

	db := database.DB.Model(&model.Record{})
	var records []model.FaceRecord

	db.Order("date").Limit(num).Find(&records)
	if len(records) < 1 {
		return 0
	}
	// 删除记录
	for _, record := range records {
		filePath := global.System.Face.Dir + "/" + record.TaskId + "/" + record.SmallPictureFilename
		err := os.Remove(filePath)
		if err != nil {
			logger.Error("无法删除文件：%v\n", err)
		}
		// 删除记录
		if err := db.Delete(&record).Error; err != nil {
			logger.Error("无法删除记录：%v\n", err)
		}
		bigFilePath := global.System.Face.Dir + "/" + record.TaskId + "/" + record.BigPictureFilename
		if _, err := os.Stat(bigFilePath); err == nil {
			err := os.Remove(bigFilePath)
			if err != nil {
				logger.Error("无法删除文件：%v\n", err)
			}
		}
	}
	return len(records)
}

func (b *FaceAlarmApi) getSize() (int64, int64) {
	size, err := calculateSize(global.System.Face.Dir)
	if err != nil {
		logger.Error("文件大小计算错误：%v", err)
		return 0, 0
	}

	return size, global.System.Face.MaxSize * 1024 * 1024
}

func (b *FaceAlarmApi) GetImage(c *gin.Context) {
	taskId := c.Query("taskId")
	fileName := c.Query("fileName")

	var file *os.File
	var err error

	filePath := ""
	if taskId == "" {
		filePath = global.System.Face.Dir + "/" + fileName
	} else {
		filePath = global.System.Face.Dir + "/" + taskId + "/" + fileName
	}
	file, err = os.Open(filePath)

	if err != nil {
		c.String(http.StatusInternalServerError, "Unable to open image")
		return
	}
	defer file.Close()

	c.Header("Content-Type", "image/jpeg")
	_, err = io.Copy(c.Writer, file)
	if err != nil {
		c.String(http.StatusInternalServerError, "Unable to send image")
		return
	}
}

func (b *FaceAlarmApi) queryRecord(faceQuery dto.FaceQuery) ([]model.FaceRecord, int) {
	db := database.DB.Model(&model.FaceRecord{})
	var records []model.FaceRecord
	query := ""
	var vars []interface{}

	if faceQuery.BeginTime != 0 {
		query = query + " and date >= ?  "
		vars = append(vars, &faceQuery.BeginTime)
	}
	if faceQuery.EndTime != 0 {
		query = query + " and date <= ?  "
		vars = append(vars, &faceQuery.EndTime)
	}

	if faceQuery.TaskID != "" {
		query = query + " and task_id = ? "
		vars = append(vars, faceQuery.TaskID)
	}

	if faceQuery.SrcID != "" {
		query = query + " and src_id = ? "
		vars = append(vars, faceQuery.SrcID)
	}

	var total int
	if len(query) > 0 {
		db.Where(query[4:], vars...).Count(&total)
		db.Where(query[4:], vars...).Offset((faceQuery.PageNo - 1) * faceQuery.PageSize).Limit(faceQuery.PageSize).Order("date desc").Find(&records)
	} else {
		db.Count(&total)
		db.Where(query, vars...).Offset((faceQuery.PageNo - 1) * faceQuery.PageSize).Limit(faceQuery.PageSize).Order("date desc").Find(&records)
	}
	return records, total
}

func (b *FaceAlarmApi) querySearchRecord(searchQuery dto.SearchQuery) ([]model.SearchResultFullRecord, int) {
	db := database.DB.Model(&model.SearchResultRecord{}).
		Select("search_result_record.comparison_task_id, search_result_record.score, face_record.*").
		Joins("left join face_record on face_record.id = search_result_record.face_id")
	var records []model.SearchResultFullRecord
	query := "comparison_task_id = ?"

	var total int
	db.Where(query, searchQuery.ComparisonTaskID).Count(&total)
	db.Where(query, searchQuery.ComparisonTaskID).Offset((searchQuery.PageNo - 1) * searchQuery.PageSize).Limit(searchQuery.PageSize).Order("score desc").Scan(&records)
	logger.Info("len: %d, %d", len(records), total)
	return records, total
}
