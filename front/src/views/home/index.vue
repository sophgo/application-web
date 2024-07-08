<template>
  <el-container class="system-status">
    <el-header>
      <h1>系统状态</h1>
    </el-header>
    <el-main>
      <el-row :gutter="20">
        <el-col :span="12">
          <el-card class="box-card">
            <div slot="header" class="clearfix">
              <span>CPU使用率</span>
            </div>
            <el-progress :percentage="+systemBaord?.result.currentInfo.cpuUsedPercent.toFixed(2)" :color="customColor(systemBaord?.result.currentInfo.cpuUsedPercent)"></el-progress>
            <br />
            <el-popover trigger="hover" placement="top" width="200">
              <div>
                <div v-for="(cpuPercent, index) in systemBaord?.result.currentInfo.cpuPercent" :key="index">CPU {{ index + 1 }}: {{ cpuPercent.toFixed(2) }}%</div>
              </div>
            </el-popover>
            <el-popover>
              <template #reference>
                <span style="background-color: rgb(194, 235, 249); padding: 5px">{{ systemBaord?.result.cpuCores }}核 - {{ systemBaord?.result.cpuModelName }}</span>
              </template>
              <template #default>
                <div style="background-color: rgb(180, 222, 246); margin: 6px" v-for="(cpuPercent, index) in systemBaord?.result.currentInfo.cpuPercent" :key="index">
                  CPU {{ index + 1 }}: {{ cpuPercent.toFixed(2) }}%
                </div>
              </template>
            </el-popover>
          </el-card>
        </el-col>
        <el-col :span="12">
          <el-card class="box-card">
            <div slot="header" class="clearfix">
              <span>系统内存使用率</span>
            </div>
            <el-progress :percentage="+systemBaord?.result.currentInfo.memoryUsedPercent.toFixed(2)" :color="customColor(systemBaord?.result.currentInfo.memoryUsedPercent)"></el-progress>
            <br />
            <span style="background-color: rgb(241, 215, 146); padding: 5px">
              {{ formatDiskCapacity(systemBaord?.result.currentInfo.memoryUsed) }} / {{ formatDiskCapacity(systemBaord?.result.currentInfo.memoryTotal) }}
            </span>
          </el-card>
        </el-col>
        <el-col :span="12">
          <el-card class="box-card">
            <div slot="header" class="clearfix">
              <span>系统负载</span>
            </div>
            <el-progress :percentage="+(systemBaord?.result.currentInfo.loadUsagePercent / 2).toFixed(2)" :color="customColor(systemBaord?.result.currentInfo.loadUsagePercent / 2)"></el-progress>
            <br />
            <el-popover trigger="hover" placement="top">
              <div></div>
            </el-popover>
            <el-popover>
              <template #reference>
                <span style="background-color: rgb(226, 240, 231); padding: 5px">{{ loadStatus(systemBaord?.result.currentInfo.loadUsagePercent / 2) }}</span>
              </template>
              <template #default>
                <div>最近 1 分钟平均负载: {{ systemBaord?.result.currentInfo.load1.toFixed(2) }}</div>
                <div>最近 5 分钟平均负载: {{ systemBaord?.result.currentInfo.load5.toFixed(2) }}</div>
                <div>最近 15 分钟平均负载: {{ systemBaord?.result.currentInfo.load15.toFixed(2) }}</div>
              </template>
            </el-popover>
          </el-card>
        </el-col>
        <el-col :span="12">
          <el-card class="box-card">
            <div slot="header" class="clearfix">
              <span>TPU使用率</span>
            </div>
            <el-progress :percentage="+systemBaord?.result.currentInfo.tpuUsed" :color="customColor(systemBaord?.result.currentInfo.tpuUsed)"></el-progress>
            <br />
            <span style="padding: 5px"></span>
          </el-card>
        </el-col>
      </el-row>
      <el-row :gutter="20">
        <el-col :span="24">
          <el-card class="box-card">
            <div slot="header" class="clearfix">
              <span>TPU设备内存</span>
            </div>
            <br />
            <el-progress
              :percentage="+((100 * systemBaord?.result.currentInfo.npuMemoryUsed) / systemBaord?.result.currentInfo.npuMemoryTotal).toFixed(2)"
              :color="customColor((100 * systemBaord?.result.currentInfo.npuMemoryUsed) / systemBaord?.result.currentInfo.npuMemoryTotal)"
            ></el-progress>
            <span style="padding: 5px">
              NPU内存： {{ formatDiskCapacity(systemBaord?.result.currentInfo.npuMemoryUsed * 1024 * 1024) }} / {{ formatDiskCapacity(systemBaord?.result.currentInfo.npuMemoryTotal * 1024 * 1024) }}
            </span>
            <br />
            <br />
            <el-progress
              :percentage="+((100 * systemBaord?.result.currentInfo.vppMemoryUsed) / systemBaord?.result.currentInfo.vppemoryTotal).toFixed(2)"
              :color="customColor((100 * systemBaord?.result.currentInfo.vppMemoryUsed) / systemBaord?.result.currentInfo.vppemoryTotal)"
            ></el-progress>
            <span style="padding: 5px">
              VPP内存： {{ formatDiskCapacity(systemBaord?.result.currentInfo.vppMemoryUsed * 1024 * 1024) }} / {{ formatDiskCapacity(systemBaord?.result.currentInfo.vppemoryTotal * 1024 * 1024) }}
            </span>
            <br />
            <br />
            <el-progress
              v-if="systemBaord?.result.currentInfo.vpuMemoryTotal !== 0"
              :percentage="+((100 * systemBaord?.result.currentInfo.vpuMemoryUsed) / systemBaord?.result.currentInfo.vpuMemoryTotal).toFixed(2)"
              :color="customColor((100 * systemBaord?.result.currentInfo.vpuMemoryUsed) / systemBaord?.result.currentInfo.vpuMemoryTotal)"
            ></el-progress>
            <span v-if="systemBaord?.result.currentInfo.vpuMemoryTotal !== 0" style="padding: 5px">
              VPU内存： {{ formatDiskCapacity(systemBaord?.result.currentInfo.vpuMemoryUsed * 1024 * 1024) }} / {{ formatDiskCapacity(systemBaord?.result.currentInfo.vpuMemoryTotal * 1024 * 1024) }}
            </span>
          </el-card>
        </el-col>
      </el-row>
      <el-row :gutter="20">
        <el-col :span="24">
          <el-card class="box-card">
            <div slot="header" class="clearfix">
              <span>磁盘使用率</span>
            </div>
            <el-table :data="systemBaord?.result.currentInfo.diskData" style="width: 100%">
              <el-table-column prop="path" label="挂载目录"></el-table-column>
              <el-table-column prop="device" label="磁盘"></el-table-column>
              <el-table-column prop="容量" label="容量">
                <template #default="{ row }">
                  <span>{{ formatDiskCapacity(row.used) }} / {{ formatDiskCapacity(row.total) }}</span>
                </template>
              </el-table-column>
              <el-table-column prop="使用率" label="使用率 (%)">
                <template #default="{ row }">
                  <el-progress :percentage="+row.usedPercent.toFixed(2)" color="rgb(103, 194, 58)"></el-progress>
                </template>
              </el-table-column>
            </el-table>
          </el-card>
          <el-card class="box-card system-info-card">
            <div slot="header" class="clearfix card-header">
              <span>系统信息</span>
            </div>
            <div class="card-content">
              <p>
                <strong>主机名：</strong>
                {{ systemBaord?.result.hostname }}
              </p>
              <p>
                <strong>操作系统：</strong>
                {{ systemBaord?.result.platform }}:{{ systemBaord?.result.platformVersion }}
              </p>
              <p>
                <strong>SDK版本：</strong>
                <span class="sdk-info" v-html="formatSDKInfo(systemBaord?.result.SDKVersion)"></span>
              </p>
              <p>
                <strong>内核版本：</strong>
                {{ systemBaord?.result.kernelVersion }}
              </p>
              <p>
                <strong>系统时间：</strong>
                {{ timeFormat(systemBaord?.result.currentInfo.shotTime) }}
              </p>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </el-main>
  </el-container>
