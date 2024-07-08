import type { RouteRecordRaw } from 'vue-router'
export interface UserState {
  token: string | null
  menuRoutes: RouteRecordRaw[]
  username: string
  avatar: string
  buttons: string[]
}

export interface CategoryState {
  c1Id: string | number
  c2Id: string | number
  c3Id: string | number
}
