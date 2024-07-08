package service

import (
	"encoding/json"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"application-web/logger"
	"application-web/pkg/dto"
	"application-web/pkg/utils/cmd"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

type DashboardService struct{}

type IDashboardService interface {
	LoadBaseInfo() (dto.DashboardBase, error)
	LoadCurrentInfo() *dto.DashboardCurrent
}

func NewIDashboardService() IDashboardService {
	return &DashboardService{}
}

func (u *DashboardService) LoadBaseInfo() (dto.DashboardBase, error) {
	var baseInfo dto.DashboardBase
	hostInfo, err := host.Info()
	if err != nil {
		return baseInfo, err
	}
	baseInfo.Hostname = hostInfo.Hostname
	baseInfo.OS = hostInfo.OS
	baseInfo.Platform = hostInfo.Platform
	baseInfo.PlatformFamily = hostInfo.PlatformFamily
	baseInfo.PlatformVersion = hostInfo.PlatformVersion
	baseInfo.KernelArch = hostInfo.KernelArch
	baseInfo.KernelVersion = hostInfo.KernelVersion
	ss, _ := json.Marshal(hostInfo)
	baseInfo.VirtualizationSystem = string(ss)

	cpuInfo, err := cpu.Info()
	if err == nil {
		baseInfo.CPUModelName = cpuInfo[0].ModelName
	}

	baseInfo.CPUCores, _ = cpu.Counts(false)
	baseInfo.CPULogicalCores, _ = cpu.Counts(true)

	baseInfo.SDKVersion, _ = cmd.Exec("bm_version")

	baseInfo.CurrentInfo = *u.LoadCurrentInfo()
	return baseInfo, nil
}

func (u *DashboardService) LoadCurrentInfo() *dto.DashboardCurrent {
	var currentInfo dto.DashboardCurrent
	hostInfo, _ := host.Info()
	currentInfo.Uptime = hostInfo.Uptime
	currentInfo.TimeSinceUptime = time.Now().Add(-time.Duration(hostInfo.Uptime) * time.Second).Format("2006-01-02 15:04:05")
	currentInfo.Procs = hostInfo.Procs

	currentInfo.CPUTotal, _ = cpu.Counts(true)
	totalPercent, _ := cpu.Percent(0, false)
	if len(totalPercent) == 1 {
		currentInfo.CPUUsedPercent = totalPercent[0]
		currentInfo.CPUUsed = currentInfo.CPUUsedPercent * 0.01 * float64(currentInfo.CPUTotal)
	}
	currentInfo.CPUPercent, _ = cpu.Percent(0, true)

	loadInfo, _ := load.Avg()
	currentInfo.Load1 = loadInfo.Load1
	currentInfo.Load5 = loadInfo.Load5
	currentInfo.Load15 = loadInfo.Load15
	currentInfo.LoadUsagePercent = loadInfo.Load1 / (float64(currentInfo.CPUTotal)) * 100

	memoryInfo, _ := mem.VirtualMemory()
	currentInfo.MemoryTotal = memoryInfo.Total
	currentInfo.MemoryAvailable = memoryInfo.Available
	currentInfo.MemoryUsed = memoryInfo.Used
	currentInfo.MemoryUsedPercent = memoryInfo.UsedPercent

	swapInfo, _ := mem.SwapMemory()
	currentInfo.SwapMemoryTotal = swapInfo.Total
	currentInfo.SwapMemoryAvailable = swapInfo.Free
	currentInfo.SwapMemoryUsed = swapInfo.Used
	currentInfo.SwapMemoryUsedPercent = swapInfo.UsedPercent

	loadTpuMemory(&currentInfo)

	currentInfo.DiskData = loadDiskInfo()

	diskInfo, _ := disk.IOCounters()
	for _, state := range diskInfo {
		currentInfo.IOReadBytes += state.ReadBytes
		currentInfo.IOWriteBytes += state.WriteBytes
		currentInfo.IOCount += (state.ReadCount + state.WriteCount)
		currentInfo.IOReadTime += state.ReadTime
		currentInfo.IOWriteTime += state.WriteTime
	}

	netInfo, _ := net.IOCounters(false)
	if len(netInfo) != 0 {
		currentInfo.NetBytesSent = netInfo[0].BytesSent
		currentInfo.NetBytesRecv = netInfo[0].BytesRecv
	}

	currentInfo.ShotTime = time.Now()
	return &currentInfo
}

type diskInfo struct {
	Type   string
	Mount  string
	Device string
}

func loadDiskInfo() []dto.DiskInfo {
	var datas []dto.DiskInfo
	stdout, err := cmd.ExecWithTimeOut("df -hT -P|grep '/'|grep -v tmpfs|grep -v 'snap/core'|grep -v udev", 2*time.Second)
	if err != nil {
		return datas
	}
	lines := strings.Split(stdout, "\n")

	var mounts []diskInfo
	var excludes = []string{"/mnt/cdrom", "/boot", "/boot/efi", "/dev", "/dev/shm", "/run/lock", "/run", "/run/shm", "/run/user"}
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 7 {
			continue
		}
		if fields[1] == "tmpfs" {
			continue
		}
		// if strings.Contains(fields[2], "M") || strings.Contains(fields[2], "K") {
		// 	continue
		// }
		if fields[6] == "/media/root-rw" {
			continue
		}
		if strings.Contains(fields[6], "docker") {
			continue
		}

		isExclude := false
		for _, exclude := range excludes {
			if exclude == fields[6] {
				isExclude = true
			}
		}
		if isExclude {
			continue
		}
		if fields[0] == "overlay" {
			fields[0] = "/dev/mmcblk0p5"
		}
		mounts = append(mounts, diskInfo{Type: fields[1], Device: fields[0], Mount: fields[6]})
	}

	var (
		wg sync.WaitGroup
		mu sync.Mutex
	)
	wg.Add(len(mounts))
	for i := 0; i < len(mounts); i++ {
		go func(timeoutCh <-chan time.Time, mount diskInfo) {
			defer wg.Done()

			var itemData dto.DiskInfo
			itemData.Path = mount.Mount
			itemData.Type = mount.Type
			itemData.Device = mount.Device
			select {
			case <-timeoutCh:
				mu.Lock()
				datas = append(datas, itemData)
				mu.Unlock()
				logger.Error("load disk info from %s failed, err: timeout", mount.Mount)
			default:
				state, err := disk.Usage(mount.Mount)
				if err != nil {
					mu.Lock()
					datas = append(datas, itemData)
					mu.Unlock()
					logger.Error("load disk info from %s failed, err: %v", mount.Mount, err)
					return
				}
				itemData.Total = state.Total
				itemData.Free = state.Free
				itemData.Used = state.Used
				itemData.UsedPercent = state.UsedPercent
				itemData.InodesTotal = state.InodesTotal
				itemData.InodesUsed = state.InodesUsed
				itemData.InodesFree = state.InodesFree
				itemData.InodesUsedPercent = state.InodesUsedPercent
				mu.Lock()
				datas = append(datas, itemData)
				mu.Unlock()
			}
		}(time.After(5*time.Second), mounts[i])
	}
	wg.Wait()

	sort.Slice(datas, func(i, j int) bool {
		return datas[i].Path < datas[j].Path
	})
	return datas
}

