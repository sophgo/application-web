import request from '@/utils/request'
import type { ReqData, TaskListResponseData, ResAbilityList, Item, ResponseData, TaskId, UpdateImg } from './type'
enum API {
  TASKLIST_URL = 'api/task/list',
  GETABILITIES_URL = 'api/task/abilities',
  ADDTASK_URL = 'api/task/add',
  MODTASK_URL = 'api/task/modify',
  DELETE_URL = 'api/task/delete',
  STARTTASK_URL = 'api/task/start',
  STOPTASK_URL = 'api/task/stop',
  GETPIC_URL = 'api/task/image',
}

export const getAbilites = () => request.get<any, ResAbilityList>(API.GETABILITIES_URL)

export const reqTaskList = (data: ReqData) => request.post<any, TaskListResponseData>(API.TASKLIST_URL, data)

export const addTask = (data: Item) => request.post<any, ResponseData>(API.ADDTASK_URL, data)

export const modTask = (data: Item) => request.post<any, ResponseData>(API.MODTASK_URL, data)

export const deleteTask = (data: TaskId) => request.post<any, ResponseData>(API.DELETE_URL, data)

export const startTask = (data: TaskId) => request.post<any, ResponseData>(API.STARTTASK_URL, data)

export const stopTask = (data: TaskId) => request.post<any, ResponseData>(API.STOPTASK_URL, data)

export const getPic = (data: UpdateImg) =>
  request.post<any, ResponseData>(API.GETPIC_URL, data, {
    timeout: 15000, // 设置超时时间为15秒
  })
