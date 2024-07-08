package system

import (
	"application-web/database"
	"application-web/global"
	"application-web/logger"
	"application-web/pkg/dto"
	"application-web/pkg/handle"
	"application-web/pkg/model"
	"application-web/pkg/utils/common"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type ReceiveAlarmApi struct{}

func (b *ReceiveAlarmApi) AlarmRev(c *gin.Context) {
	// reqBody, _ := io.ReadAll(c.Request.Body)
	// logger.Info("接受告警:%s", string(reqBody))

	// body, _ := io.ReadAll(c.Request.Body)
	// os.WriteFile("/data/test.txt", body, 0644)

	var alarmDates []dto.AlarmDate
	if err := c.ShouldBindJSON(&alarmDates); err != nil {
		logger.Error("告警接收失败，参数错误:%v", err)
		c.JSON(http.StatusOK, handle.FailWithMsg(-1, "参数错误"))
		return
	}

	// logger.Info("接受告警，%s", alarmDates[0].TaskID)

	for _, alarmDate := range alarmDates {
		dir := global.System.Picture.Dir + "/" + alarmDate.TaskID
		if !common.FileIsExisted(dir) {
			os.MkdirAll(dir, os.ModePerm)
		}

		if len(alarmDate.AnalyzeEvents) < 1 {
			logger.Error("没有目标信息")
			c.JSON(http.StatusOK, handle.FailWithMsg(-1, "参数错误"))
			return
		}

		now := time.Now()
		bigPicName := fmt.Sprintf("%s%d", common.RandStr(4), now.Unix()) + ".jpg"
		if err := jpegSave(&alarmDate.SceneImageBase64, dir+"/"+bigPicName); err != nil {
			logger.Error("图片保存错误:%v", err)
			c.JSON(http.StatusOK, handle.Fail(-1, "picture save error"))
			return
		}

		for _, event := range alarmDate.AnalyzeEvents {
			eventDir := global.System.Picture.Dir + "/" + alarmDate.TaskID + "/" + strconv.Itoa(event.Type)
			if !common.FileIsExisted(eventDir) {
				os.MkdirAll(eventDir, os.ModePerm)
			}
			smallPicName := fmt.Sprintf("%s%d", common.RandStr(4), now.Unix()) + ".jpg"
			if err := jpegSave(&event.ImageBase64, eventDir+"/"+smallPicName); err != nil {
				logger.Error("图片保存错误:%v", err)
				continue
			}

			jsonData, err := json.Marshal(event.Extend)
			if err != nil {
				logger.Error("解析Extend错误")
			}

			record := model.Record{
				TaskId:             alarmDate.TaskID,
				SrcID:              alarmDate.SrcID,
				FrameIndex:         alarmDate.FrameIndex,
				BigPictureFilename: bigPicName,

				Type:                 event.Type,
				Date:                 now.Unix(),
				SamllPictureFilename: smallPicName,
				LeftTopY:             event.Box.LeftTopY,
				RightBtmY:            event.Box.RightBtmY,
				LeftTopX:             event.Box.LeftTopX,
				RightBtmX:            event.Box.RightBtmX,
				Extend:               string(jsonData),
			}
			if err := SaveRecord(record); err != nil {
				logger.Error("写数据库错误:%v", err)
				continue
			}
		}

	}

	c.JSON(http.StatusOK, handle.Ok())
}

func jpegSave(base64ImageData *string, name string) error {
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
	options := jpeg.Options{Quality: int(global.System.Picture.Quality)}

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
func SaveRecord(record model.Record) error {
	db := database.DB.Create(&record)
	if err := db.Error; err != nil {
		return err
	}
	return nil
}
