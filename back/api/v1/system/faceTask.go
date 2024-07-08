package system

import (
	"application-web/database"
	"application-web/global"
	"application-web/logger"
	"application-web/pkg/dto"
	"application-web/pkg/handle"
	"application-web/pkg/model"
	"application-web/pkg/utils/cmd"
	"application-web/pkg/utils/common"
	"application-web/pkg/utils/httpclient"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
	"unsafe"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type FaceTaskApi struct{}

var TaskMgr sync.Map
var comparisonTasks []string

const (
	facetaskList      = "/task/list"
	facetaskCancle    = "/task/delete"
	facetaskSetup     = "/task/create"
	facetaskQuery     = "/task/query"
	facetaskFeature   = "/image/feature"
	faceGalleryCreate = "/gallery/create"
	searchTaskCreate  = "/search/add"
)

func FaceTaskInit() {

	offset := 0
	limit := 10
	var wg sync.WaitGroup
	for {
		var batchTask []model.FaceAlgoTaskSql
		if err := database.DB.Offset(offset).Limit(limit).Find(&batchTask).Error; err != nil {
			logger.Error("查询任务失败")
			return
		}

		if len(batchTask) == 0 {
			break
		}

		go func() {
			wg.Add(1)
			for _, task := range batchTask {
				TaskMgr.Store(task.TaskId, task.Status)
			}
			wg.Done()
		}()

		offset += limit
	}
	wg.Wait()

	go StartKeep()
}

func StartKeep() {

	ticker := time.NewTicker(time.Second * 20)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			TaskMgr.Range(func(key, value interface{}) bool {
				taskid := key.(string)
				status := value.(int)
				if status == 1 {
					if QueryTaskStatus(taskid) != 1 {
						res := startFaceTask(taskid)
						if res != 0 {
							TaskMgr.Store(taskid, 0)
							database.DB.Model(&model.FaceAlgoTaskSql{}).Where("task_id = ?", taskid).Update("status", 0)
						}
					}
				}
				return true
			})
		}
	}
}

func (b *FaceTaskApi) CreateTask(c *gin.Context) {

	req := dto.NewFaceCreateTaskReq()
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, handle.FailWithMsg(-1, "参数错误"))
		return
	}
	req.Reporting.ReportUrlList = []string{"http://" + global.System.FaceUploadHost + "/api/face/upload"}

	// 参数判断
	if req.TaskID == "" || req.InputSrc.SrcID == "" || req.InputSrc.StreamSrc.Address == "" {
		logger.Error("添加任务参数错误")
		c.JSON(http.StatusOK, handle.FailWithMsg(-1, "参数错误，TaskID、InputSrc不能为空"))
		return
	}

	if cmd.CheckIllegal(req.TaskID) {
		logger.Error("任务名称不合法")
		c.JSON(http.StatusOK, handle.FailWithMsg(-1, "任务名称不能包含符号，只能为中英文和数字"))
		return
	}

	var task model.FaceAlgoTaskSql
	result := database.DB.Where("task_id = ?", req.TaskID).First(&task)

	if result.Error != gorm.ErrRecordNotFound {
		logger.Error("任务已存在")
		c.JSON(http.StatusOK, handle.FailWithMsg(-1, "任务已存在"))
		return
	}

	codeName, xMax, yMax := getStreamInfo(req.InputSrc.StreamSrc.Address, global.System.Face.Dir, req.TaskID)
	if codeName == "" || xMax == 0 || yMax == 0 {
		c.JSON(http.StatusOK, handle.FailWithMsg(-1, "视频流拉取失败，请检查地址"))
		return
	}

	request, _ := common.StructToString(req)
	reporturls, _ := common.StructToString(req.Reporting.ReportUrlList)
	task = model.FaceAlgoTaskSql{
		TaskId:     req.TaskID,
		SrcId:      req.InputSrc.SrcID,
		SrcAddress: req.InputSrc.StreamSrc.Address,
		Request:    request,
		ReportUrls: reporturls,
		Status:     0,
	}

	_ = database.DB.Create(&task)

	TaskMgr.Store(req.TaskID, 0)
	logger.Info("任务创建成功\n%s", common.StructPrint(task))

	c.JSON(http.StatusOK, handle.OkWithMsg("创建成功"))

}

