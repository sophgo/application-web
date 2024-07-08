// 登录接口需要携带参数ts类型
export interface LoginFormData {
  username?: string
  password?: string
}

export interface LogoutFormData {
  token: string | null
}

export interface ResponseData {
  code?: number
  msg?: string
}

export interface ResultData {
  token: string
}

export interface LoginResponseData extends ResponseData {
  result?: ResultData
}

export interface userInfoResponseData extends ResponseData {
  data: {
    routes: string[]
    buttons: string[]
    roles: string[]
    name: string
    avatar: string
  }
}
