import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      redirect: { name: 'login' },  // 默认重定向到登录页面
      // component: HomeView,  // 登录页面组件
    },
    {
      path: '/login',
      name: 'login',
      component: HomeView,  // 登录页面组件
    },
  ],
})

export default router