func (b *FaceTaskApi) CreatePersonInfo(c *gin.Context) {

	var req dto.PersonInfoReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, handle.FailWithMsg(-1, "参数错误"))
		return
	}

	features, status := b.GetFeatureByFile(req.ImageFile)
	if status != 0 {
		logger.Error("获取特征值失败")
		c.JSON(http.StatusOK, handle.FailWithMsg(-1, "未检测到特征值"))
		return
	}

	if len(features) != 1 {
		logger.Error("特征值个数不为1")
		c.JSON(http.StatusOK, handle.FailWithMsg(-1, "请上传仅包含一张人脸的照片"))
		return
	}

	dir, file := filepath.Split(req.ImageFile)
	ext := filepath.Ext(file)
	name := req.Name + ext
	newImagePath := filepath.Join(dir, name)

	if err := os.Rename(req.ImageFile, newImagePath); err != nil {
		logger.Error("人脸库添加失败-文件重命名失败")
		c.JSON(http.StatusOK, handle.FailWithMsg(-1, "内部错误"))
		return
	}

	p := model.PersonLibrary{
		Name:      req.Name,
		ImageFile: newImagePath,
	}

	result := database.DB.Where("name = ?", req.Name).First(&p)

	if result.Error != gorm.ErrRecordNotFound {
		logger.Error("人脸已存在")
		c.JSON(http.StatusOK, handle.FailWithMsg(-1, "人脸已存在"))
		return
	}

	_ = database.DB.Create(&p)

	logger.Info("人脸库记录创建成功\n%s", common.StructPrint(p))
	c.JSON(http.StatusOK, handle.OkWithMsg("创建成功"))

}

func (b *FaceTaskApi) RemovePersonInfo(c *gin.Context) {

	var req dto.PersonInfoReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, handle.FailWithMsg(-1, "参数错误"))
		return
	}

	var p model.PersonLibrary
	db := database.DB.Where("name = ?", req.Name).First(&p)

	if db.Error == gorm.ErrRecordNotFound {
		logger.Error("人脸不存在")
		c.JSON(http.StatusOK, handle.FailWithMsg(-1, "人脸不存在"))
		return
	}

	db = database.DB.Delete(&p)
	if db.Error != nil {
		logger.Error("删除数据时出错:%v", db.Error)
		// return
	}

	os.Remove(p.ImageFile)
	logger.Info("人脸删除成功%v", p)

	c.JSON(http.StatusOK, handle.OkWithMsg("删除成功"))
}

func (b *FaceTaskApi) CreateSearchTask(c *gin.Context) {

	var req dto.SearchTaskReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, handle.FailWithMsg(-1, "参数错误"))
		return
	}
	req.ReportUrl = "http://" + global.System.FaceUploadHost + "/api/face/upload/search"

	req.ComparisonTaskID = common.GenerateRandomString(10)
	// 参数判断
	if req.ImageFile == "" {
		logger.Error("添加任务参数错误")
		c.JSON(http.StatusOK, handle.FailWithMsg(-1, "参数错误，ImageFile不能为空"))
		return
	}

	if cmd.CheckIllegal(req.ComparisonTaskID) {
		logger.Error("任务名称不合法")
		c.JSON(http.StatusOK, handle.FailWithMsg(-1, "任务名称不能包含符号，只能为中英文和数字"))
		return
	}

	task := model.SearchTaskSql{
		SearchTaskReq: req,
	}
	result := database.DB.Where("comparison_task_id = ?", req.ComparisonTaskID).First(&task)

	if result.Error != gorm.ErrRecordNotFound {
		logger.Error("任务已存在")
		c.JSON(http.StatusOK, handle.FailWithMsg(-1, "任务已存在"))
		return
	}

	rec := dto.AlgoReq{
		Code: -1,
	}
	request, _ := common.StructToString(req)
	data := httpclient.NewRequestWithHeaders(global.System.FaceAlgoHost+searchTaskCreate, "POST", HEADER, []byte(request))
	json.Unmarshal(data, &rec)

	if rec.Code != 0 {
		logger.Error("搜图任务启动失败")
		return
	}

	_ = database.DB.Create(&task)

	logger.Info("搜图任务创建成功\n%s", common.StructPrint(task))
	comparisonTasks = append(comparisonTasks, req.ComparisonTaskID)
	c.JSON(http.StatusOK, handle.OkWithMsg("创建成功"))
}

