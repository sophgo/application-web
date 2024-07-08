export interface ResponseData {
  code: number
  msg: string
}

export interface ResConfig extends ResponseData {
  result: ConfigResult
}

export interface ConfigResult {
  device: Device
  algorithms: Algorithm[]
}

export interface Device {
  codeName: string
  name: string
  resolution: string
  url: string
  width: number
  height: number
}

export interface Algorithm {
  Type: number
  TrackInterval: number
  DetectInterval: number
  AlarmInterval: number
  threshold: number
  TargetSize: TargetSize
  DetectInfos: DetectInfo[]
}

export interface ModAlgorithm extends Algorithm {
  taskId: string
  Algorithm: Algorithm
}

export interface TargetSize {
  MinDetect: number
  MaxDetect: number
}

export interface Point {
  X: number
  Y: number
}

export interface DetectInfo {
  TripWire: TripWire
  HotArea: Point[]
}
export interface TripWire {
  LineStart: Point
  LineEnd: Point
  DirectStart: Point
  DirectEnd: Point
}

export interface TaskId {
  taskId: string
}
