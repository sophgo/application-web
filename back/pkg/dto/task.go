package dto

// 新增算法任务接收结构体
type AlgoTask struct {
	TaskId     string `json:"taskId" validate:"required"`
	DeviceName string `json:"deviceName" validate:"required"`
	Url        string `json:"url" validate:"required"`
	Abilities  []int  `json:"types" validate:"required"`
}

type UpdateImg struct {
	TaskId string `json:"taskId" validate:"required"`
	Url    string `json:"url" validate:"required"`
}

// 任务下发参数
type Task struct {
	TaskId    string      `json:"TaskID"`
	InputSrc  InputSource `json:"InputSrc"`
	Algorithm []Algorithm `json:"Algorithm"`
	Reporting Reporting   `json:"Reporting"`
}

type InputSource struct {
	SrcID     string       `json:"SrcID"`
	StreamSrc StreamSource `json:"StreamSrc"`
}

type StreamSource struct {
	Address string `json:"Address"`
}

type Algorithm struct {
	Type           int          `json:"Type"`
	TrackInterval  int          `json:"TrackInterval"`
	DetectInterval int          `json:"DetectInterval"`
	AlarmInterval  int          `json:"AlarmInterval"`
	Threshold      float64      `json:"threshold"`
	TargetSize     TargetSize   `json:"TargetSize,omitempty"`
	DetectInfos    []DetectInfo `json:"DetectInfos,omitempty"`
	Extend         interface{}  `json:"Extend,omitempty"` // 使用 interface{} 以支持任何类型的扩展字段
}

type TargetSize struct {
	MinDetect int `json:"MinDetect"`
	MaxDetect int `json:"MaxDetect"`
}

type DetectInfo struct {
	TripWire TripWire  `json:"TripWire,omitempty"`
	HotArea  []Point2D `json:"HotArea"`
}

type TripWire struct {
	LineStart   Point2D `json:"LineStart"`
	LineEnd     Point2D `json:"LineEnd"`
	DirectStart Point2D `json:"DirectStart"`
	DirectEnd   Point2D `json:"DirectEnd"`
}

type Point2D struct {
	X int `json:"X"`
	Y int `json:"Y"`
}

type Reporting struct {
	ReportUrlList []string `json:"ReportUrlList"`
}

// 任务参数修改
type TaskConf struct {
	TaskId    string    `json:"TaskID"`
	Algorithm Algorithm `json:"Algorithm"`
}

// 查询算法服务的任务状态列表
type TaskList struct {
	Code   int        `json:"Code"`
	Msg    string     `json:"Msg"`
	Result []TaskInfo `json:"Result"`
}

type TaskInfo struct {
	TaskId string `json:"TaskID"`
	Status int    `json:"Status"`
}

// 返回给前端的任务列表
type TaskListRes struct {
	Total     int        `json:"total"`
	PageSize  int        `json:"pageSize"`
	PageCount int        `json:"pageCount"`
	PageNo    int        `json:"pageNo"`
	Items     []TaskItem `json:"items"`
}

type TaskItem struct {
	TaskId      string   `json:"taskId"`
	DeviceName  string   `json:"deviceName"`
	Url         string   `json:"url"`
	Status      int      `json:"status"`
	ErrorReason string   `json:"errorReason"`
	Abilities   []string `json:"abilities"`
	Types       []int    `json:"types"`
	Width       int      `json:"width"`
	Height      int      `json:"height"`
	CodeName    string   `json:"codeName"`
}

type AbilityList struct {
	Type int    `json:"type"`
	Name string `json:"name"`
}

// 算法配置信息，包括视频通道信息和能力列表
type AlgorithmConf struct {
	Device     Device      `json:"device"`
	Algorithms []Algorithm `json:"algorithms"`
}

// 流媒体添加设备结构体
type Device struct {
	CodeName   string `json:"codeName"`
	DeviceName string `json:"name"`
	Resolution string `json:"resolution"`
	Url        string `json:"url"`
	Width      int    `json:"width"`
	Height     int    `json:"height"`
}

// 接受前端请求内容
type TaskReq struct {
	TaskId string `json:"taskId"`
}

// 算法应用响应
type AlgoReq struct {
	Code int    `json:"Code"`
	Msg  string `json:"Msg"`
}

// 向算法应用发送请求
type TaskRes struct {
	TaskId string `json:"TaskID"`
}

// 接受告警信息
type AlarmDate struct {
	TaskID           string         `json:"TaskID" validate:"required"`
	SceneImageBase64 string         `json:"SceneImageBase64" validate:"required"`
	SrcID            string         `json:"SrcID"`
	FrameIndex       int            `json:"FrameIndex"`
	AnalyzeEvents    []AnalyzeEvent `json:"AnalyzeEvents" validate:"required"`
}

type AnalyzeEvent struct {
	ImageBase64 string      `json:"ImageBase64"`
	Type        int         `json:"Type"`
	Box         Box         `json:"Box"`
	Extend      interface{} `json:"Extend"`
}

type Box struct {
	LeftTopY  int `json:"LeftTopY"`  //左上⻆ y 坐标
	RightBtmY int `json:"RightBtmY"` //右下⻆ y 坐标
	LeftTopX  int `json:"LeftTopX"`  //左上⻆ x 坐标
	RightBtmX int `json:"RightBtmX"` //右下⻆ x 坐标
}

// 告警信息查询请求
type AlgoQuery struct {
	BeginTime int64  `json:"beginTime"`
	EndTime   int64  `json:"endTime"`
	PageNo    int    `json:"pageNo"`
	PageSize  int    `json:"pageSize"`
	TaskID    string `json:"taskID"`
	SrcID     string `json:"srcID"`
	Types     []int  `json:"types"`
}

// 告警信息查询响应
type AlarmReq struct {
	Total     int            `json:"total"`
	PageNo    int            `json:"pageNo"`
	PageSize  int            `json:"pageSize"`
	PageCount int            `json:"pageCount"`
	Items     []AlarmReqItem `json:"items"`
	UsedSize  string         `json:"usedSize"`
	MaxSize   string         `json:"maxSize"`
}

type AlarmReqItem struct {
	Id         int    `json:"id"`
	TaskId     string `json:"taskId"`
	SrcID      string `json:"srcId"`
	FrameIndex int    `json:"frameIndex"`
	Type       int    `json:"type"`
	SmallImage string `json:"smallImage"`
	BigImage   string `json:"bigImage"`
	Time       int64  `json:"time"`
	Box        Box    `json:"box"`
	Extend     string `json:"Extend"`
}

type QueryInfo struct {
	Types   []int    `json:"types"`
	TaskIds []string `json:"taskIds"`
	SrcIds  []string `json:"srcIds"`
}
