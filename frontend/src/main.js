import { createApp } from 'vue'
import { createStore } from 'vuex'
import PrimeVue from 'primevue/config'
import axios from 'axios'
import router from './router'
import App from './App.vue'

// PrimeVue styles
import 'primevue/resources/themes/saga-blue/theme.css'
import 'primevue/resources/primevue.min.css'
import 'primeicons/primeicons.css'
// Bootstrap
import 'bootstrap/dist/css/bootstrap.min.css'
import 'bootstrap/dist/js/bootstrap.bundle.min.js'
// Font Awesome
import '@fortawesome/fontawesome-free/css/all.min.css'

// Axios configuration
const baseURL = window.location.hostname.startsWith('admin.') 
  ? import.meta.env.VITE_API_BASE_URL
  : window.location.hostname === 'localhost' 
    ? import.meta.env.VITE_API_BASE_URL
    : `${window.location.protocol}//${window.location.hostname.replace('admin.', '')}`

axios.defaults.baseURL = baseURL
axios.defaults.withCredentials = true

// Add token to Axios headers if exists
const token = localStorage.getItem('token')
if (token) {
  axios.defaults.headers.common['Authorization'] = `Bearer ${token}`
}

// Add Axios interceptor for handling 401 errors
axios.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response && error.response.status === 401) {
      // Token expired or invalid, logout user
      store.dispatch('logout')
      router.push('/login')
    }
    return Promise.reject(error)
  }
)

// Vuex store configuration
const store = createStore({
  state() {
    return {
      isAuthenticated: !!localStorage.getItem('token'),
      isAdminSite: window.location.hostname.startsWith('admin.'),
      user: null,
      products: [],
      services: [],
      about: {},
      contact: {},
      users: []
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
    },
    setUsers(state, users) {
      state.users = users
    }
  },
  actions: {
    async login({ commit }, credentials) {
      try {
        const response = await axios.post('/api/auth/login', credentials)
        const token = response.data.token
        localStorage.setItem('token', token)
        axios.defaults.headers.common['Authorization'] = `Bearer ${token}`
        commit('setAuth', true)
        return response
      } catch (error) {
        throw error
      }
    },
    logout({ commit }) {
      localStorage.removeItem('token')
      delete axios.defaults.headers.common['Authorization']
      commit('setAuth', false)
      commit('setUser', null)
      commit('setProducts', [])
      commit('setServices', [])
      commit('setAbout', {})
      commit('setContact', {})
    },
    async fetchUsers({ commit }) {
      try {
        const response = await axios.get('/api/admin/users')
        commit('setUsers', response.data)
      } catch (error) {
        console.error('Error fetching users:', error)
        throw error
      }
    },
    async createUser({ dispatch }, userData) {
      try {
        await axios.post('/api/admin/users', userData)
        await dispatch('fetchUsers')
      } catch (error) {
        console.error('Error creating user:', error)
        throw error
      }
    },
    async updateUser({ dispatch }, { id, data }) {
      try {
        await axios.put(`/api/admin/users/${id}`, data)
        await dispatch('fetchUsers')
      } catch (error) {
        console.error('Error updating user:', error)
        throw error
      }
    },
    async deleteUser({ dispatch }, id) {
      try {
        await axios.delete(`/api/admin/users/${id}`)
        await dispatch('fetchUsers')
      } catch (error) {
        console.error('Error deleting user:', error)
        throw error
      }
    }
  }
})

const app = createApp(App)

app.use(router)
app.use(store)
app.use(PrimeVue)

app.mount('#app') 