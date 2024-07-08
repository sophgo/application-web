export interface ResponseData {
  code: number
  msg: string
}

export interface SystemInfo extends ResponseData {
  result: Result
}

interface Result {
  hostname: string
  os: string
  platform: string
  platformFamily: string
  platformVersion: string
  kernelArch: string
  kernelVersion: string
  virtualizationSystem: string
  cpuCores: number
  cpuLogicalCores: number
  cpuModelName: string
  currentInfo: CurrentInfo
  SDKVersion: string
}

interface CurrentInfo {
  uptime: number
  timeSinceUptime: string
  procs: number
  load1: number
  load5: number
  load15: number
  loadUsagePercent: number
  cpuPercent: number[]
  cpuUsedPercent: number
  cpuUsed: number
  cpuTotal: number
  memoryTotal: number
  memoryAvailable: number
  memoryUsed: number
  memoryUsedPercent: number
  swapMemoryTotal: number
  swapMemoryAvailable: number
  swapMemoryUsed: number
  swapMemoryUsedPercent: number
  npuMemoryTotal: number
  npuMemoryUsed: number
  vppemoryTotal: number
  vppMemoryUsed: number
  vpuMemoryTotal: number
  vpuMemoryUsed: number
  tpuUsed: number
  ioReadBytes: number
  ioWriteBytes: number
  ioCount: number
  ioReadTime: number
  ioWriteTime: number
  diskData: DiskData[]
  netBytesSent: number
  netBytesRecv: number
  shotTime: string
}

interface DiskData {
  path: string
  type: string
  device: string
  total: number
  free: number
  used: number
  usedPercent: number
  inodesTotal: number
  inodesUsed: number
  inodesFree: number
  inodesUsedPercent: number
}
