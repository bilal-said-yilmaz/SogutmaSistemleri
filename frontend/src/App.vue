<template>
  <div id="app">
    <!-- Top Contact Bar -->
    <div class="top-contact-bar" v-if="!isAdminSite">
      <div class="container">
        <div class="contact-buttons">
          <a href="https://wa.me/+905385153191" target="_blank" class="contact-button whatsapp-button">
            <i class="fab fa-whatsapp"></i>
            WhatsApp'tan Ulaşın
          </a>
          <a href="tel:+905385153191" class="contact-button call-button">
            <i class="fas fa-phone-alt"></i>
            +90 (538) 515 3191
          </a>
        </div>
      </div>
    </div>

    <!-- Navigation -->
    <nav class="navbar navbar-expand-lg navbar-dark fixed-top" id="mainNav" :class="{ 'admin-nav': isAdminSite }">
      <div class="container">
        <a class="navbar-brand" :href="isAdminSite ? mainSiteUrl : '/'">
          <img :src="logo" alt="Klima Kozan Logo">
          <span class="brand-text">KOZAN İKLİMLENDİRME</span>
        </a>
        <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarResponsive">
          <span class="navbar-toggler-icon"></span>
        </button>
        <div class="collapse navbar-collapse" id="navbarResponsive">
          <ul class="navbar-nav text-uppercase ms-auto py-4 py-lg-0">
            <li class="nav-item">
              <a class="nav-link" :href="mainSiteUrl" target="_blank" rel="noopener noreferrer" v-if="isAdminSite">Ana Sayfa</a>
              <router-link class="nav-link" to="/" v-else>Ana Sayfa</router-link>
            </li>
            <template v-if="$route.name === 'home'">
              <li class="nav-item">
                <a class="nav-link" href="#services">Hizmetler</a>
              </li>
              <li class="nav-item">
                <a class="nav-link" href="#products">Ürünler</a>
              </li>
              <li class="nav-item">
                <a class="nav-link" href="#about">Hakkımızda</a>
              </li>
              <li class="nav-item">
                <a class="nav-link" href="#contact">İletişim</a>
              </li>
            </template>
            <!-- Admin sitesinde ve giriş yapmış kullanıcı için çıkış butonu -->
            <li class="nav-item" v-if="isAdminSite && isAuthenticated">
              <a class="nav-link" href="#" @click.prevent="handleLogout">
                <i class="fas fa-sign-out-alt me-1"></i>
                Çıkış Yap
              </a>
            </li>
          </ul>
        </div>
      </div>
    </nav>

    <!-- Main Content -->
    <router-view></router-view>
  </div>
</template>

<script setup>
import { computed, onMounted, onUnmounted } from 'vue'
import { useStore } from 'vuex'
import { useRouter } from 'vue-router'
import logo from './assets/logo.png'

const store = useStore()
const router = useRouter()

const isAuthenticated = computed(() => store.state.isAuthenticated)
const isAdminSite = computed(() => store.state.isAdminSite)
const mainSiteUrl = computed(() => {
  if (window.location.hostname === 'localhost' || window.location.hostname === 'admin.localhost') {
    return 'http://localhost:5173'
  }
  return `${window.location.protocol}//${window.location.hostname.replace('admin.', '')}`
})

const handleLogout = () => {
  store.dispatch('logout')
  // Login sayfasına yönlendir
  router.push('/login')
}

// Navbar scroll effect
const handleScroll = () => {
  const mainNav = document.getElementById('mainNav')
  const sections = ['services', 'products', 'about', 'contact']
  
  if (mainNav) {
    if (window.scrollY > 50) {
      mainNav.classList.add('navbar-shrink')
    } else {
      mainNav.classList.remove('navbar-shrink')
    }
  }

  // Aktif bölümü belirle
  sections.forEach(sectionId => {
    const section = document.getElementById(sectionId)
    const navLink = document.querySelector(`a[href="#${sectionId}"]`)
    
    if (section && navLink) {
      const sectionTop = section.offsetTop - 100
      const sectionBottom = sectionTop + section.offsetHeight
      
      if (window.scrollY >= sectionTop && window.scrollY < sectionBottom) {
        navLink.classList.add('active')
      } else {
        navLink.classList.remove('active')
      }
    }
  })
}

onMounted(() => {
  window.addEventListener('scroll', handleScroll)
  // Sayfa yüklendiğinde de kontrol et
  handleScroll()
})

onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll)
})
</script>

<style>
#app {
  font-family: 'Roboto', sans-serif;
}

/* Navigation */
#mainNav {
  padding: 0;
  background-color: rgba(33, 37, 41, 0.98);
  backdrop-filter: blur(10px);
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease-in-out;
  top: 0;
  height: 70px;
  display: flex;
  align-items: center;
  position: fixed;
  width: 100%;
}

#mainNav .container {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 100%;
  padding: 0 1.5rem;
}

#mainNav .navbar-brand {
  color: #fff;
  font-weight: 700;
  letter-spacing: 0.0625em;
  text-transform: uppercase;
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 0;
}

#mainNav .navbar-brand img {
  height: 50px;
  width: auto;
  transition: all 0.4s ease-in-out;
  object-fit: contain;
  filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.1));
}

#mainNav .navbar-brand .brand-text {
  font-size: 1.4rem;
  white-space: nowrap;
  color: #1e90ff;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  font-weight: 800;
}

