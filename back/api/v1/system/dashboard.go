package system

import (
	"application-web/pkg/handle"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DashBoardApi struct{}

func (b *DashBoardApi) LoadDashboardBaseInfo(c *gin.Context) {

	data, err := dashboardService.LoadBaseInfo()
	if err != nil {
		handle.FailWithMsg(-1, err.Error())
		return
	}
	c.JSON(http.StatusOK, handle.Success(data))

}
