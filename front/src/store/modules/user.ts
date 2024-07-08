import { defineStore } from 'pinia'
import router from '@/router'
import { reqLogin, reqUserInfo, reqLogOut } from '@/api/user'
import type { LoginFormData, LoginResponseData, userInfoResponseData } from '@/api/user/type'
import type { UserState } from './types/types'
import { SET_TOKEN, GET_TOKEN, REMOVE_TOKEN } from '@/utils/token'
import { constantRoute, asyncRoute, anyRoute } from '@/router/routes'

// @ts-ignore
import cloneDeep from 'lodash/cloneDeep'
import setting from '@/setting'

let dynamicRoutes: any = []

function filterAsyncRoute(asyncRoute: any, routes: any) {
  return asyncRoute.filter((item: any) => {
    if (routes.includes(item.name)) {
      if (item.children && item.children.length > 0) {
        item.children = filterAsyncRoute(item.children, routes)
      }
      return true
    }
  })
}

let useUserStore = defineStore('User', {
  // 小仓库存储数据的地方
  state: (): UserState => {
    return {
      token: GET_TOKEN()!,
      menuRoutes: constantRoute,
      username: '',
      avatar: setting.admin,
      buttons: [],
    }
  },
  // 异步|逻辑的地方
  actions: {
    //用户登录方法
    async userLogin(data: LoginFormData) {
      let res: LoginResponseData = await reqLogin(data)
      // success=>token
      // error=>error.message
      if (res.code === 0) {
        this.token = res.result?.token as string
        // 持久化
        SET_TOKEN(this.token)
        // console.log(this.token)
        return 'ok'
      } else {
        return Promise.reject(new Error(res.msg))
      }
    },

    async userLogout() {
      let logoutForm = {
        token: this.token,
      }
      REMOVE_TOKEN()
      let res = await reqLogOut(logoutForm)
      if (res.code === 0) {
        this.token = ''
        this.username = ''
        this.avatar = ''
        REMOVE_TOKEN()
        dynamicRoutes.forEach((route: any) => {
          router.removeRoute(route.name)
        })
      } else {
        return Promise.reject(new Error(res.message))
      }
    },
  },
  getters: {},
})

export default useUserStore
