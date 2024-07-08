import type { RouteRecordRaw } from 'vue-router'

export const constantRoute : RouteRecordRaw[]= [
  {
    path: '/login',
    component: () => import('@/views/login/index.vue'),
    name: 'login',
    meta: {
      title: 'login',
      hidden: true,
    },
  },
  {
    path: '/',
    component: () => import('@/layout/index.vue'),
    name: 'layout',
    meta: {
      title: '',
      hidden: false,
      icon: '',
    },
    redirect: '/home',
    children: [
      {
        path: '/home',
        component: () => import('@/views/home/index.vue'),
        meta: {
          title: '首页',
          hidden: false,
          icon: 'HomeFilled',
        },
      },
    ],
  },
  {
    path: '/algo',
    component: () => import('@/layout/index.vue'),
    name: 'Algo',
    meta: {
      title: '算法业务',
      hidden: false,
      icon: 'SetUp',
    },
    redirect: '/algo/task',
    children: [
      {
        path: '/algo/task',
        component: () => import('@/views/algo/task/index.vue'),
        name: 'Task',
        meta: {
          title: '任务管理',
          hidden: false,
          icon: 'Files',
        },
      },
      {
        path: '/algo/alarm',
        component: () => import('@/views/algo/alarm/index.vue'),
        name: 'Alarm',
        meta: {
          title: '告警检索',
          hidden: false,
          icon: 'Picture',
        },
      },
      {
        path: '/algo/config',
        component: () => import('@/views/algo/config/index.vue'),
        name: 'Config',
        meta: {
          title: '算法参数配置',
          hidden: false,
          icon: 'List',
        },
      },
    ],
  },
  {
    path: '/face',
    component: () => import('@/layout/index.vue'),
    name: 'Face',
    meta: {
      title: '人脸业务',
      hidden: false,
      icon: 'Avatar',
    },
    children: [
      {
        path: '/face/task',
        component: () => import('@/views/face/task/index.vue'),
        name: 'FaceTask',
        meta: {
          title: '任务管理',
          hidden: false,
          icon: 'Files',
        },
      },
      {
        path: '/face/alarm',
        component: () => import('@/views/face/alarm/index.vue'),
        name: 'FaceAlarm',
        meta: {
          title: '人脸检索',
          hidden: false,
          icon: 'Picture',
        },
      },
      {
        path: '/face/search',
        component: () => import('@/views/face/search/index.vue'),
        name: 'FaceSearch',
        meta: {
          title: '以图搜图',
          hidden: false,
          icon: 'Aim',
        },
      },
      {
        path: '/face/compare',
        component: () => import('@/views/face/compare/index.vue'),
        name: 'FaceCompare',
        meta: {
          title: '人脸图片比对',
          hidden: false,
          icon: 'List',
        },
      },
      {
        path: '/face/library',
        component: () => import('@/views/face/library/index.vue'),
        name: 'FaceLibrary',
        meta: {
          title: '人脸库管理',
          hidden: false,
          icon: 'Avatar',
        },
      },
      // {
      //   path: '/face/gallery',
      //   component: () => import('@/views/face/gallery/index.vue'),
      //   name: 'FaceGallery',
      //   meta: {
      //     title: '人脸库管理',
      //     hidden: false,
      //     icon: 'List',
      //   },
      // },
      // {
      //   path: '/face/test',
      //   component: () => import('@/views/face/test/index.vue'),
      //   name: 'test',
      //   meta: {
      //     title: '测试',
      //     hidden: false,
      //     icon: 'List',
      //   },
      // },
    ],
  },
  {
    path: '/404',
    component: () => import('@/views/404/index.vue'),
    name: '404',
    meta: {
      title: '404',
      hidden: true,
    },
  },
  {
    path: '/:pathMatch(.*)*',
    redirect: '/404',
    name: 'Any',
    meta: {
      title: '任意路由',
      hidden: true,
    },
  },
]

export const asyncRoute = []

export const anyRoute = {
  path: '/:pathMatch(.*)*',
  redirect: '/404',
  name: 'Any',
  meta: {
    title: '任意路由',
    hidden: true,
  },
}