func (b *FaceTaskApi) ListPersonInfo(c *gin.Context) {

	var personList []model.PersonLibrary

	if err := database.DB.Order("id desc").Find(&personList).Error; err != nil {
		logger.Error("查询人脸库列表失败")
		c.JSON(http.StatusOK, handle.FailWithMsg(-1, "查询人脸库列表失败"))
		return
	}

	c.JSON(http.StatusOK, handle.Success(personList))
}

func (b *FaceTaskApi) StartTask(c *gin.Context) {
	var req dto.FaceTaskReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, handle.FailWithMsg(-1, "参数错误"))
		return
	}

	startRes := startFaceTask(req.TaskID)
	switch startRes {
	case 0:
		c.JSON(http.StatusOK, handle.OkWithMsg("启动成功"))
		return
	case 1:
		c.JSON(http.StatusOK, handle.FailWithMsg(-1, "任务不存在"))
		return
	case 2:
		c.JSON(http.StatusOK, handle.FailWithMsg(-1, "任务正在运行"))
		return
	case 3:
		c.JSON(http.StatusOK, handle.FailWithMsg(-1, "算法任务配置错误"))
		return
	case 4:
		c.JSON(http.StatusOK, handle.FailWithMsg(-1, "任务启动失败"))
		return
	default:
		c.JSON(http.StatusOK, handle.FailWithMsg(-1, "未知错误"))
		return
	}
}

func startFaceTask(taskid string) int {
	var task model.FaceAlgoTaskSql
	db := database.DB.Where("task_id = ?", taskid).First(&task)

	if db.Error == gorm.ErrRecordNotFound {
		logger.Error("任务不存在")
		return 1
	}

	var request dto.FaceCreateTaskReq
	if json.Unmarshal([]byte(task.Request), &request) != nil {
		logger.Error("算法任务配置错误")
		return 3
	}

	logger.Info("启动任务:%s", task.Request)

	rec := dto.AlgoReq{
		Code: -1,
	}
	data := httpclient.NewRequestWithHeaders(global.System.FaceAlgoHost+facetaskSetup, "POST", HEADER, []byte(task.Request))
	json.Unmarshal(data, &rec)

	if rec.Code != 0 {
		logger.Error("任务启动失败")
		return 4
	}
	task.Status = 1
	_ = database.DB.Save(&task)
	TaskMgr.Store(task.TaskId, 1)

	return 0
}

func (b *FaceTaskApi) StopTask(c *gin.Context) {

	var req dto.FaceTaskReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, handle.FailWithMsg(-1, "参数错误"))
		return
	}

	if err := stopTaskByID(req.TaskID); err != 0 {
		switch err {
		case 1:
			c.JSON(http.StatusOK, handle.FailWithMsg(-1, "任务不存在"))
			return
		case 2:
			c.JSON(http.StatusOK, handle.FailWithMsg(-1, "任务停止请求构建失败"))
			return
		case 3:
			c.JSON(http.StatusOK, handle.FailWithMsg(-1, "任务停止请求发送失败"))
			return
		case 4:
			c.JSON(http.StatusOK, handle.FailWithMsg(-1, "任务停止失败"))
			return
		}
	}
	c.JSON(http.StatusOK, handle.OkWithMsg("任务停止成功"))

}

