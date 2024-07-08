package system

import (
	"application-web/global"
	"application-web/pkg/handle"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type FileApi struct{}

func init() {
}

func (b *FileApi) Upload(c *gin.Context) {

	FileUploadPath := global.System.Face.Dir + "/upload"

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusOK, handle.FailWithMsg(-1, "文件传输错误"))
		return
	}

	newFilePath := filepath.Join(FileUploadPath, file.Filename)
	newFile, err := os.Create(newFilePath)
	if err != nil {
		c.JSON(http.StatusOK, handle.FailWithMsg(-1, "文件创建失败"))
		return
	}

	defer newFile.Close()

	if err := c.SaveUploadedFile(file, newFilePath); err != nil {
		c.JSON(http.StatusOK, handle.FailWithMsg(-1, "文件保存失败"))
		return
	}

	c.JSON(http.StatusOK, handle.Success(newFilePath))
}
