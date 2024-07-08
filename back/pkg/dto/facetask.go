package dto

type FaceCreateTaskReq struct {
	TaskID   string `json:"TaskID"`
	InputSrc struct {
		SrcID     string `json:"SrcID"`
		StreamSrc struct {
			Address string `json:"Address"`
		} `json:"StreamSrc"`
	} `json:"InputSrc"`
	Algorithm FaceAlgorithm `json:"Algorithm"`
	Reporting struct {
		ReportUrlList []string `json:"ReportUrlList"`
	}
}

type FaceAlgorithm struct {
	TrackInterval    int  `json:"TrackInterval"`
	DetectInterval   int  `json:"DetectInterval"`
	AttributeInclude bool `json:"AttributeInclude"`
	FeatureInclude   bool `json:"FeatureInclude"`
	TargetSize       struct {
		MinDetect int `json:"MinDetect"`
		MaxDetect int `json:"MaxDetect"`
	} `json:"TargetSize"`
	HotRegion struct {
		LeftTopY  int `json:"LeftTopY"`
		RightBtmY int `json:"RightBtmY"`
		LeftTopX  int `json:"LeftTopX"`
		RightBtmX int `json:"RightBtmX"`
	} `json:"HotRegion"`
}

type FaceTaskReq struct {
	TaskID string `json:"TaskID"`
}

type SearchTaskReq struct {
	ImageFile        string  `json:"ImageFile"`
	ComparisonTaskID string  `json:"ComparisonTaskID"`
	SrcID            string  `json:"SrcID"`
	Threshold        float32 `json:"Threshold"`
	BeginTime        int64   `json:"BeginTime"`
	EndTime          int64   `json:"EndTime"`
	Top              int     `json:"Top"`
	ReportUrl        string  `json:"ReportUrl"`
}

type PersonInfoReq struct {
	Name      string `json:"Name"`
	ImageFile string `json:"ImageFile"`
}

type FaceFeatureReq struct {
	ImageData string `json:"ImageData"`
}

type FaceFeatureRes struct {
	Code   int                    `json:"Code"`
	Msg    string                 `json:"Msg"`
	Result []FaceFeatureResResult `json:"Result"`
}

type CommonRes struct {
	Code int    `json:"Code"`
	Msg  string `json:"Msg"`
}

type FaceFeatureResResult struct {
	FeatureData string `json:"FeatureData"`
}
type FaceTaskList struct {
	Code   int                `json:"code"`
	Msg    string             `json:"msg"`
	Result FaceTaskListResult `json:"result"`
}

type FaceTaskQueryRsp struct {
	Code   int               `json:"code"`
	Msg    string            `json:"msg"`
	Result FaceTaskStatusRsp `json:"result"`
}

type FaceTaskStatusRsp struct {
	TaskId string `json:"TaskID"`
	Status int    `json:"Status"`
}

type FaceTaskListResult struct {
	Total     int `json:"total"`
	PageSize  int `json:"pageSize"`
	PageCount int `json:"pageCount"`
	PageNo    int `json:"pageNo"`
	Items     []FaceTaskListResultItem
}

type FaceTaskListResultItem struct {
	FaceCreateTaskReq
	Status int `json:"Status"`
}

type FaceUpload struct {
	TaskId           string `json:"TaskID"`
	SceneImageBase64 string `json:"SceneImageBase64"`
	SrcID            string `json:"SrcID"`
	FrameIndex       int    `json:"FrameIndex"`
	FaceList         []FaceUploadList
}

type SearchUpload struct {
	ComparisonTaskID string       `json:"ComparisonTaskID"`
	Data             []SearchData `json:"Data"`
}

type SearchData struct {
	Id    int     `json:"id"`
	Score float32 `json:"score"`
}

