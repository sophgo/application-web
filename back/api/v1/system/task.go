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
	"application-web/pkg/utils/ffmpeg"
	"application-web/pkg/utils/httpclient"
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// 算法应用api
const (
	taskFind   = "/task/list"
	taskCancle = "/task/delete"
	taskSetup  = "/task/create"
	taskQuery  = "/task/query"
)

var HEADER = map[string]string{"content-type": "application/json;charset=UTF-8"}

type TaskApi struct{}

// 给前端返回算法能力列表，示例如下
//
//	"1": "吸烟检测",
//	"2": "机动车违停",
//	"3": "未戴口罩",
//	"4": "非机动车乱停",
//	"5": "火焰监测",
//	"6": "突发性事件",
//	"7": "占道经营",
//	"8": "电动车检测",
//	"9": "离岗检测",
//	"10": "施工占道"
func (b *TaskApi) AbilitiesList(c *gin.Context) {
	var res []dto.AbilityList
	for key, value := range global.System.Abilities {
		res = append(res, dto.AbilityList{
			Type: key,
			Name: value,
		})
	}

	c.JSON(http.StatusOK, handle.Success(res))
}

func (b *TaskApi) GetImage(c *gin.Context) {
	picPath := "/data/pictures/" + c.Query("taskId")
	file, err := os.Open(picPath + ".jpg")
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
func (b *TaskApi) UpdateImage(c *gin.Context) {
	req := dto.UpdateImg{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, handle.FailWithMsg(-1, "参数错误"))
		return
	}

	codeName, w, h := getStreamInfo(req.Url, global.System.Picture.Dir, req.TaskId)
	if codeName == "" || w == 0 || h == 0 {
		c.JSON(http.StatusOK, handle.FailWithMsg(-1, "视频流拉取失败，请检查地址"))
		return
	}

	file, err := os.Open(global.System.Picture.Dir + "/" + req.TaskId + ".jpg")
	if err != nil {
		c.JSON(http.StatusOK, handle.FailWithMsg(-1, "图片打开失败"))
		return
	}
	defer file.Close()

	c.JSON(http.StatusOK, handle.Ok())
}

// 添加算法任务
func (b *TaskApi) AddTask(c *gin.Context) {
	algoTask := dto.AlgoTask{}
	if err := c.ShouldBindJSON(&algoTask); err != nil {
		c.JSON(http.StatusOK, handle.FailWithMsg(-1, "参数错误"))
		return
	}
	// 参数判断
	if algoTask.DeviceName == "" || algoTask.Url == "" || algoTask.TaskId == "" {
		logger.Error("添加任务参数错误")
		c.JSON(http.StatusOK, handle.FailWithMsg(-1, "参数错误，不能为空"))
		return
	}

	if len(algoTask.Abilities) == 0 {
		logger.Error("算法能力为空")
		c.JSON(http.StatusOK, handle.FailWithMsg(-1, "算法能力为空"))
		return
	}

	if cmd.CheckIllegal(algoTask.TaskId) {
		logger.Error("任务名称不合法")
		c.JSON(http.StatusOK, handle.FailWithMsg(-1, "任务名称不能包含符号，只能为中英文和数字"))
		return
	}

	// 通过 TaskId 查找数据
	var task model.AlgoTaskSql
	result := database.DB.Where("task_id = ?", algoTask.TaskId).First(&task)

	if result.Error != gorm.ErrRecordNotFound {
		logger.Error("任务已存在")
		c.JSON(http.StatusOK, handle.FailWithMsg(-1, "任务已存在"))
		return
	}

	codeName, xMax, yMax := getStreamInfo(algoTask.Url, global.System.Picture.Dir, algoTask.TaskId)
	if codeName == "" || xMax == 0 || yMax == 0 {
		c.JSON(http.StatusOK, handle.FailWithMsg(-1, "视频流拉取失败，请检查地址"))
		return
	}

	detectInfos := []dto.DetectInfo{
		{
			HotArea: []dto.Point2D{
				{
					X: 0,
					Y: 0,
				},
				{
					X: xMax,
					Y: 0,
				},
				{
					X: xMax,
					Y: yMax,
				},
				{
					X: 0,
					Y: yMax,
				},
			},
		},
	}

	var algorithms []dto.Algorithm
	for _, algoType := range algoTask.Abilities {
		algorithm := dto.Algorithm{
			Type:           algoType,
			TrackInterval:  20,
			DetectInterval: 3,
			AlarmInterval:  300,
			Threshold:      50,
			TargetSize: dto.TargetSize{
				MinDetect: 30,
				MaxDetect: 250,
			},
			DetectInfos: detectInfos,
		}
		algorithms = append(algorithms, algorithm)
	}

	str, _ := common.StructToString(algorithms)

	algoTaskSql := model.AlgoTaskSql{
		TaskId:      algoTask.TaskId,
		Status:      0,
		DeviceName:  algoTask.DeviceName,
		Url:         algoTask.Url,
		VideoWidth:  xMax,
		VideoHeight: yMax,
		CodeName:    codeName,
		Abilities:   str,
	}

	// 保存到数据库
	_ = database.DB.Create(&algoTaskSql)

	logger.Info("任务创建成功\n%s", common.StructPrint(algoTaskSql))

	c.JSON(http.StatusOK, handle.OkWithMsg("创建成功"))
}

