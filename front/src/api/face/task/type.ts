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

export interface WebStream {
  socket: WebSocket
  taskid: string
  domain: string
}

export class WebStreamImpl implements WebStream {
  socket: WebSocket;
  taskid: string;
  domain: string;

  init(domain: string, taskId: string = '') {
    this.taskid = taskId;
    this.domain = domain;
  }

  connect() {
    this.socket = new WebSocket(this.domain);
  }

  disconnect() {
    this.socket.close()
  }
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
  TaskID: string
}

export interface UpdateImg extends TaskId {
  url: string
}

export interface TaskListResponseData extends ResponseData {
  result: {
    Items: RequestBody[]
    total: number
    pageSize: number
    pageCount: number
    pageNo: number
  }
}
export interface RequestBody {
  TaskID: string
  InputSrc: {
    SrcID: string
    StreamSrc: {
      Address: string
    }
  }
  Algorithm: {
    TrackInterval: number
    DetectInterval: number
    AttributeInclude: boolean
    FeatureInclude: boolean
    TargetSize: {
      MinDetect: number
      MaxDetect: number
    }
  }
  Status: number
}