func (b *FaceTaskApi) RemoveTask(c *gin.Context) {

	var req dto.FaceTaskReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, handle.FailWithMsg(-1, "参数错误"))
		return
	}

	var task model.FaceAlgoTaskSql
	db := database.DB.Where("task_id = ?", req.TaskID).First(&task)

	taskid := task.TaskId
	if db.Error == gorm.ErrRecordNotFound {
		logger.Error("任务不存在")
		c.JSON(http.StatusOK, handle.FailWithMsg(-1, "任务不存在"))
		return
	}

	if task.Status == 1 {
		stopTask(&task)
	}

	db = database.DB.Delete(&task)
	if db.Error != nil {
		logger.Error("删除数据时出错:%v", db.Error)
		// return
	}

	b.deleteRecords(taskid)

	logger.Info("任务删除成功%v", task)

	c.JSON(http.StatusOK, handle.OkWithMsg("删除成功"))
}

func (b *FaceTaskApi) deleteRecords(taskid string) {

	var records []model.FaceRecord
	db := database.DB.Where("task_id = ?", taskid).Find(&records)

	if db.Error == gorm.ErrRecordNotFound {
		logger.Debug("任务%s不存在记录", taskid)
		return
	}

	fids := make([]uint, len(records))
	for idx, record := range records {
		fids[idx] = record.FeatureId
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
			err2 := os.Remove(bigFilePath)
			if err2 != nil {
				logger.Error("无法删除文件：%v\n", err2)
			}
		} else {
			logger.Error("无法找到文件：%v\n", err)
		}
	}

	var features []model.FaceFeature
	fdb := database.DB.Where("id IN ?", fids).Find(&features)
	if fdb.Error != gorm.ErrRecordNotFound {
		logger.Error("删除特征失败%s", fdb.Error.Error())
	}
	fdb = database.DB.Delete(&features)
	if fdb.Error != nil {
		logger.Error("删除特征失败%s", fdb.Error.Error())
	}

	db = database.DB.Delete(&records)
	if db.Error != nil {
		logger.Error("删除记录失败%s", db.Error.Error())
		return
	}

}

func (b *FaceTaskApi) ListTask(c *gin.Context) {

	var req struct {
		PageNo   int `json:"pageNo"`
		PageSize int `json:"pageSize"`
	}

	if err := c.ShouldBind(&req); err != nil {

		logger.Error("参数错误")
		c.JSON(http.StatusOK, handle.FailWithMsg(-1, "参数错误"))
		return
	}

	if req.PageNo < 1 {
		req.PageNo = 1
	}

	if req.PageSize < 1 {
		req.PageSize = 10
	}

	var taskList []model.FaceAlgoTaskSql

	if err := database.DB.Order("id desc").Offset((req.PageNo - 1) * req.PageSize).Limit(req.PageSize).Find(&taskList).Error; err != nil {
		logger.Error("查询任务列表失败")
		c.JSON(http.StatusOK, handle.FailWithMsg(-1, "查询任务列表失败"))
		return
	}

	var count int64
	database.DB.Model(&model.FaceAlgoTaskSql{}).Count(&count)

	total := int(count)/req.PageSize + 1

	var rsp dto.FaceTaskList
	rsp.Code = 0
	rsp.Msg = "查询成功"
	rsp.Result.Total = total
	rsp.Result.PageSize = req.PageSize
	rsp.Result.PageNo = req.PageNo
	for _, t := range taskList {

		if QueryTaskStatus(t.TaskId) != 1 {
			TaskMgr.Store(t.TaskId, 0)
			if t.Status == 1 {
				database.DB.Model(&model.FaceAlgoTaskSql{}).Where("task_id = ?", t.TaskId).Update("status", 0)
			}
		} else {
			TaskMgr.Store(t.TaskId, 1)
			if t.Status == 0 {
				database.DB.Model(&model.FaceAlgoTaskSql{}).Where("task_id = ?", t.TaskId).Update("status", 1)
			}
		}
		var taskInfo dto.FaceCreateTaskReq
		if err := json.Unmarshal([]byte(t.Request), &taskInfo); err != nil {
			logger.Error("任务:" + t.TaskId + " 解析失败")
			rsp.Result.Items = append(rsp.Result.Items, dto.FaceTaskListResultItem{
				FaceCreateTaskReq: dto.FaceCreateTaskReq{
					TaskID: t.TaskId,
				},
				Status: t.Status,
			})
			continue
		}
		rsp.Result.Items = append(rsp.Result.Items, dto.FaceTaskListResultItem{
			FaceCreateTaskReq: taskInfo,
			Status:            t.Status,
		})
	}

	c.JSON(http.StatusOK, rsp)
}