func (b *TaskApi) ModTask(c *gin.Context) {
	algoTask := dto.AlgoTask{}
	if err := c.ShouldBindJSON(&algoTask); err != nil {
		c.JSON(http.StatusOK, handle.FailWithMsg(-1, "参数错误"))
		return
	}

	// 通过 TaskId 查找数据
	var task model.AlgoTaskSql
	db := database.DB.Where("task_id = ?", algoTask.TaskId).First(&task)

	if db.Error == gorm.ErrRecordNotFound {
		logger.Error("任务不存在")
		c.JSON(http.StatusOK, handle.FailWithMsg(-1, "任务不存在"))
		return
	}

	var abilities []dto.Algorithm
	if err := json.Unmarshal([]byte(task.Abilities), &abilities); err != nil {
		c.JSON(http.StatusOK, handle.FailWithMsg(-1, "请求失败"+err.Error()))
		return
	}

	if algoTask.Url != task.Url {
		codeName, xMax, yMax := getStreamInfo(algoTask.Url, global.System.Picture.Dir, algoTask.TaskId)
		if codeName == "" || xMax == 0 || yMax == 0 {
			c.JSON(http.StatusOK, handle.FailWithMsg(-1, "视频流拉取失败，请检查地址"))
			return
		}
		task.VideoWidth = xMax
		task.VideoHeight = yMax
		task.CodeName = codeName
		task.Url = algoTask.Url
	}
	task.DeviceName = algoTask.DeviceName

	detectInfos := []dto.DetectInfo{
		{
			HotArea: []dto.Point2D{
				{
					X: 0,
					Y: 0,
				},
				{
					X: task.VideoWidth,
					Y: 0,
				},
				{
					X: task.VideoWidth,
					Y: task.VideoHeight,
				},
				{
					X: 0,
					Y: task.VideoHeight,
				},
			},
		},
	}

	var algorithms []dto.Algorithm
	for _, algoType := range algoTask.Abilities {
		if exist, item := hasAlgorithmWithType(abilities, algoType); exist {
			algorithms = append(algorithms, item)
			continue
		}

		algorithm := dto.Algorithm{
			Type:           algoType,
			TrackInterval:  1,
			DetectInterval: 300,
			Threshold:      0.5,
			TargetSize: dto.TargetSize{
				MinDetect: 30,
				MaxDetect: 250,
			},
			DetectInfos: detectInfos,
		}
		algorithms = append(algorithms, algorithm)
	}

	str, _ := common.StructToString(algorithms)

	task.Abilities = str
	task.Status = 0

	// 保存到数据库
	_ = database.DB.Save(&task)
	logger.Info("任务修改成功%s", common.StructPrint(task))

	stopAlgoTask(task.TaskId)

	c.JSON(http.StatusOK, handle.OkWithMsg("修改成功"))
}

func (b *TaskApi) DeleteTask(c *gin.Context) {
	var req dto.TaskReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, handle.FailWithMsg(-1, "参数错误"))
		return
	}

	if req.TaskId == "" {
		c.JSON(http.StatusOK, handle.FailWithMsg(-1, "TaskId不能为空"))
		return
	}

	// 通过 TaskName 查找数据
	var task model.AlgoTaskSql
	database.DB.Where("task_id = ?", req.TaskId).First(&task)

	if task.Status == 1 {
		stopAlgoTask(task.DeviceName)
	}

	// 删除找到的数据
	db := database.DB.Delete(&task)
	if db.Error != nil {
		logger.Error("删除数据时出错:%v", db.Error)
		return
	}
	logger.Info("任务删除成功%v", task)

	c.JSON(http.StatusOK, handle.OkWithMsg("删除成功"))
}