</template>

<script setup lang="ts">
import { onMounted, ref, onBeforeUnmount } from 'vue'
import type { SystemInfo } from '@/api/home/type'
import { reqSystemInfo } from '@/api/home'
const customColor = (percentage: number) => [
  { color: '#1989fa', percentage: 20 },
  { color: '#5cb87a', percentage: 40 },
  { color: '#e6a23c', percentage: 60 },
  { color: '#6f7ad3', percentage: 80 },
  { color: '#f56c6c', percentage: 100 },
]
const systemBaord = ref<SystemInfo>({
  code: -1,
  msg: '',
  result: {
    hostname: '',
    os: '',
    platform: '',
    platformFamily: '',
    platformVersion: '',
    kernelArch: '',
    kernelVersion: '',
    virtualizationSystem: '',
    cpuCores: 4,
    cpuLogicalCores: 8,
    cpuModelName: '',
    SDKVersion: '',
    currentInfo: {
      uptime: 3600,
      timeSinceUptime: '',
      procs: 100,
      load1: 0.5,
      load5: 0.8,
      load15: 1.2,
      loadUsagePercent: 80,
      cpuPercent: [10, 20, 30],
      cpuUsedPercent: 40,
      cpuUsed: 4,
      cpuTotal: 8,
      memoryTotal: 8192,
      memoryAvailable: 4096,
      memoryUsed: 4096,
      memoryUsedPercent: 50,
      swapMemoryTotal: 4096,
      swapMemoryAvailable: 2048,
      swapMemoryUsed: 2048,
      swapMemoryUsedPercent: 50,
      npuMemoryTotal: 0,
      npuMemoryUsed: 0,
      vppemoryTotal: 0,
      vppMemoryUsed: 0,
      vpuMemoryTotal: 0,
      vpuMemoryUsed: 0,
      tpuUsed: 0,
      ioReadBytes: 1024,
      ioWriteBytes: 512,
      ioCount: 1000,
      ioReadTime: 100,
      ioWriteTime: 50,
      diskData: [
        {
          path: '',
          type: 'ext4',
          device: '',
          total: 0,
          free: 0,
          used: 0,
          usedPercent: 0,
          inodesTotal: 100000,
          inodesUsed: 80000,
          inodesFree: 20000,
          inodesUsedPercent: 80,
        },
      ],
      netBytesSent: 1024,
      netBytesRecv: 512,
      shotTime: '',
    },
  },
})
let timer: NodeJS.Timer | null = null
let isActive = ref(true)

