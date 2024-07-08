export interface ResponseData {
  code: number
  msg: string
}

export interface ReqData {
  pageNo: number
  pageSize: number
}

export interface AbilityList {
  type: number
  name: string
}

export interface ResAbilityList extends ResponseData {
  result: AbilityList[]
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

export interface TaskId {
  taskId: string
}

export interface UpdateImg extends TaskId {
  url: string
}

export interface TaskListResponseData extends ResponseData {
  result: {
    items: Item[]
    total: number
    pageSize: number
    pageCount: number
    pageNo: number
  }
}
