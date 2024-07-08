import request from '@/utils/request'
import type { ResConfig, TaskId, ModAlgorithm, ResponseData } from './type'
enum API {
  MODCONFIG_URL = 'api/config/mod',
  GETCONFIG_URL = 'api/config/get',
}

export const reqTaskConfig = (data: TaskId) => request.post<any, ResConfig>(API.GETCONFIG_URL, data)

export const ModTaskConfig = (data: ModAlgorithm) => request.post<any, ResponseData>(API.MODCONFIG_URL, data)
