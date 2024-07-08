import type { UploadFile, UploadUserFile } from 'element-plus'

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

export interface Features extends ResponseData {
  result: string[]
}

export interface CompareResult extends ResponseData {
  result: {
    Similarity: number
  }
}

export interface FileID {
  FileID: string
}

export interface FileIDs {
  FileID: string[]
}

export interface UploadConfig {
  fileList: UploadUserFile[]
  fileName: string
  action: string
  listType: string
  limit: number
  disabled: boolean
}

export interface localFile extends UploadFile {
  id: number
}
