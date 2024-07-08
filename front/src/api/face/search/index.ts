import request from '@/utils/request'
import type { ResAlarmList, ReqData, SearchReqData ,ResInfo, Item, ResponseData, TaskId, AlarmNum, SearchResReq } from './type'
enum API {
  ALARMLIST_URL = 'api/face/alarm/list',
  GETINFOS_URL = 'api/face/alarm/info',
  DELETEALARM_URL = 'api/face/alarm/delete',
  ADDTASK_URL = 'api/face/task/add',
  DELETE_URL = 'api/face/task/delete',

  ADDSEARCHTASK_URL = 'api/face/search/add',
  SEARCHRESULTLIST_URL = 'api/face/alarm/search/res/list'
}

export const getInfo = () => request.get<any, ResInfo>(API.GETINFOS_URL)

export const reqAlarmList = (data: ReqData) => request.post<any, ResAlarmList>(API.ALARMLIST_URL, data)

export const reqSearchResList = (data: SearchResReq) => request.post<any, ResAlarmList>(API.SEARCHRESULTLIST_URL, data)

export const reqAddSearchTask = (data: SearchReqData) => request.post<any, ResponseData>(API.ADDSEARCHTASK_URL, data)

export const addTask = (data: Item) => request.post<any, ResponseData>(API.ADDTASK_URL, data)

export const reqDeleteTask = (data: TaskId) => request.post<any, any>(API.DELETE_URL, data)

export const reqDeleteAlarm = (data: AlarmNum) => request.post<any, any>(API.DELETEALARM_URL, data)

export const reqSearchLimit = (data: ReqData) => request.post<any, ResAlarmList>(API.ALARMLIST_URL, data)
