package ffmpeg

import (
	"application-web/logger"
	"application-web/pkg/utils/cmd"
	"application-web/pkg/utils/common"
	"encoding/json"
	"fmt"
	"strings"
)

// 定义与 JSON 输出匹配的结构
type Streams struct {
	CodecName string `json:"codec_name"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
}

type VedioInfo struct {
	Data struct {
		CodecName string `json:"codeName"`
		Width     int    `json:"width"`
		Height    int    `json:"height"`
	} `json:"data"`
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func GetVediaInfo(url, dir, name string) (string, int, int) {
	comStr := fmt.Sprintf("/var/lib/application-web/get_frame -d %s -n %s -s %s -f %d", dir, name, url, 1)
	out, err := cmd.Exec(comStr)
	if err != nil {
		logger.Error("get_frame执行失败：%s", out)
		return "", 0, 0
	}

	logger.Info("%s\n %s", comStr, out)

	if out = extractStreamJSON(out); out == "" {
		logger.Error("get_frame执行失败：执行失败：%s", out)
		return "", 0, 0
	}
	logger.Info("%s", out)

	// 解析 JSON 输出
	var vedioInfo VedioInfo
	if err := json.Unmarshal([]byte(out), &vedioInfo); err != nil {
		logger.Error("解析JSON失败：%s", err.Error())
	}

	if vedioInfo.Code != 0 {
		return "", 0, 0
	}
	return vedioInfo.Data.CodecName, vedioInfo.Data.Width, vedioInfo.Data.Height
}

func HandleStream(url string) (string, int, int) {
	ffmpegBinPath := ""
	if common.IsDirectoryExists("/opt/sophon/sophon-ffmpeg-latest/bin/") {
		ffmpegBinPath = "/opt/sophon/sophon-ffmpeg-latest/bin/"
	}
	comStr := ffmpegBinPath + "ffprobe -v error  -show_entries stream=width,height,codec_name -of default=noprint_wrappers=1:nokey=1 -v quiet -of json -i " + url
	out, err := cmd.Exec(comStr)
	if err != nil {
		logger.Error("ffprobe执行失败：%s", out)
		return "", 0, 0
	}

	if out = extractJSON(out); out == "" {
		logger.Error("ffprobe执行失败：%s", out)
		return "", 0, 0
	}
	// logger.Info("%s", out)

	// 解析 JSON 输出
	var ffprobeOutput Streams
	if err := json.Unmarshal([]byte(out), &ffprobeOutput); err != nil {
		logger.Error("解析JSON失败：%s", err.Error())
	}
	return ffprobeOutput.CodecName, ffprobeOutput.Width, ffprobeOutput.Height
}

func OutPic(url, picPath string) error {
	ffmpegBinPath := ""
	if common.IsDirectoryExists("/opt/sophon/sophon-ffmpeg-latest/bin/") {
		ffmpegBinPath = "/opt/sophon/sophon-ffmpeg-latest/bin/"
	}
	comStr := ffmpegBinPath + "ffmpeg -i " + url + " -vframes 1 " + picPath + ".jpg -y "
	out, err := cmd.Exec(comStr)
	if err != nil {
		logger.Error("拉流失败：%s", out)
	}

	return err
}

func extractJSON(input string) string {
	streamsIndex := strings.Index(input, `"streams": [`)
	if streamsIndex == -1 {
		return ""
	}
	// 调整开始索引到 "streams": [ 之后
	startBracketIndex := strings.Index(input[streamsIndex:], "{")
	if startBracketIndex == -1 {
		return ""
	}
	// 计算 "{" 在原始字符串中的位置
	start := streamsIndex + startBracketIndex

	// 从 "{" 开始寻找 "}"
	endBracketIndex := strings.Index(input[start:], "}")
	if endBracketIndex == -1 {
		return ""
	}
	// 计算 "}" 在原始字符串中的位置
	end := start + endBracketIndex

	// 返回从 "{" 到 "}" 的子字符串，包括这两个字符
	return input[start : end+1]
}

func extractStreamJSON(input string) string {
	start := strings.Index(input, `{"code`)
	if start == -1 {
		return ""
	}

	// 从 "{" 开始寻找 "}"
	end := strings.Index(input, `success"}`)
	if end == -1 {
		return ""
	}
	return input[start : end+9]
}
