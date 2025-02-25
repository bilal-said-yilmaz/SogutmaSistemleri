<template>
  <div id="app">
    <link href="https://fonts.googleapis.com/css2?family=Montserrat:wght@600;700;800&display=swap" rel="stylesheet">
    <link href="https://fonts.googleapis.com/css2?family=Poppins:wght@400;500;600;700&display=swap" rel="stylesheet">
    <!-- Top Contact Bar -->
    <div class="top-contact-bar" v-if="!isAdminSite">
      <div class="container">
        <div class="contact-buttons">
          <a :href="'https://wa.me/' + (contact?.phone?.replace(/\D/g, '') || '+905385153191')" target="_blank" class="contact-button whatsapp-button">
            <i class="fab fa-whatsapp"></i>
            WhatsApp'tan Ulaşın
          </a>
          <a :href="'tel:' + (contact?.phone || '+905385153191')" class="contact-button call-button">
            <i class="fas fa-phone-alt"></i>
            {{ contact?.phone || '+90 (538) 515 3191' }}
          </a>
        </div>
      </div>
    </div>

    <!-- Navigation -->
    <nav class="navbar navbar-expand-lg navbar-dark fixed-top" id="mainNav" :class="{ 'admin-nav': isAdminSite }">
      <div class="container">
        <a class="navbar-brand" :href="isAdminSite ? mainSiteUrl : '/'" :target="isAdminSite ? '_blank' : '_self'" :rel="isAdminSite ? 'noopener noreferrer' : ''">
          <img :src="logo" alt="Klima Kozan Logo">
          <span class="brand-text">KOZAN İKLİMLENDİRME</span>
        </a>
        <button 
          class="navbar-toggler" 
          type="button" 
          @click="toggleMenu"
          :class="{ 'collapsed': !isMenuOpen }"
        >
          <i class="fas" :class="isMenuOpen ? 'fa-xmark' : 'fa-bars-staggered'"></i>
        </button>
        <div class="collapse navbar-collapse" id="navbarResponsive">
          <ul class="navbar-nav text-uppercase ms-auto py-4 py-lg-0">
            <li class="nav-item">
              <a class="nav-link" :href="mainSiteUrl" target="_blank" rel="noopener noreferrer" v-if="isAdminSite" @click="closeMenu">Ana Sayfa</a>
              <router-link class="nav-link" to="/" v-else @click="closeMenu">Ana Sayfa</router-link>
            </li>
            <template v-if="$route.name === 'home'">
              <li class="nav-item">
                <a class="nav-link" href="#services" @click="closeMenu">Hizmetler</a>
              </li>
              <li class="nav-item">
                <a class="nav-link" href="#products" @click="closeMenu">Ürünler</a>
              </li>
              <li class="nav-item">
                <a class="nav-link" href="#about" @click="closeMenu">Hakkımızda</a>
              </li>
              <li class="nav-item">
                <a class="nav-link" href="#contact" @click="closeMenu">İletişim</a>
              </li>
            </template>
            <li class="nav-item" v-if="isAdminSite && isAuthenticated">
              <a class="nav-link" href="#" @click.prevent="handleLogout(); closeMenu()">
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

    <!-- Footer -->
    <footer class="footer">
      <div class="container">
        <div class="row">
          <div class="col-md-4 footer-section">
            <h3>Kozan İklimlendirme</h3>
            <p>15 yılı aşkın tecrübemizle Kozan'da profesyonel klima ve beyaz eşya servisi hizmeti veriyoruz.</p>
          </div>
          <div class="col-md-4 footer-section">
            <h3>Hızlı Linkler</h3>
            <ul class="list-unstyled">
              <li><a href="#services"><i class="fas fa-chevron-right"></i>Hizmetlerimiz</a></li>
              <li><a href="#products"><i class="fas fa-chevron-right"></i>Ürünlerimiz</a></li>
              <li><a href="#about"><i class="fas fa-chevron-right"></i>Hakkımızda</a></li>
              <li><a href="#contact"><i class="fas fa-chevron-right"></i>İletişim</a></li>
            </ul>
          </div>
          <div class="col-md-4 footer-section">
            <h3>İletişim</h3>
            <ul class="list-unstyled">
              <li><i class="fas fa-phone"></i>{{ contact?.phone || '+90 (538) 515 3191' }}</li>
              <li><i class="fas fa-envelope"></i>{{ contact?.email || 'kozaniklimlendirme@gmail.com' }}</li>
              <li><i class="fas fa-map-marker-alt"></i>{{ contact?.address || 'Kozan, Adana' }}</li>
            </ul>
          </div>
        </div>
        <div class="footer-bottom text-center mt-4">
          <p>&copy; {{ new Date().getFullYear() }} Kozan İklimlendirme. Tüm hakları saklıdır.</p>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup>
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { useStore } from 'vuex'
import { useRouter } from 'vue-router'
import { Collapse } from 'bootstrap'
import logo from './assets/logo.png'
import axios from 'axios'