func (b *TaskApi) StartTask(c *gin.Context) {
	var req dto.TaskReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, handle.FailWithMsg(-1, "参数错误"))
		return
	}

	// 通过 TaskId 查找任务
	var task model.AlgoTaskSql
	db := database.DB.Where("task_id = ?", req.TaskId).First(&task)
	if db.Error == gorm.ErrRecordNotFound {
		logger.Error("任务不存在")
		c.JSON(http.StatusOK, handle.FailWithMsg(-1, "任务不存在"))
		return
	}

	if task.Status != 0 {
		c.JSON(http.StatusOK, handle.FailWithMsg(-1, "任务正在运行"))
		return
	}

	var algorithmAbilites []dto.Algorithm
	var err error
	if algorithmAbilites, err = ToAbilities(task.Abilities); err != nil {
		logger.Error("算法能力配置错误")
		c.JSON(http.StatusOK, handle.FailWithMsg(-1, "算法能力配置错误"))
		return
	}
	taskInfo := dto.Task{
		TaskId: task.TaskId,
		InputSrc: dto.InputSource{
			SrcID: task.DeviceName,
			StreamSrc: dto.StreamSource{
				Address: task.Url,
			},
		},
		Algorithm: algorithmAbilites,
		Reporting: dto.Reporting{
			ReportUrlList: []string{"http://" + global.System.UploadHost + "/api/upload"},
		},
	}

	requestBody, _ := common.StructToString(taskInfo)
	logger.Info("启动任务:%s", requestBody)

	rec := dto.AlgoReq{
		Code: -1,
	}
	data := httpclient.NewRequestWithHeaders(global.System.AlgorithmHost+taskSetup, "POST", HEADER, []byte(requestBody))
	json.Unmarshal(data, &rec)

	if rec.Code != 0 {
		if rec.Msg == "" {
			c.JSON(http.StatusOK, handle.FailWithMsg(-1, "任务启动失败"))
			return
		}
		c.JSON(http.StatusOK, handle.FailWithMsg(-1, rec.Msg))
		return
	}
	task.Status = 1
	_ = database.DB.Save(&task)
	c.JSON(http.StatusOK, handle.OkWithMsg("启动成功"))

}

// 算法任务数据库字符串解析到变量，重新计算坐标
func StrToAbilities(jsonStr string, height, width int) ([]dto.Algorithm, error) {

	var algorithmAbilites []dto.Algorithm
	err := json.Unmarshal([]byte(jsonStr), &algorithmAbilites)
	if err != nil {
		return algorithmAbilites, err
	}

	logger.Info("%v", algorithmAbilites)
	for i := range algorithmAbilites {
		for j := range algorithmAbilites[i].DetectInfos {
			tmp0 := algorithmAbilites[i].DetectInfos[j]
			if tmp0.HotArea[0].X == 0 && tmp0.HotArea[0].Y == 0 && tmp0.HotArea[1].X == width {
				continue
			}
			for k := range algorithmAbilites[i].DetectInfos[j].HotArea {
				tmp := algorithmAbilites[i].DetectInfos[j].HotArea[k]
				x, y := common.PointCalculation(width, height, tmp.X, tmp.Y)
				algorithmAbilites[i].DetectInfos[j].HotArea[k].X = x
				algorithmAbilites[i].DetectInfos[j].HotArea[k].Y = y
			}

			trw := &algorithmAbilites[i].DetectInfos[j].TripWire
			if trw.LineStart.X != 0 && trw.LineEnd.X != 0 && trw.DirectStart.X != 0 && trw.DirectEnd.X != 0 {
				trw.LineStart.X, trw.LineStart.Y = common.PointCalculation(width, height, trw.LineStart.X, trw.LineStart.Y)
				trw.LineEnd.X, trw.LineEnd.Y = common.PointCalculation(width, height, trw.LineEnd.X, trw.LineEnd.Y)
				trw.DirectStart.X, trw.DirectStart.Y = common.PointCalculation(width, height, trw.DirectStart.X, trw.DirectStart.Y)
				trw.DirectEnd.X, trw.DirectEnd.Y = common.PointCalculation(width, height, trw.DirectEnd.X, trw.DirectEnd.Y)
			}

		}
	}

	return algorithmAbilites, err
}

func ToAbilities(jsonStr string) ([]dto.Algorithm, error) {

	var algorithmAbilites []dto.Algorithm
	err := json.Unmarshal([]byte(jsonStr), &algorithmAbilites)
	if err != nil {
		return algorithmAbilites, err
	}

	logger.Info("%v", common.StructPrint(algorithmAbilites))

	return algorithmAbilites, nil
}