func (b *FaceTaskApi) ReceiveResult(c *gin.Context) {

	var req []dto.FaceUpload
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("参数错误")
		c.JSON(http.StatusOK, handle.FailWithMsg(-1, "参数错误"))
		return
	}

	for _, r := range req {
		go b.ProcessUpload(r)
	}

	c.JSON(http.StatusOK, handle.Ok())
}

func (b *FaceTaskApi) ReceiveSearch(c *gin.Context) {

	var req dto.SearchUpload
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("参数错误")
		c.JSON(http.StatusOK, handle.FailWithMsg(-1, "参数错误"))
		return
	}

	go b.ProcessSearchUpload(req)

	c.JSON(http.StatusOK, handle.Ok())
}

func (b *FaceTaskApi) GetFeature(c *gin.Context) {

	var req dto.FileFeature
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("参数错误")
		c.JSON(http.StatusOK, handle.FailWithMsg(-1, "参数错误"))
		return
	}

	features, errCode := b.GetFeatureByFile(req.FileID)
	if errCode != 0 {
		logger.Error("获取特征值失败")
		c.JSON(http.StatusOK, handle.FailWithMsg(-1, "未检测到特征值"))
		return
	}
	c.JSON(http.StatusOK, handle.Success(features))
}

func (b *FaceTaskApi) CreateGallery(c *gin.Context) {

	var req dto.FileFeature
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("参数错误")
		c.JSON(http.StatusOK, handle.FailWithMsg(-1, "参数错误"))
		return
	}

	if reqBody, err := json.Marshal(req); err != nil {
		logger.Error("请求构建失败")
		c.JSON(http.StatusOK, handle.FailWithMsg(-1, "请求构建失败"))
		return
	} else {
		data := httpclient.NewRequestWithHeaders(global.System.FaceAlgoHost+faceGalleryCreate, "POST", HEADER, reqBody)
		rec := dto.CommonRes{
			Code: -1,
		}
		if json.Unmarshal(data, &rec) != nil {
			logger.Error("请求发送失败")
			c.JSON(http.StatusOK, handle.FailWithMsg(-1, "请求发送失败"))
			return
		} else {
			if rec.Code != 0 {
				logger.Error("人脸库创建失败")
				c.JSON(http.StatusOK, handle.FailWithMsg(-1, "人脸库创建失败"))
				return
			}
		}
	}

	c.JSON(http.StatusOK, handle.Ok())
}

