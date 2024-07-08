import { defineStore } from 'pinia'
import { CategoryState } from './types/types'

let useCategoryStore = defineStore('Category', {
  state: (): CategoryState => {
    return {
      c1Id: '',
      c2Id: '',
      c3Id: '',
    }
  },
  actions: {},
  getters: {},
})

export default useCategoryStore
