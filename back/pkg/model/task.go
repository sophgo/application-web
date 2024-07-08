package model

// 告警图片数据库存储
type Record struct {
	ID                   uint `gorm:"primary_key"`
	TaskId               string
	SrcID                string
	FrameIndex           int
	Type                 int
	Date                 int64
	SamllPictureFilename string
	BigPictureFilename   string
	LeftTopY             int
	RightBtmY            int
	LeftTopX             int
	RightBtmX            int
	Extend               string `gorm:"type:text"`
	Number               int
}

// 算法任务数据库存储
type AlgoTaskSql struct {
	ID          uint   `gorm:"primary_key"`
	TaskId      string `json:"taskId"`
	Status      int    `json:"status"`
	DeviceName  string `json:"deviceName"`
	Url         string `json:"url"`
	VideoWidth  int    `json:"VideoWidth"`
	VideoHeight int    `json:"VideoHeight"`
	CodeName    string `json:"codeName"`
	Abilities   string `json:"abilities" gorm:"type:text"`
}