func (b *FaceTaskApi) Compare(c *gin.Context) {

	var req dto.FaceCompare
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("参数错误")
		c.JSON(http.StatusOK, handle.FailWithMsg(-1, "参数错误"))
		return
	}

	if len(req.FileID) != 2 {
		logger.Error("参数错误, 需要两个fileid")
		c.JSON(http.StatusOK, handle.FailWithMsg(-1, "参数错误, 需要两个fileid"))
		return
	}

	req.FileID[0] = "/data/face/upload/" + req.FileID[0]
	req.FileID[1] = "/data/face/upload/" + req.FileID[1]

	feature1, errCode1 := b.GetFeatureByFile(req.FileID[0])
	feature2, errCode2 := b.GetFeatureByFile(req.FileID[1])
	if errCode1 != 0 || errCode2 != 0 {
		logger.Error("获取特征值失败")
		c.JSON(http.StatusOK, handle.FailWithMsg(-1, "获取特征值失败"))
		return
	}

	featureBytes1, err1 := base64.StdEncoding.DecodeString(feature1[0])
	featureBytes2, err2 := base64.StdEncoding.DecodeString(feature2[0])

	if err1 != nil || err2 != nil {
		logger.Error("特征值解码失败")
		c.JSON(http.StatusOK, handle.FailWithMsg(-1, "特征值解码失败"))
		return
	}
	if len(featureBytes1) != len(featureBytes2) {
		logger.Error("特征值长度不一致")
		c.JSON(http.StatusOK, handle.FailWithMsg(-1, "特征值长度不一致"))
		return
	}

	numFloats := len(featureBytes1) / 4
	// featureFloats1 := make([]float32, numFloats)
	// var ptr *[numFloats]float32
	ptr := unsafe.Slice(&featureBytes1[0], numFloats)
	featureFloats1 := *(*[]float32)(unsafe.Pointer(&ptr))

	ptr = unsafe.Slice(&featureBytes2[0], numFloats)
	featureFloats2 := *(*[]float32)(unsafe.Pointer(&ptr))

	for i := 0; i < numFloats; i++ {
		logger.Errorln("feature" + strconv.Itoa(1) + ": " + strconv.Itoa(i) + strconv.FormatFloat(float64(featureFloats1[i]), 'f', 6, 32))
	}
	sim := dotProduct(featureFloats1, featureFloats2)

	c.JSON(http.StatusOK, handle.Success(map[string]float32{
		"Similarity": sim,
	}))

	// featureFloats1 := make([]float32, len(featureBytes1))
	// featureFloats2 := make([]float32, len(featureBytes2))
	// reader := bytes.NewReader(featureBytes1)
	// for i := 0; i < len(featureFloats1); i++ {
	// 	var f float32
	// 	if err := binary.Read(reader, binary.LittleEndian, &f); err != nil {
	// 		logger.Error("特征值1读取失败")
	// 		c.JSON(http.StatusOK, handle.FailWithMsg(-1, "特征值1读取失败"))
	// 		return
	// 	}
	// 	featureFloats1[i] = f
	// }

}

func (b *FaceTaskApi) CompareTest(c *gin.Context) {

	var req dto.FaceComparetest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("参数错误")
		c.JSON(http.StatusOK, handle.FailWithMsg(-1, "参数错误"))
		return
	}

	feature1 := req.F1
	feature2 := req.F2

	featureFloats1, numfloats1 := decodeFeature(feature1)
	featureFloats2, numfloats2 := decodeFeature(feature2)
	if numfloats1 == 0 || numfloats2 == 0 {
		c.JSON(http.StatusOK, handle.FailWithMsg(-1, "特征值解码失败"))
		return
	}
	if numfloats1 != numfloats2 {
		logger.Error("特征值长度不一致")
		c.JSON(http.StatusOK, handle.FailWithMsg(-1, "特征值长度不一致"))
		return
	}

	numFloats := numfloats2
	// featureFloats1 := make([]float32, numFloats)
	// var ptr *[numFloats]float32

	sim := dotProduct(featureFloats1, featureFloats2)
	res := dto.FaceCompareTestReq{
		Similarity: sim,
		Len:        uint(numFloats),
	}

	c.JSON(http.StatusOK, res)
}

