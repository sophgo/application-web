// 统一管理用户相关接口
import request from '@/utils/request'
import type { LoginFormData, LoginResponseData, userInfoResponseData, LogoutFormData } from './type'

enum API {
  LOGIN_URL = '/api/login',
  USERINFO_URL = '/admin/acl/index/info',
  LOGOUT_URL = '/api/logout',
}

export const reqLogin = (data: LoginFormData) => request.post<any, LoginResponseData>(API.LOGIN_URL, data)

export const reqUserInfo = () => request.get<any, userInfoResponseData>(API.USERINFO_URL)

export const reqLogOut = (data: LogoutFormData) => request.post<any, any>(API.LOGOUT_URL, data)
