import request from '@/utils/request'
import type { ResAlarmList, ReqData, ResInfo, Item, ResponseData, TaskId, AlarmNum } from './type'
enum API {
  ALARMLIST_URL = 'api/alarm/list',
  GETINFOS_URL = 'api/alarm/info',
  DELETEALARM_URL = 'api/alarm/delete',
  ADDTASK_URL = 'api/task/add',
  DELETE_URL = 'api/task/delete',
}

export const getInfo = () => request.get<any, ResInfo>(API.GETINFOS_URL)

export const reqAlarmList = (data: ReqData) => request.post<any, ResAlarmList>(API.ALARMLIST_URL, data)

export const addTask = (data: Item) => request.post<any, ResponseData>(API.ADDTASK_URL, data)

export const reqDeleteTask = (data: TaskId) => request.post<any, any>(API.DELETE_URL, data)

export const reqDeleteAlarm = (data: AlarmNum) => request.post<any, any>(API.DELETEALARM_URL, data)