type FaceQuery struct {
	BeginTime int64  `json:"beginTime"`
	EndTime   int64  `json:"endTime"`
	PageNo    int    `json:"pageNo"`
	PageSize  int    `json:"pageSize"`
	TaskID    string `json:"taskID"`
	SrcID     string `json:"srcID"`
}

type SearchQuery struct {
	PageNo           int    `json:"pageNo"`
	PageSize         int    `json:"pageSize"`
	ComparisonTaskID string `json:"ComparisonTaskID"`
}

type FaceQueryRsp struct {
	Total     int                `json:"total"`
	PageNo    int                `json:"pageNo"`
	PageSize  int                `json:"pageSize"`
	PageCount int                `json:"pageCount"`
	Items     []FaceQueryRspItem `json:"items"`
	UsedSize  string             `json:"usedSize"`
	MaxSize   string             `json:"maxSize"`
}

type FaceQueryRspItem struct {
	Id          int     `json:"id"`
	TaskId      string  `json:"taskId"`
	SrcID       string  `json:"srcId"`
	FrameIndex  int     `json:"frameIndex"`
	SmallImage  string  `json:"smallImage"`
	BigImage    string  `json:"bigImage"`
	Time        int64   `json:"time"`
	Box         FaceBox `json:"box"`
	GenderCode  int     `json:"GenderCode"`
	GlassExtend int     `json:"GlassExtend"`
}

type SearchResultQueryRspItem struct {
	FaceQueryRspItem
	Score float32 `json:"score"`
}

type SearchResultQueryRsp struct {
	Total     int                        `json:"total"`
	PageNo    int                        `json:"pageNo"`
	PageSize  int                        `json:"pageSize"`
	PageCount int                        `json:"pageCount"`
	Items     []SearchResultQueryRspItem `json:"items"`
	UsedSize  string                     `json:"usedSize"`
	MaxSize   string                     `json:"maxSize"`
}

type FaceQueryInfo struct {
	TaskIds []string `json:"taskIds"`
	SrcIds  []string `json:"srcIds"`
}

type FaceCompare struct {
	FileID []string `json:"FileID"`
}

type FileFeature struct {
	FileID string `json:"FileID"`
}

type FaceComparetest struct {
	F1 string `json:"Feature1"`
	F2 string `json:"Feature2"`
}

type FaceCompareReq struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	Result struct {
		Similarity float32 `json:"Similarity"`
	} `json:"result"`
}

type FaceCompareTestReq struct {
	Similarity float32 `json:"Similarity"`
	Len        uint    `json:"Len"`
}

type FaceUploadList struct {
	FaceBox         FaceBox `json:"FaceBox"`
	FaceImageBase64 string  `json:"FaceImageBase64"`
	FeatureBase64   string  `json:"FeatureBase64"`
	GenderCode      int     `json:"GenderCode"`
	GlassExtend     int     `json:"GlassExtend"`
}

type FaceBox struct {
	LeftTopY  int `json:"LeftTopY"`
	RightBtmY int `json:"RightBtmY"`
	LeftTopX  int `json:"LeftTopX"`
	RightBtmX int `json:"RightBtmX"`
}

func NewFaceCreateTaskReq() *FaceCreateTaskReq {

	req := &FaceCreateTaskReq{
		Algorithm: FaceAlgorithm{
			TrackInterval:    1,
			DetectInterval:   1,
			AttributeInclude: true,
			FeatureInclude:   true,
			TargetSize: struct {
				MinDetect int `json:"MinDetect"`
				MaxDetect int `json:"MaxDetect"`
			}{
				MinDetect: 0,
				MaxDetect: 250,
			},
			HotRegion: struct {
				LeftTopY  int `json:"LeftTopY"`
				RightBtmY int `json:"RightBtmY"`
				LeftTopX  int `json:"LeftTopX"`
				RightBtmX int `json:"RightBtmX"`
			}{
				LeftTopY:  200,
				RightBtmY: 400,
				LeftTopX:  1,
				RightBtmX: 3,
			},
		},
	}
	return req
}
