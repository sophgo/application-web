package model

import "application-web/pkg/dto"

// 算法任务数据库存储
type FaceAlgoTaskSql struct {
	ID         uint `gorm:"primary_key"`
	TaskId     string
	SrcId      string
	SrcAddress string
	Request    string `gorm:"type:text"`
	ReportUrls string
	Status     int
}

type SearchTaskSql struct {
	dto.SearchTaskReq
}

type PersonLibrary struct {
	ID        uint `gorm:"primary_key"`
	Name      string
	ImageFile string
}

type FaceRecord struct {
	ID                   uint `gorm:"primary_key"`
	TaskId               string
	SrcID                string
	FrameIndex           int
	Date                 int64
	SmallPictureFilename string
	BigPictureFilename   string
	FeatureId            uint
	LeftTopY             int
	RightBtmY            int
	LeftTopX             int
	RightBtmX            int
	GenderCode           int
	GlassExtend          int
}

type SearchResultRecord struct {
	ID               uint `gorm:"primary_key"`
	ComparisonTaskID string
	FaceId           int
	Score            float32
}

type SearchResultFullRecord struct {
	ComparisonTaskID string
	Score            float32
	FaceRecord
}
type FaceFeature struct {
	ID      uint   `gorm:"primary_key"`
	Feature string `gorm:"type:text"`
}