func decodeFeature(feature string) ([]float32, int) {

	featureBytes, err := base64.StdEncoding.DecodeString(feature)
	if err != nil {
		logger.Error("特征值解码失败")
		return nil, 0
	}

	numFloats := len(featureBytes) / 4

	ptr := unsafe.Slice(&featureBytes[0], numFloats)
	featureFloats := *(*[]float32)(unsafe.Pointer(&ptr))

	return featureFloats, numFloats
}
func dotProduct(featureFloats1, featureFloats2 []float32) float32 {
	sum := float32(0)
	for i := range featureFloats1 {
		sum += featureFloats1[i] * featureFloats2[i]
	}
	return sum
}

func (b *FaceTaskApi) GetFeatureByFile(fileID string) ([]string, int) {
	img, err := os.Open(fileID)

	if err != nil {
		logger.Error("文件打开失败")
		return nil, 1
	}

	fileBytes, err1 := io.ReadAll(img)

	if err1 != nil {
		logger.Error("文件读取失败")
		return nil, 2
	}

	imgBase64 := base64.StdEncoding.EncodeToString(fileBytes)
	if len(imgBase64) == 0 {
		logger.Error("文件编码失败")
		return nil, 3
	}

	r := dto.FaceFeatureReq{
		ImageData: imgBase64,
	}

	var features []string
	if reqBody, err := json.Marshal(r); err != nil {
		logger.Error("请求构建失败")
		return nil, 4
	} else {
		data := httpclient.NewRequestWithHeaders(global.System.FaceAlgoHost+facetaskFeature, "POST", HEADER, reqBody)
		rec := dto.FaceFeatureRes{
			Code: -1,
		}
		if json.Unmarshal(data, &rec) != nil {
			logger.Error("请求发送失败")
			return nil, 5
		} else {
			if rec.Code != 0 {
				logger.Error("获取特征值失败")
				return nil, 6
			}
		}

		for _, r := range rec.Result {

			if len(r.FeatureData) == 0 {
				logger.Error("特征值解析失败")
				continue
			}

			if len(r.FeatureData)%4 != 0 {
				logger.Error("特征值长度错误")
				continue
			}

			features = append(features, r.FeatureData)
		}
	}

	if len(features) == 0 {
		logger.Error("未检测到特征值")
		return nil, 6
	}

	return features, 0

}

func (b *FaceTaskApi) ProcessUpload(req dto.FaceUpload) {

	now := time.Now()

	taskid := req.TaskId

	fileDir := global.System.Face.Dir + "/" + taskid
	_, err := os.Stat(fileDir)
	if os.IsNotExist(err) {
		os.MkdirAll(fileDir, 0755)
	}

	sceneFileName := fmt.Sprintf("%d_%d_%d_%d_scene.jpg", now.Year(), now.Month(), now.Day(), req.FrameIndex)
	if err := faceJpegSave(&req.SceneImageBase64, fileDir+"/"+sceneFileName); err != nil {
		logger.Error("场景图存储失败， frameid: %d", req.FrameIndex)
	}

	for i := range req.FaceList {

		var feature model.FaceFeature
		feature.Feature = req.FaceList[i].FeatureBase64

		_, numfloats := decodeFeature(feature.Feature)

		logger.Debug("Receive feature, length: %d", numfloats)
		if err := database.DB.Create(&feature).Error; err != nil {
			logger.Error(err.Error())
			continue
		}

		faceFileName := fmt.Sprintf("%d_%d_%d_%d_face_%d.jpg", now.Year(), now.Month(), now.Day(), req.FrameIndex, i)
		if err := faceJpegSave(&req.FaceList[i].FaceImageBase64, fileDir+"/"+faceFileName); err != nil {
			logger.Error("图片存储失败")
			continue
		}

		record := model.FaceRecord{
			TaskId:               req.TaskId,
			SrcID:                req.SrcID,
			FrameIndex:           req.FrameIndex,
			SmallPictureFilename: faceFileName,
			BigPictureFilename:   sceneFileName,
			Date:                 now.Unix(),
			FeatureId:            feature.ID,
			LeftTopX:             req.FaceList[i].FaceBox.LeftTopX,
			RightBtmY:            req.FaceList[i].FaceBox.RightBtmY,
			LeftTopY:             req.FaceList[i].FaceBox.LeftTopY,
			RightBtmX:            req.FaceList[i].FaceBox.RightBtmX,
			GenderCode:           req.FaceList[i].GenderCode,
			GlassExtend:          req.FaceList[i].GlassExtend,
		}

		if err := database.DB.Create(&record).Error; err != nil {
			logger.Error(err.Error())
			continue
		}
	}

}

