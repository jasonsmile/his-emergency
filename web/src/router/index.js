import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'Login',
      component: () => import('@/views/Login.vue'),
      meta: { title: '登录' },
    },
    {
      path: '/',
      component: () => import('@/layout/MainLayout.vue'),
      redirect: '/registration',
      meta: { requiresAuth: true },
      children: [
        {
          path: 'registration',
          name: 'RegistrationWorkbench',
          component: () => import('@/views/registration/Workbench.vue'),
          meta: { title: '挂号工作台' },
        },
        {
          path: 'registration/list',
          name: 'EncounterList',
          component: () => import('@/views/registration/EncounterList.vue'),
          meta: { title: '就诊列表' },
        },
        {
          path: 'charge',
          name: 'ChargeWorkbench',
          component: () => import('@/views/charge/Workbench.vue'),
          meta: { title: '收费工作台' },
        },
        {
          path: 'charge/report',
          name: 'DailyReport',
          component: () => import('@/views/charge/DailyReport.vue'),
          meta: { title: '日结报表' },
        },
      ],
    },
    { path: '/:pathMatch(.*)*', redirect: '/' },
  ],
})

router.beforeEach(to => {
  if (to.meta.requiresAuth && !localStorage.getItem('token')) {
    return { name: 'Login', query: { redirect: to.fullPath }, replace: true }
  }
})

export default router
