export interface ResponseData {
  code: number
  msg: string
}

export interface ReqData {
  pageNo: number
  pageSize: number
  beginTime?: number
  endTime?: number
  srcId?: string
  taskId?: string
  trackId?: number
  types?: number[]
}

export interface AlarmRes {
  pageNo: number
  pageSize: number
  total: number
  pageCount: number
  items: AlarmItem[]
}

export interface AlarmItem {
  id: number
  taskId: string
  srcId: string
  status?: number
  frameIndex?: number
  abilities?: string[]
  type: number
  smallImage: string
  bigImage: string
  time: number
  Extend?: string
  box: Box
}
export interface Box {
  LeftTopY: number
  RightBtmY: number
  LeftTopX: number
  RightBtmX: number
}
export interface Info {
  types: number[]
  taskIds: string[]
  srcIds: string[]
}
export interface ResAlarmList extends ResponseData {
  result: AlarmRes
}

export interface ResInfo extends ResponseData {
  result: Info
}

export interface TaskId {
  taskId: string
}

export interface AlarmNum {
  number: number
}

export interface Item {
  taskId: string
  deviceName: string
  url: string
  status?: number
  errorReason?: string
  abilities?: string[]
  types: number[]
}