func (b *FaceTaskApi) ProcessSearchUpload(req dto.SearchUpload) {

	taskid := req.ComparisonTaskID

	for _, d := range req.Data {

		record := model.SearchResultRecord{
			ComparisonTaskID: taskid,
			FaceId:           d.Id,
			Score:            d.Score,
		}
		if err := database.DB.Create(&record).Error; err != nil {
			logger.Error(err.Error())
			continue
		}
	}
}

func stopTaskByID(taskid string) int {

	var req dto.FaceTaskReq
	req.TaskID = taskid
	var task model.FaceAlgoTaskSql
	db := database.DB.Where("task_id = ?", taskid).First(&task)

	if db.Error == gorm.ErrRecordNotFound {
		logger.Error("任务不存在")
		return 1
	}

	return stopTask(&task)
}

func stopTask(task *model.FaceAlgoTaskSql) int {

	req := dto.FaceTaskReq{
		TaskID: task.TaskId,
	}
	if reqBody, err := json.Marshal(req); err != nil {
		logger.Error("请求构建失败")
		return 2
	} else {
		data := httpclient.NewRequestWithHeaders(global.System.FaceAlgoHost+facetaskCancle, "POST", HEADER, reqBody)
		rec := dto.AlgoReq{
			Code: -1,
		}
		if json.Unmarshal(data, &rec) != nil {
			logger.Error("请求发送失败")
			return 3
		} else {
			if rec.Code != 0 {
				logger.Error("任务停止失败")
				return 4
			}
		}
	}
	task.Status = 0
	TaskMgr.Store(task.TaskId, 0)
	_ = database.DB.Save(&task)
	return 0
}

func QueryTaskStatus(taskid string) int {

	rspData := httpclient.NewRequestWithHeaders(global.System.FaceAlgoHost+facetaskQuery, "POST", HEADER, []byte("{\"TaskID\":"+"\""+taskid+"\"}"))
	var rsp dto.FaceTaskQueryRsp
	if err := json.Unmarshal(rspData, &rsp); err != nil {
		logger.Error("返回解析失败")
		return 0
	}

	if rsp.Code != 0 {
		logger.Error("请求失败" + rsp.Msg)
		return 0
	}

	return 1
}

func faceJpegSave(base64ImageData *string, name string) error {
	// 将Base64数据解码成字节数组
	imageData, err := base64.StdEncoding.DecodeString(*base64ImageData)
	if err != nil {
		logger.Error("解码Base64数据失败:%v", err)
		return err
	}

	img, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		logger.Error("解码图片失败:%v", err)
		return err
	}

	// 设置 JPEG 编码器的选项，包括图像质量（1-100，100表示最高质量）,数值越高，图片越清晰，磁盘占用也越高
	options := jpeg.Options{Quality: int(global.System.Face.Quality)}

	// 图片保存
	outputFile, err := os.Create(name)
	if err != nil {
		logger.Error("创建输出文件失败:", err)
		return err
	}
	defer outputFile.Close()

	// 保存图片为JPEG格式
	err = jpeg.Encode(outputFile, img, &options)
	if err != nil {
		logger.Error("保存图片失败:", err)
		return err
	}
	return nil
}

// 保存告警
func SaveFaceRecord(record model.FaceRecord) error {
	db := database.DB.Create(&record)
	if err := db.Error; err != nil {
		return err
	}
	return nil
}