func loadTpuMemory(d *dto.DashboardCurrent) {
	hasMemoryContext := false
	heapSizeRegex := regexp.MustCompile(`heap size:([0-9]+) bytes`)
	usedRegex := regexp.MustCompile(`used:([0-9]+) bytes`)
	var heapFiles []string
	memoryCtxs := []float32{0, 0, 0, 0, 0, 0}

	if _, err := os.Stat("/sys/kernel/debug/ion/cvi_npu_heap_dump/summary"); err == nil {
		heapFiles = []string{
			"/sys/kernel/debug/ion/cvi_npu_heap_dump/summary",
			"/sys/kernel/debug/ion/cvi_vpp_heap_dump/summary",
		}
	} else if _, err := os.Stat("/sys/kernel/debug/ion/bm_npu_heap_dump/summary"); err == nil {
		heapFiles = []string{
			"/sys/kernel/debug/ion/bm_npu_heap_dump/summary",
			"/sys/kernel/debug/ion/bm_vpp_heap_dump/summary",
			"/sys/kernel/debug/ion/bm_vpu_heap_dump/summary",
		}
	}

	for index, fileName := range heapFiles {
		if contex, err := os.ReadFile(fileName); err == nil {
			hasMemoryContext = true
			contex_lines := strings.Split(string(contex), "\n")
			for _, line := range contex_lines {
				if heapSizeMatches := heapSizeRegex.FindStringSubmatch(line); heapSizeMatches != nil {
					if val, err2 := strconv.Atoi(heapSizeMatches[1]); err2 == nil {
						memoryCtxs[index*2] = float32(val) / (1024 * 1024)
					}
				}
				if usedMatches := usedRegex.FindStringSubmatch(line); usedMatches != nil {
					if val, err2 := strconv.Atoi(usedMatches[1]); err2 == nil {
						memoryCtxs[index*2+1] = float32(val) / (1024 * 1024)
					}
				}
			}
		}
	}

	if hasMemoryContext {
		d.NpuMemoryTotal = uint64(memoryCtxs[0])
		d.NpuMemoryUsed = uint64(memoryCtxs[1])
		d.VppMemoryTotal = uint64(memoryCtxs[2])
		d.VppMemoryUsed = uint64(memoryCtxs[3])
		d.VpuMemoryTotal = uint64(memoryCtxs[4])
		d.VpuMemoryUsed = uint64(memoryCtxs[5])
	}

	if contex, err := os.ReadFile("/sys/class/bm-tpu/bm-tpu0/device/npu_usage"); err == nil {
		if contex[len(contex)-1] == '\n' {
			contex = contex[:len(contex)-1]
		}
		contex_lines := strings.Split(string(contex), "\n")
		util := 0
		re := regexp.MustCompile(`usage:([0-9]+)`)
		for _, line := range contex_lines {
			if match := re.FindStringSubmatch(line); match != nil {
				if usage, err2 := strconv.Atoi(match[1]); err2 == nil {
					util += usage
				}
			}
		}

		util /= len(contex_lines)
		d.TpuUsed = uint64(util)
	}
}
