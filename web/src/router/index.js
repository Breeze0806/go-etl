import Vue from 'vue'
import VueRouter from 'vue-router'
import Home from '../views/Home.vue'
import Login from '../views/Login.vue'
import Register from '../views/Register.vue'

Vue.use(VueRouter)

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: Login,
    meta: { requiresAuth: false }
  },
  {
    path: '/register',
    name: 'Register',
    component: Register,
    meta: { requiresAuth: false }
  },
  {
    path: '/',
    name: 'Dashboard',
    component: Home,
    meta: { requiresAuth: true }
  },
  {
    path: '/datasources',
    name: 'DataSources',
    component: () => import('../views/DataSources.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/synctasks',
    name: 'SyncTasks',
    component: () => import('../views/SyncTasks.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/synctasks/new',
    name: 'NewSyncTask',
    component: () => import('../views/SyncTaskWizard.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '*',
    redirect: '/'
  }
]

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes
})

// Auth guard
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('jwt_token')
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth)

  if (requiresAuth && !token) {
    next({
      path: '/login',
      query: { redirect: to.fullPath }
    })
  } else if (to.path === '/login' && token) {
    next('/')
  } else {
    next()
  }
})

export default router
