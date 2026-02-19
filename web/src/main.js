import Vue from 'vue'
import App from './App.vue'
import router from './router'
import Element from 'element-ui'
import 'element-ui/lib/theme-chalk/index.css'
import axios from 'axios'

Vue.config.productionTip = false

axios.defaults.baseURL = process.env.VUE_APP_API_BASE_URL || 'http://localhost:8080'
axios.interceptors.request.use(
  config => {
    const token = localStorage.getItem('jwt_token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

axios.interceptors.response.use(
  response => response,
  error => {
    if (error.response && error.response.status === 401) {
      localStorage.removeItem('jwt_token')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

Vue.prototype.$http = axios

Vue.use(Element, {
  size: 'medium',
  zIndex: 3000
})

new Vue({
  router,
  render: h => h(App)
}).$mount('#app')
