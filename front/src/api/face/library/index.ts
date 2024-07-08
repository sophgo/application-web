import request from '@/utils/request'
import type { ResAlarmList, ReqData, SearchReqData ,ResInfo, reqPersonInfo, Item, PersonInfo, ResponseData, TaskId, AlarmNum, SearchResReq, ResPersonList } from './type'
enum API {
  ALARMLIST_URL = 'api/face/alarm/list',
  GETINFOS_URL = 'api/face/alarm/info',
  DELETEALARM_URL = 'api/face/alarm/delete',
  ADDPERSON_URL = 'api/face/library/add',
  LISTPERSON_URL = 'api/face/library/list',
  DELETE_URL = 'api/face/library/delete',

  ADDSEARCHTASK_URL = 'api/face/search/add',
  SEARCHRESULTLIST_URL = 'api/face/alarm/search/res/list'
}

export const getInfo = () => request.get<any, ResInfo>(API.GETINFOS_URL)

export const reqAlarmList = (data: ReqData) => request.post<any, ResAlarmList>(API.ALARMLIST_URL, data)

export const reqSearchResList = (data: SearchResReq) => request.post<any, ResAlarmList>(API.SEARCHRESULTLIST_URL, data)

export const reqAddSearchTask = (data: SearchReqData) => request.post<any, ResponseData>(API.ADDSEARCHTASK_URL, data)

export const reqAddPerson = (data: reqPersonInfo) => request.post<any, ResponseData>(API.ADDPERSON_URL, data)

export const reqListPerson = () => request.get<any, ResPersonList>(API.LISTPERSON_URL)

export const reqDeletePerson = (data: PersonInfo) => request.post<any, ResponseData>(API.DELETE_URL, data)

export const reqDeleteAlarm = (data: AlarmNum) => request.post<any, any>(API.DELETEALARM_URL, data)

export const reqSearchLimit = (data: ReqData) => request.post<any, ResAlarmList>(API.ALARMLIST_URL, data)
