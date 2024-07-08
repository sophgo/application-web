import { createRouter, createWebHashHistory } from 'vue-router'
import { constantRoute, anyRoute } from './routes'

const router = createRouter({
  history: createWebHashHistory(),
  routes: constantRoute, //[...constantRoute, anyRoute]
  // 滚动行为
  scrollBehavior() {
    return {
      left: 0,
      top: 0,
    }
  },
})

export default router
