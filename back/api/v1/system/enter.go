package system

import "application-web/pkg/service"

type ApiGroup struct {
	BaseApi
	BasicApi
	TaskApi
	FileApi
	FaceAlarmApi
	QueryApi
	ReceiveAlarmApi
	FaceTaskApi
	DashBoardApi
}

var dashboardService = service.NewIDashboardService()
