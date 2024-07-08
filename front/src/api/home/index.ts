// 统一管理用户相关接口
import request from '@/utils/request'
import type { SystemInfo } from './type'

enum API {
  SYSTEMINFO_URL = '/api/dashboard/info',
}
export const reqSystemInfo = () => request.get<any, SystemInfo>(API.SYSTEMINFO_URL)
