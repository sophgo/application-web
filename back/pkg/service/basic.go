package service

import (
	"application-web/logger"
	"application-web/pkg/utils/cmd"
	"bufio"
	"fmt"
	"os"
	"strings"
)

type BasicService struct {
}

func (b *BasicService) ReplaceNTPValue(value string) error {
	input, err := os.Open("/etc/systemd/timesyncd.conf")
	if err != nil {
		return err
	}
	defer input.Close()

	// 创建一个临时文件
	output, err := os.CreateTemp("", "timesyncd.conf")
	if err != nil {
		return err
	}
	defer output.Close()

	scanner := bufio.NewScanner(input)
	found := false

	for scanner.Scan() {
		line := scanner.Text()
		// 检查是否为 NTP= 行
		if strings.HasPrefix(line, "NTP=") || strings.HasPrefix(line, "#NTP=") {
			found = true
			// 替换行内容
			line = "NTP=" + value
		}
		// 写入新行到临时文件
		if _, err := output.WriteString(line + "\n"); err != nil {
			return err
		}
	}

	if !found {
		logger.Error("没有找到 NTP= 行")
	}

	// 检查扫描是否有错误
	if err := scanner.Err(); err != nil {
		return err
	}

	// 关闭输出文件以确保所有数据都已写入
	if err := output.Close(); err != nil {
		return err
	}

	// 替换原始文件
	return os.Rename(output.Name(), "/etc/systemd/timesyncd.conf")
}

func (b *BasicService) GetNTPValue() (string, error) {
	file, err := os.Open("/etc/systemd/timesyncd.conf")
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "NTP=") {
			return strings.TrimPrefix(line, "NTP="), nil
		}
		if strings.HasPrefix(line, "#NTP=") {
			return strings.TrimPrefix(line, "#NTP="), nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	// 如果文件中没有找到 NTP= 行
	return "", fmt.Errorf("配置文件/etc/systemd/timesyncd.conf 错误")
}

// SE6和SE5 ntp配置不一样
func (b *BasicService) GetSE6NTPValue() (string, error) {
	ntpPool := `awk '/^pool/ {
		for (i=2; i<NF; i++) {
			print $i
		}
	}' /etc/ntp.conf`
	out, err := cmd.Exec(ntpPool)
	if err != nil {
		logger.Error("out %s, err:%v ", out, err)
		return "", err
	}
	out = strings.ReplaceAll(out, "\n", " ")
	out = strings.ReplaceAll(out, "\r\n", " ")
	logger.Info("get /etc/ntp.conf %s ", out)

	return out, err
}

func (b *BasicService) SE6Ntpd(hosts string) error {
	deleteAll := "sed -i '/^pool/d' /etc/ntp.conf"
	out, err := cmd.Exec(deleteAll)
	if err != nil {
		return err
	}
	logger.Info("delete /etc/ntp.conf %s ", out)

	appendPool := "sed -i '$a pool %s iburst' /etc/ntp.conf"

	addHost := fmt.Sprintf(appendPool, hosts)
	out, err = cmd.Exec(addHost)
	if err != nil {
		return err
	}
	logger.Info("add /etc/ntp.conf %s ", out)

	restart := "systemctl restart ntp"
	out, err = cmd.Exec(restart)
	if err != nil {
		return err
	}
	logger.Info("restart /etc/ntp.conf %s ", out)

	return nil
}