const store = useStore()
const router = useRouter()
const isMenuOpen = ref(false)
const contact = ref(null)

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

const toggleMenu = () => {
  const navbarCollapse = document.getElementById('navbarResponsive')
  if (navbarCollapse) {
    if (isMenuOpen.value) {
      // Menü açıksa kapat
      navbarCollapse.classList.remove('show')
      isMenuOpen.value = false
      document.body.style.overflow = 'auto'
    } else {
      // Menü kapalıysa aç
      navbarCollapse.classList.add('show')
      isMenuOpen.value = true
      if (window.innerWidth < 992) {
        document.body.style.overflow = 'hidden'
      }
    }
  }
}

const closeMenu = () => {
  const navbarCollapse = document.getElementById('navbarResponsive')
  if (navbarCollapse) {
    navbarCollapse.classList.remove('show')
    isMenuOpen.value = false
    document.body.style.overflow = 'auto'
  }
}

// Sayfa değiştiğinde menüyü kapat
const handleNavigation = () => {
  closeMenu()
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

// Fetch contact data
const fetchContact = async () => {
  try {
    const response = await axios.get('/api/contact')
    contact.value = response.data
  } catch (error) {
    console.error('İletişim bilgileri yüklenirken hata oluştu:', error)
  }
}

onMounted(() => {
  window.addEventListener('scroll', handleScroll)
  handleScroll()
  fetchContact()
  
  // Collapse olaylarını dinle
  const navbarCollapse = document.getElementById('navbarResponsive')
  if (navbarCollapse) {
    navbarCollapse.addEventListener('show.bs.collapse', () => {
      isMenuOpen.value = true
      if (window.innerWidth < 992) {
        document.body.style.overflow = 'hidden'
      }
    })
    
    navbarCollapse.addEventListener('hide.bs.collapse', () => {
      isMenuOpen.value = false
      document.body.style.overflow = 'auto'
    })
    
    // Link tıklamalarını dinle
    const navLinks = navbarCollapse.querySelectorAll('.nav-link')
    navLinks.forEach(link => {
      link.addEventListener('click', () => {
        if (window.innerWidth < 992) {
          navbarCollapse.classList.remove('show')
          isMenuOpen.value = false
          document.body.style.overflow = 'auto'
        }
      })
    })
  }

  // Sayfa değişikliklerini dinle
  window.addEventListener('hashchange', handleNavigation)
  
  // Ekran boyutu değişikliklerini dinle
  window.addEventListener('resize', () => {
    if (window.innerWidth >= 992 && isMenuOpen.value) {
      closeMenu()
    }
  })
})

onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll)
  window.removeEventListener('hashchange', handleNavigation)
  window.removeEventListener('resize', () => {})
  
  // Event listener'ları temizle
  const navbarCollapse = document.getElementById('navbarResponsive')
  if (navbarCollapse) {
    navbarCollapse.removeEventListener('show.bs.collapse', () => {})
    navbarCollapse.removeEventListener('hide.bs.collapse', () => {})
    
    // Link event listener'larını temizle
    const navLinks = navbarCollapse.querySelectorAll('.nav-link')
    navLinks.forEach(link => {
      link.removeEventListener('click', () => {})
    })
  }
})
</script>

<style>
#app {
  font-family: 'Roboto', sans-serif;
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

/* Navigation */
#mainNav {
  padding: 0;
  background-color: #000000;
  backdrop-filter: blur(10px);
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.2);
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
  font-family: 'Montserrat', sans-serif;
  font-size: 1.4rem;
  white-space: nowrap;
  text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.15);
  font-weight: 800;
  letter-spacing: 0.02em;
  text-transform: uppercase;
  background: linear-gradient(135deg, rgba(255, 255, 255, 1) 0%, rgba(240, 240, 240, 0.9) 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  transition: all 0.3s ease;
  position: relative;
}

#mainNav .navbar-brand .brand-text::after {
  content: '';
  position: absolute;
  bottom: -4px;
  left: 0;
  width: 100%;
  height: 2px;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.8), transparent);
  transform: scaleX(0);
  transition: transform 0.3s ease;
}