const getSystemInfo = async () => {
  try {
    let res: SystemInfo = await reqSystemInfo()

    if (res.code === 0) {
      systemBaord.value = res
      console.log(systemBaord.value)
    } else {
      ElMessage({
        type: 'error',
        message: '请求失败',
      })
    }
  } catch (error) {
    // 在这里处理超时情况
    ElMessage({
      type: 'error',
      message: '服务器异常',
    })
  }
}
function loadStatus(val: number) {
  if (val < 30) {
    return '运行流畅'
  }
  if (val < 70) {
    return '运行正常'
  }
  if (val < 80) {
    return '运行缓慢'
  }
  return '运行阻塞'
}
const formatDiskCapacity = (bytes: number) => {
  const K = 1024
  if (bytes >= K * K * K) {
    return (bytes / (K * K * K)).toFixed(2) + ' GB'
  } else if (bytes >= K * K) {
    return (bytes / (K * K)).toFixed(2) + ' MB'
  } else {
    return bytes + ' bytes'
  }
}
const timeFormat = (timestamp: string) => {
  const date = new Date(timestamp)
  return date.toLocaleString()
}
const onFocus = () => {
  isActive.value = true
}
const onBlur = () => {
  isActive.value = false
}
onMounted(async () => {
  window.addEventListener('focus', onFocus)
  window.addEventListener('blur', onBlur)
  getSystemInfo()
  timer = setInterval(async () => {
    if (isActive.value) {
      await getSystemInfo()
    }
  }, 2000)
})
onBeforeUnmount(() => {
  window.removeEventListener('focus', onFocus)
  window.removeEventListener('blur', onBlur)
  clearInterval(Number(timer))
  timer = null
})
const formatSDKInfo = (info: string): string => {
  return info.replace(/\n/g, '<br>')
}
</script>

<style scoped>
.el-header {
  background-color: #b3c0d1;
  color: #333;
  padding: 30px;
}

.box-card {
  margin-bottom: 20px;
}

.card-header {
  background-color: #f0f0f0;
  padding: 10px;
  border-bottom: 1px solid #ccc;
}

.card-content {
  padding: 20px;
}

p {
  margin-bottom: 10px;
}
.sdk-info {
  /* display: block; */
  text-indent: 2em; /* 设置首行缩进为2em，可以根据需要调整 */
  line-height: 1.5; /* 设置行间距为1.5，可以根据需要调整 */
}

strong {
  font-weight: bold;
}
</style>