#mainNav .navbar-nav .nav-item .nav-link {
  font-size: 1rem;
  color: #1e90ff;
  letter-spacing: 0.0625em;
  padding: 1.25rem 1rem;
  position: relative;
  transition: all 0.3s ease-in-out;
}

#mainNav .navbar-nav .nav-item .nav-link::after {
  content: '';
  position: absolute;
  bottom: 1rem;
  left: 1rem;
  right: 1rem;
  height: 2px;
  background: #1e90ff;
  transform: scaleX(0);
  transition: transform 0.3s ease-in-out;
  transform-origin: right;
}

#mainNav .navbar-nav .nav-item .nav-link:hover::after,
#mainNav .navbar-nav .nav-item .nav-link.active::after {
  transform: scaleX(1);
  transform-origin: left;
}

#mainNav .navbar-nav .nav-item .nav-link:hover,
#mainNav .navbar-nav .nav-item .nav-link.active {
  color: #0077e6;
}

#mainNav .navbar-nav .nav-item {
  position: relative;
  overflow: hidden;
}

#mainNav .navbar-nav .nav-item::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(30, 144, 255, 0.1);
  transform: translateY(100%);
  transition: transform 0.3s ease-in-out;
  z-index: -1;
}

#mainNav.navbar-shrink {
  padding: 0.8rem 0;
  background-color: rgba(33, 37, 41, 0.98);
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.1);
}

#mainNav.navbar-shrink .navbar-brand img {
  height: 55px;
}

/* Footer */
.footer {
  text-align: center;
  font-size: 0.9rem;
  background-color: #f8f9fa;
  padding: 4rem 0;
}

.btn-social {
  height: 2.5rem;
  width: 2.5rem;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 0;
  border-radius: 100%;
  transition: all 0.3s ease-in-out;
}

.btn-social:hover {
  transform: translateY(-3px);
  box-shadow: 0 0.5rem 1rem rgba(0, 0, 0, 0.15);
}

/* Global Styles */
section {
  padding: 6rem 0;
}

.section-heading {
  font-size: 2.5rem;
  margin-top: 0;
  margin-bottom: 1rem;
}

.section-subheading {
  font-size: 1rem;
  font-weight: 400;
  font-style: italic;
  margin-bottom: 4rem;
}

@media (max-width: 991.98px) {
  #mainNav {
    height: 70px;
  }

  #mainNav .container {
    padding: 0 1rem;
  }

  #mainNav .navbar-brand {
    padding: 0;
  }

  #mainNav .navbar-brand img {
    height: 45px;
  }
  
  #mainNav .navbar-brand .brand-text {
    font-size: 1.2rem;
  }

  #mainNav .navbar-collapse {
    position: absolute;
    top: 70px;
    left: 0;
    right: 0;
    background-color: rgba(33, 37, 41, 0.98);
    padding: 1rem;
    max-height: calc(100vh - 70px);
    overflow-y: auto;
  }

  #mainNav .navbar-nav .nav-item .nav-link {
    padding: 0.75rem 1rem;
    text-align: center;
  }
}

@media (max-width: 767.98px) {
  #mainNav {
    height: 70px;
  }

  #mainNav .container {
    padding: 0 0.75rem;
  }

  #mainNav .navbar-brand img {
    height: 40px;
  }
  
  #mainNav .navbar-brand .brand-text {
    font-size: 1.1rem;
  }

  #mainNav.navbar-shrink {
    height: 70px;
  }

  #mainNav.navbar-shrink .navbar-brand img {
    height: 40px;
  }

  .top-contact-bar {
    top: 70px;
  }
}

.brand-logo {
  display: none;
}

/* Top Contact Bar Styles */
.top-contact-bar {
  background: rgba(255, 255, 255, 0);

  padding: 0.5rem 0;
  position: fixed;
  top: 70px;
  left: 0;
  right: 0;
  z-index: 1030;
}

.contact-buttons {
  display: flex;
  justify-content: center;
  gap: 2rem;
  max-width: 600px;
  margin: 0 auto;
}

.contact-button {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.6rem 1.5rem;
  border-radius: 50px;
  color: #fff;
  text-decoration: none;
  font-weight: 600;
  transition: all 0.3s ease;
  font-size: 1rem;
  min-width: 180px;
  justify-content: center;
}

.whatsapp-button {
  background-color: rgba(37, 211, 102, 0.9);
  border: 1px solid rgba(37, 211, 102, 0.2);
}

.whatsapp-button:hover {
  background-color: rgba(34, 197, 94, 1);
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(37, 211, 102, 0.3);
}

.call-button {
  background-color: rgba(30, 144, 255, 0.9);
  border: 1px solid rgba(30, 144, 255, 0.2);
}

.call-button:hover {
  background-color: rgba(0, 119, 230, 1);
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(30, 144, 255, 0.3);
}

.contact-button i {
  font-size: 1.1rem;
}

@media (max-width: 768px) {
  .contact-buttons {
    gap: 1rem;
  }
  
  .contact-button {
    padding: 0.5rem 1rem;
    font-size: 0.9rem;
    min-width: 140px;
  }
  
  #mainNav {
    top: 0;
  }
}

/* Admin Panel Specific Styles */
#mainNav.admin-nav {
  top: 0;
  height: 70px;
}

#mainNav.admin-nav .container {
  padding: 0 1.5rem;
}
</style>