import { createApp } from 'vue'
import { createRouter, createWebHashHistory } from 'vue-router'
import { createStore } from 'vuex'
import PrimeVue from 'primevue/config'
import axios from 'axios'
import 'primevue/resources/themes/saga-blue/theme.css'
import 'primevue/resources/primevue.min.css'
import 'primeicons/primeicons.css'


// Axios yapılandırması
axios.defaults.baseURL = 'http://localhost:3000'

// Token varsa Axios headers'ına ekle
const token = localStorage.getItem('token')
if (token) {
  axios.defaults.headers.common['Authorization'] = `Bearer ${token}`
}

import App from './App.vue'
import Home from './views/Home.vue'
import Admin from './views/Admin.vue'
import Login from './views/Login.vue'

// Router yapılandırması
const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    {
      path: '/',
      name: 'Home',
      component: Home
    },
    {
      path: '/admin',
      name: 'Admin',
      component: Admin,
      meta: { requiresAuth: true }
    },
    {
      path: '/admin/login',
      name: 'Login',
      component: Login
    }
  ]
})

// Vuex store yapılandırması
const store = createStore({
  state() {
    return {
      isAuthenticated: !!localStorage.getItem('token'),
      user: null,
      products: [],
      services: [],
      about: {},
      contact: {}
    }
  },
  mutations: {
    setAuth(state, status) {
      state.isAuthenticated = status
    },
    setUser(state, user) {
      state.user = user
    },
    setProducts(state, products) {
      state.products = products
    },
    setServices(state, services) {
      state.services = services
    },
    setAbout(state, about) {
      state.about = about
    },
    setContact(state, contact) {
      state.contact = contact
    }
  }
})

// Auth kontrolü
router.beforeEach((to, from, next) => {
  if (to.matched.some(record => record.meta.requiresAuth)) {
    if (!store.state.isAuthenticated) {
      next({
        path: '/admin/login',
        query: { redirect: to.fullPath }
      })
    } else {
      next()
    }
  } else {
    next()
  }
})

const app = createApp(App)

app.use(router)
app.use(store)
app.use(PrimeVue)

app.mount('#app') 