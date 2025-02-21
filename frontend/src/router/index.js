import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import Admin from '../views/Admin.vue'
import Login from '../views/Login.vue'

const isAdminSubdomain = window.location.hostname.startsWith('admin.')

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
      beforeEnter: (to, from, next) => {
        if (isAdminSubdomain) {
          next('/admin')
        } else {
          next()
        }
      }
    },
    {
      path: '/admin',
      name: 'admin',
      component: Admin,
      meta: { requiresAuth: true },
      beforeEnter: (to, from, next) => {
        if (!isAdminSubdomain) {
          window.location.href = `https://admin.${window.location.hostname.replace('www.', '')}${to.fullPath}`
        } else {
          next()
        }
      }
    },
    {
      path: '/login',
      name: 'login',
      component: Login,
      meta: { guest: true },
      beforeEnter: (to, from, next) => {
        if (!isAdminSubdomain) {
          window.location.href = `https://admin.${window.location.hostname.replace('www.', '')}${to.fullPath}`
        } else {
          next()
        }
      }
    }
  ],
  scrollBehavior(to, from, savedPosition) {
    if (to.hash) {
      return {
        el: to.hash,
        behavior: 'smooth'
      }
    }
    return { top: 0 }
  }
})

// Navigation guard
router.beforeEach((to, from, next) => {
  const isAuthenticated = localStorage.getItem('token')
  
  if (to.matched.some(record => record.meta.requiresAuth)) {
    // Auth gerektiren route için
    if (!isAuthenticated) {
      next({
        path: '/login',
        query: { redirect: to.fullPath }
      })
    } else {
      next()
    }
  } else if (to.matched.some(record => record.meta.guest)) {
    // Misafir route'u için (login gibi)
    if (isAuthenticated) {
      next('/admin')
    } else {
      next()
    }
  } else {
    next()
  }
})

export default router 