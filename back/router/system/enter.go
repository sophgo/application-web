package system

type RouterGroup struct {
	BaseRouter
	TaskRouter
	QueryRouter
	ReceiveAlarmRouter
	FaceTaskRouter
	FaceAlarmRouter
	FileRouter
	DashBoardRouter
}
