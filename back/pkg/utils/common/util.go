package common

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// 结构体转字符串
func StructToString(data interface{}) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

// 字符串json解析到变量
func ParseJSONString(jsonStr string, target interface{}) error {
	err := json.Unmarshal([]byte(jsonStr), target)
	return err
}

// 打印结构体
func StructPrint(s interface{}) string {
	val := reflect.ValueOf(s)
	// 如果是指针，获取其所指向的元素
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// 确保给定的参数确实是结构体
	if val.Kind() != reflect.Struct {
		return "provided value is not a struct"
	}

	var result string
	// 遍历结构体的所有字段
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		typeField := val.Type().Field(i)
		fieldName := typeField.Name
		fieldValue := field.Interface()
		result += fmt.Sprintf("%s: %+v\n", fieldName, fieldValue)
	}
	return result
}

func FileIsExisted(name string) bool {
	_, err := os.Stat(name)
	if err != nil {
		return false
	}
	return true
}

// 是否包含字符串
func ContainsString(slice []string, target string) bool {
	for _, element := range slice {
		if element == target {
			return true
		}
	}
	return false
}

// 设备名称合法性检测
func ContainsSpecialCharacters(input string) bool {
	specialCharacters := "&/?=#\\+-_:;*\" "

	for _, char := range specialCharacters {
		if strings.Contains(input, string(char)) {
			return true
		}
	}

	return false
}

func ExtractIP(hostPort string) (string, string, error) {
	host, ip, err := net.SplitHostPort(hostPort)
	if err != nil {
		return "", "", err
	}
	return host, ip, nil
}

func ExtractResolution(resolution string) (int, int, error) {
	parts := strings.Split(resolution, "*")

	x, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, err
	}

	y, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, err
	}

	return x, y, nil
}

const windowWidth = 1600 // 窗口宽度
const windowHeight = 900 // 窗口高度
// 前端视频窗口大小是固定的16:9大小，需要根据视频实际分辨率进行转换
func PointCalculation(videoWidth, videoHeight, pointX, pointY int) (int, int) {
	if videoWidth <= 0 || videoHeight <= 0 {
		return -1, -1
	}

	// 计算缩放比例
	scaleX := float64(windowWidth) / float64(videoWidth)
	scaleY := float64(windowHeight) / float64(videoHeight)
	scale := scaleX
	if scaleY < scaleX {
		scale = scaleY
	}

	// 计算视频在窗口中的实际尺寸
	actualVideoWidth := int(float64(videoWidth) * scale)
	actualVideoHeight := int(float64(videoHeight) * scale)

	// 计算视频在窗口中的居中位置
	offsetX := (windowWidth - actualVideoWidth) / 2
	offsetY := (windowHeight - actualVideoHeight) / 2

	// 计算检测区域在视频中的实际位置
	detectionAreaXInVideo := int(float64(pointX-offsetX) / scale)
	detectionAreaYInVideo := int(float64(pointY-offsetY) / scale)

	// 考虑存在空白的情况
	if detectionAreaXInVideo < 0 {
		detectionAreaXInVideo = 0
	} else if detectionAreaXInVideo > videoWidth {
		detectionAreaXInVideo = videoWidth
	}

	if detectionAreaYInVideo < 0 {
		detectionAreaYInVideo = 0
	} else if detectionAreaYInVideo > videoHeight {
		detectionAreaYInVideo = videoHeight
	}

	return int(detectionAreaXInVideo), int(detectionAreaYInVideo)
}

func GenerateRandomString(length int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	b := make([]byte, length)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