func (b *TaskApi) StopTask(c *gin.Context) {
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

	if !stopAlgoTask(task.TaskId) {
		c.JSON(http.StatusOK, handle.FailWithMsg(-1, "停止任务失败"))
		return
	}

	task.Status = 0
	_ = database.DB.Save(&task)

	c.JSON(http.StatusOK, handle.OkWithMsg("停止成功"))
}

func stopAlgoTask(taskId string) bool {

	body := dto.TaskRes{
		TaskId: taskId,
	}
	data, _ := json.Marshal(body)

	logger.Info("停止任务：%s", string(data))

	req := dto.AlgoReq{
		Code: -1,
	}
	res := httpclient.NewRequestWithHeaders(global.System.AlgorithmHost+taskCancle, "POST", HEADER, data)
	json.Unmarshal(res, &req)

	return req.Code == 0
}

func (b *TaskApi) List(c *gin.Context) {
	var page struct {
		PageNo   int `json:"pageNo"`
		PageSize int `json:"pageSize"`
	}
	reqBody, _ := io.ReadAll(c.Request.Body)
	json.Unmarshal(reqBody, &page)

	var algoTasks []model.AlgoTaskSql
	database.DB.Find(&algoTasks)

	var getTaskLists dto.TaskList

	res := httpclient.NewRequestWithHeaders(global.System.AlgorithmHost+taskFind, "POST", HEADER, []byte("{}"))
	err := json.Unmarshal(res, &getTaskLists)
	if err != nil {
		logger.Error("向算法服务请求  %s  失败", global.System.AlgorithmHost+taskFind)
	}

	var taskList dto.TaskListRes
	var items []dto.TaskItem
	taskList.Total = len(algoTasks)
	taskList.PageCount = taskList.Total / page.PageSize
	taskList.PageSize = page.PageSize

	if page.PageNo <= 0 {
		taskList.PageCount = 1
		taskList.PageSize = taskList.Total
		page.PageNo = 1
		page.PageSize = taskList.Total
	}
	taskList.PageCount = page.PageNo
	taskList.PageNo = page.PageNo

	for i := taskList.PageSize * (page.PageNo - 1); i < page.PageSize*page.PageNo && i < taskList.Total && i >= 0; i++ {
		status := containsTask(getTaskLists.Result, algoTasks[i].TaskId)
		// logger.Info("%s %d", algoTasks[i].TaskId, status)

		if algoTasks[i].Status != status {
			algoTasks[i].Status = status
			database.DB.Save(&algoTasks[i])
		}

		var algorithms []dto.Algorithm
		common.ParseJSONString(algoTasks[i].Abilities, &algorithms)
		var abilities []string
		var types []int

		for _, algorithm := range algorithms {
			types = append(types, algorithm.Type)
			abilities = append(abilities, global.System.Abilities[algorithm.Type])
		}
		item := dto.TaskItem{
			TaskId:      algoTasks[i].TaskId,
			DeviceName:  algoTasks[i].DeviceName,
			Status:      algoTasks[i].Status,
			Url:         algoTasks[i].Url,
			ErrorReason: "",
			Abilities:   abilities,
			Types:       types,
			Width:       algoTasks[i].VideoWidth,
			Height:      algoTasks[i].VideoHeight,
			CodeName:    algoTasks[i].CodeName,
		}
		items = append(items, item)

	}
	taskList.Items = items
	c.JSON(http.StatusOK, handle.Success(taskList))
}

func containsTask(taskInfos []dto.TaskInfo, taskId string) int {
	for _, item := range taskInfos {
		if item.TaskId == taskId {
			return item.Status
		}
	}
	return 0
}

// 判断是否存在指定 Type 的 Algorithm
func hasAlgorithmWithType(algorithms []dto.Algorithm, targetType int) (bool, dto.Algorithm) {
	for _, algo := range algorithms {
		if algo.Type == targetType {
			return true, algo
		}
	}
	return false, dto.Algorithm{}
}

func getStreamInfo(url, dir, name string) (string, int, int) {
	codeName := ""
	xMax, yMax := 0, 0
	codeName, xMax, yMax = ffmpeg.GetVediaInfo(url, dir, name)

	if codeName == "" || xMax == 0 || yMax == 0 {
		codeName, xMax, yMax = ffmpeg.HandleStream(url)
		if codeName == "" || xMax == 0 || yMax == 0 {
			return "", 0, 0
		}
		picPath := dir + "/" + name
		err1 := ffmpeg.OutPic(url, picPath)
		if err1 != nil {
			return "", 0, 0
		}
		return codeName, xMax, yMax
	}
	return codeName, xMax, yMax
}