#mainNav .navbar-brand:hover .brand-text {
  background: linear-gradient(135deg, rgba(255, 255, 255, 1) 0%, rgba(230, 230, 230, 0.95) 100%);
  -webkit-background-clip: text;
  background-clip: text;
  transform: translateY(-1px);
  text-shadow: 3px 3px 6px rgba(0, 0, 0, 0.2);
}

#mainNav .navbar-nav .nav-item .nav-link {
  font-size: 1rem;
  color: #fff;
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
  background: #fff;
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
  color: rgba(255, 255, 255, 0.9);
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
  background-color: #000000;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.2);
}

#mainNav.navbar-shrink .navbar-brand img {
  height: 55px;
}

/* Footer */
.footer {
  font-size: 0.9rem;
  background-color: #000000;
  color: #fff;
  padding: 4rem 0 2rem;
  margin-top: auto;
  width: 100%;
  box-shadow: 0 -2px 10px rgba(0, 0, 0, 0.2);
}

.footer-section {
  text-align: left;
}

.footer-section h3 {
  font-size: 1.5rem;
  margin-bottom: 1.5rem;
  color: #fff;
  font-weight: 600;
  text-align: left;
}

.footer-section ul {
  padding-left: 0;
}

.footer-section ul li {
  margin-bottom: 0.5rem;
}

.footer-section ul li a {
  color: rgba(255, 255, 255, 0.8);
  text-decoration: none;
  transition: color 0.3s ease;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.footer-section ul li a i {
  width: 20px;
}

.footer-section ul li a:hover {
  color: #fff;
}

.footer-section ul li i {
  margin-right: 0.5rem;
  color: #fff;
}

.footer-section p {
  color: rgba(255, 255, 255, 0.8);
  line-height: 1.6;
  text-align: left;
}

.footer-bottom {
  color: rgba(255, 255, 255, 0.6);
  font-size: 0.9rem;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
  padding-top: 2rem;
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
    position: fixed;
    top: 70px;
    left: 0;
    right: 0;
    background: #000000;
    backdrop-filter: blur(10px);
    padding: 1rem;
    border-radius: 0 0 1rem 1rem;
    margin-top: 0;
    max-height: calc(100vh - 70px);
    overflow-y: auto;
    transition: all 0.3s ease-in-out;
    opacity: 0;
    transform: translateY(-10px);
    pointer-events: none;
  }

  #mainNav .navbar-collapse.show {
    opacity: 1;
    transform: translateY(0);
    pointer-events: auto;
  }

  #mainNav .navbar-nav {
    padding: 1rem 0;
  }

  #mainNav .nav-item {
    margin: 0.5rem 0;
    opacity: 0;
    transform: translateY(10px);
    transition: all 0.3s ease-in-out;
  }

  #mainNav .navbar-collapse.show .nav-item {
    opacity: 1;
    transform: translateY(0);
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

.navbar-toggler {
  border: none !important;
  padding: 0.5rem;
  color: #fff !important;
  transition: all 0.3s ease;
  outline: none !important;
  box-shadow: none !important;
  background: transparent !important;
}

.navbar-toggler:hover {
  color: rgba(255, 255, 255, 0.8) !important;
}

.navbar-toggler:focus {
  box-shadow: none !important;
  outline: none !important;
}

.navbar-toggler i {
  font-size: 1.5rem;
  transition: transform 0.3s ease;
}

.navbar-toggler:not(.collapsed) i {
  transform: rotate(180deg);
}

@media (max-width: 991.98px) {
  #mainNav .navbar-collapse {
    position: fixed;
    top: 70px;
    left: 0;
    right: 0;
    background: #000000;
    backdrop-filter: blur(10px);
    padding: 1rem;
    border-radius: 0 0 1rem 1rem;
    margin-top: 0;
    max-height: calc(100vh - 70px);
    overflow-y: auto;
    transition: all 0.3s ease-in-out;
    opacity: 0;
    transform: translateY(-10px);
    pointer-events: none;
  }

  #mainNav .navbar-collapse.show {
    opacity: 1;
    transform: translateY(0);
    pointer-events: auto;
  }

  #mainNav .navbar-nav {
    padding: 1rem 0;
  }

  #mainNav .nav-item {
    margin: 0.5rem 0;
    opacity: 0;
    transform: translateY(10px);
    transition: all 0.3s ease-in-out;
  }

  #mainNav .navbar-collapse.show .nav-item {
    opacity: 1;
    transform: translateY(0);
  }
}
</style> 