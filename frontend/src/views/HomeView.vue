<script setup lang="ts">
import { ref, onMounted, watch, nextTick, defineAsyncComponent } from 'vue'
import axios from 'axios'
import { useStore } from 'vuex'

// TypeScript interfaces
interface Product {
  id: number
  name: string
  description: string
  image: string
  category: string
  price?: number
  features?: string[]
}

interface Contact {
  title: string
  address: string
  phone: string
  email: string
  latitude: number
  longitude: number
  weekdayHours: string
  saturdayHours: string
  sundayHours: string
}

interface ProductForm {
  image: string
}

interface ServiceForm {
  image: string
}

// Services Data
const store = useStore()
const services = ref<Product[]>([])
const products = ref<Product[]>([])
const selectedProduct = ref<Product | null>(null)
const about = ref({
  title: 'Kozan İklimlendirme',
  content: '',
  image: '/placeholder-about.jpg'
})
const contact = ref<Contact | null>(null)

// Product modal state
const showProductModal = ref(false)

// Add video state management
const videoError = ref(false)
const handleVideoError = () => {
  videoError.value = true
}

// Lazy loaded components
const Dialog = defineAsyncComponent(() => import('primevue/dialog'))

// Intersection Observer için ref
const observerRef = ref(null)

// Form refs
const productForm = ref<ProductForm>({ image: '' })
const serviceForm = ref<ServiceForm>({ image: '' })

// Görsel önbelleğe alma
const preloadImages = () => {
  const imageUrls = [
    ...products.value.map(p => p.image),
    ...services.value.map(s => s.image),
    about.value?.image
  ].filter((url): url is string => url !== undefined)

  imageUrls.forEach(url => {
    const img = new Image()
    img.src = url
  })
}

// Fetch data
const fetchData = async () => {
  try {
    const [servicesRes, productsRes, aboutRes, contactRes] = await Promise.all([
      axios.get<Product[]>('/api/services'),
      axios.get<Product[]>('/api/products'),
      axios.get('/api/about'),
      axios.get<Contact>('/api/contact')
    ])

    services.value = servicesRes.data.map(service => ({
      ...service,
      image: service.image || '/placeholder-image.jpg'
    }))
    
    products.value = productsRes.data.map(product => ({
      ...product,
      image: product.image || '/placeholder-image.jpg'
    }))
    
    about.value = {
      ...aboutRes.data,
      image: aboutRes.data.image || '/placeholder-about.jpg'
    }
    
    contact.value = contactRes.data
    
    // Debug logları kaldırıldı
  } catch (error) {
    console.error('Veri yüklenirken hata oluştu:', error)
  }
}

// Helper functions
const truncateDescription = (text) => {
  if (!text) return ''
  return text.length > 100 ? text.substring(0, 100) + '...' : text
}

const formatPrice = (price) => {
  return new Intl.NumberFormat('tr-TR').format(price)
}

const openProductModal = (product) => {
  selectedProduct.value = product
}

const closeProductModal = (event) => {
  if (event.target.classList.contains('modal-overlay')) {
    selectedProduct.value = null
  }
}

onMounted(() => {
  fetchData()
  
  const observer = new IntersectionObserver(
    (entries) => {
      entries.forEach(entry => {
        if (entry.isIntersecting) {
          const img = entry.target as HTMLImageElement
          const dataSrc = img.getAttribute('data-src')
          if (dataSrc) {
            img.src = dataSrc
            img.removeAttribute('data-src')
          }
          observer.unobserve(img)
        }
      })
    },
    {
      rootMargin: '50px 0px',
      threshold: 0.1
    }
  )

  const images = document.querySelectorAll('img[data-src]')
  images.forEach(img => observer.observe(img))

  // Veri yüklendikten sonra görselleri önbelleğe al
  watch([products, services, about], () => {
    nextTick(() => {
      preloadImages()
    })
  })
})

// Performance metrics için düzeltme
const reportPerformance = () => {
  if ('performance' in window) {
    const metrics = performance.getEntriesByType('navigation')[0] as PerformanceNavigationTiming
    if (metrics) {
      // Performance metrikleri kaldırıldı
      const perfData = {
        dnsTime: metrics.domainLookupEnd - metrics.domainLookupStart,
        tcpTime: metrics.connectEnd - metrics.connectStart,
        firstByteTime: metrics.responseStart - metrics.requestStart,
        domLoadTime: metrics.domComplete - metrics.domInteractive,
        totalTime: metrics.loadEventEnd - metrics.startTime
      }
      // Metrikleri bir analytics servisine gönderebiliriz
    }
  }
}

onMounted(() => {
  window.addEventListener('load', reportPerformance)
  return () => {
    window.removeEventListener('load', reportPerformance)
  }
})

// Add function to generate Google Maps URL
const getGoogleMapsUrl = () => {
  return 'https://maps.app.goo.gl/sLDnB9FuE55URarD8'
}

// Add these methods to the script section
const openGoogleMaps = () => {
  window.open(getGoogleMapsUrl(), '_blank')
}

const makePhoneCall = () => {
  window.location.href = `tel:${contact.value?.phone}`
}

const sendEmail = () => {
  const email = contact.value?.email || 'kozaniklimlendirme@gmail.com'
  window.open(`https://mail.google.com/mail/?view=cm&fs=1&to=${email}`, '_blank')
}

// Add default phone number
const defaultPhone = '+90 (538) 515 3191'

// Add handlePhoneClick function
const handlePhoneClick = () => {
  const phoneNumber = contact.value?.phone || defaultPhone
  // Remove spaces and parentheses for the tel: link
  const formattedPhone = phoneNumber.replace(/[\s()]/g, '')
  window.location.href = `tel:${formattedPhone}`
}

// Görsel dönüştürme fonksiyonu düzeltmesi
const convertToWebP = async (file: File): Promise<string> => {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onload = (e: ProgressEvent<FileReader>) => {
      if (!e.target?.result) {
        reject(new Error('Failed to read file'))
        return
      }
      
      const img = new Image()
      img.onload = () => {
        const canvas = document.createElement('canvas')
        canvas.width = img.width
        canvas.height = img.height
        const ctx = canvas.getContext('2d')
        if (!ctx) {
          reject(new Error('Failed to get canvas context'))
          return
        }
        
        ctx.drawImage(img, 0, 0)
        canvas.toBlob((blob) => {
          if (!blob) {
            reject(new Error('Failed to create blob'))
            return
          }
          resolve(URL.createObjectURL(blob))
        }, 'image/webp', 0.8)
      }
      img.src = e.target.result as string
    }
    reader.onerror = () => reject(new Error('Failed to read file'))
    reader.readAsDataURL(file)
  })
}

// Form ve görsel işleme fonksiyonları düzeltmesi
const handleProductImage = async (event: Event) => {
  const target = event.target as HTMLInputElement
  const file = target.files?.[0]
  if (file) {
    try {
      const webpUrl = await convertToWebP(file)
      productForm.value.image = webpUrl
    } catch (error) {
      console.error('Görsel dönüştürme hatası:', error)
      const reader = new FileReader()
      reader.onload = (e: ProgressEvent<FileReader>) => {
        if (e.target?.result) {
          productForm.value.image = e.target.result as string
        }
      }
      reader.readAsDataURL(file)
    }
  }
}

const handleServiceImage = async (event: Event) => {
  const target = event.target as HTMLInputElement
  const file = target.files?.[0]
  if (file) {
    try {
      const webpUrl = await convertToWebP(file)
      serviceForm.value.image = webpUrl
    } catch (error) {
      console.error('Görsel dönüştürme hatası:', error)
      const reader = new FileReader()
      reader.onload = (e: ProgressEvent<FileReader>) => {
        if (e.target?.result) {
          serviceForm.value.image = e.target.result as string
        }
      }
      reader.readAsDataURL(file)
    }
  }
}

const handleAboutImage = (event: Event) => {
  const target = event.target as HTMLInputElement
  const file = target.files?.[0]
  if (file) {
    const reader = new FileReader()
    reader.onload = (e: ProgressEvent<FileReader>) => {
      if (e.target?.result) {
        about.value.image = e.target.result as string
      }
    }
    reader.readAsDataURL(file)
  }
}
</script>

<template>
  <div class="home">
    <!-- Hero Section -->
    <header class="masthead">
      <div class="container">
        <div class="masthead-content fade-in">
          <div class="masthead-subheading">Kozan'ın Güvenilir Klima ve Beyaz Eşya Servisi</div>
          <div class="masthead-heading text-uppercase">Profesyonel Klima, Beyaz Eşya Tamiri ve Soğutma Sistemleri Çözümleri</div>
          <a class="btn-discover" href="#services">
            <span>Hizmetlerimizi Keşfedin</span>
            <i class="fas fa-chevron-down"></i>
          </a>
        </div>
      </div>
      <div class="masthead-overlay"></div>
      <div class="masthead-background" :class="{ 'fallback-background': videoError }">
        <video 
          class="masthead-video" 
          autoplay 
          loop 
          muted 
          playsinline
          @error="handleVideoError"
        >
          <source src="../assets/videos/background.mp4" type="video/mp4">
        </video>
      </div>
    </header>

    <!-- Services Section -->
    <section class="page-section" id="services">
      <div class="container">
        <div class="text-center fade-in">
          <h1 class="section-heading text-uppercase">Klima ve Beyaz Eşya Servisi</h1>
          <h2 class="section-subheading text-muted">Uzman Ekip, Garantili Hizmet, Uygun Fiyat</h2>
        </div>
        <div class="row text-center">
          <div class="col-md-4 service-item fade-in" v-for="service in services" :key="service.id">
            <div class="service-image">
              <img 
                :src="service.image" 
                :alt="service.name + ' - Kozan İklimlendirme'" 
                loading="lazy"
                decoding="async"
                fetchpriority="low"
                width="300"
                height="200"
              />
            </div>
            <h3 class="service-heading">{{ service.name }}</h3>
            <p class="service-description">{{ service.description }}</p>
            <ul class="service-features" v-if="service.features">
              <li v-for="feature in service.features" :key="feature">{{ feature }}</li>
            </ul>
          </div>
        </div>
      </div>
    </section>

    <!-- Products Section -->
    <section class="page-section bg-light" id="products">
      <div class="container">
        <div class="text-center">
          <h2 class="section-heading text-uppercase">Klima ve Soğutma Sistemleri</h2>
          <h3 class="section-subheading text-muted">En İyi Markalar, En Uygun Fiyatlar</h3>
        </div>
        <div class="row">
          <div class="col-lg-4 col-sm-6 mb-4" v-for="product in products" :key="product.id">
            <article class="product-card" @click="openProductModal(product)">
              <div class="product-image">
                <img 
                  class="img-fluid" 
                  :src="product.image" 
                  :alt="product.name + ' - Kozan İklimlendirme'" 
                  loading="lazy"
                  decoding="async"
                  fetchpriority="low"
                  width="300"
                  height="200"
                />
                <div class="product-overlay">
                  <div class="overlay-content">
                    <i class="fas fa-search-plus"></i>
                    <span>Detayları Görüntüle</span>
                  </div>
                </div>
              </div>
              <div class="product-details">
                <h3 class="product-title">{{ product.name }}</h3>
                <p class="product-description">{{ truncateDescription(product.description) }}</p>
                <div class="product-price" itemprop="offers" itemscope itemtype="https://schema.org/Offer">
                  <meta itemprop="priceCurrency" content="TRY" />
                  <span itemprop="price">{{ formatPrice(product.price) }}</span> TL
                </div>
              </div>
            </article>
          </div>
        </div>
      </div>

      <!-- Ürün Detay Modalı -->
      <div v-if="selectedProduct" class="modal-overlay" @click="closeProductModal">
        <div class="modal-content" @click.stop>
          <button class="modal-close" @click="selectedProduct = null">
            <i class="fas fa-times"></i>
          </button>
          <div class="modal-body">
            <div class="product-modal-grid">
              <div class="product-image-container">
                <img 
                  :src="selectedProduct.image" 
                  :alt="selectedProduct.name" 
                  class="modal-product-image"
                  loading="lazy"
                  decoding="async"
                  width="400"
                  height="300"
                >
              </div>
              <div class="product-info-container">
                <h3 class="modal-title">{{ selectedProduct.name }}</h3>
                <div class="price-tag">
                  <span class="price">{{ formatPrice(selectedProduct.price) }}</span>
                  <span class="currency">TL</span>
                </div>
                <div class="product-description">
                  <p>{{ selectedProduct.description }}</p>
                </div>
                <div class="product-contact">
                  <p class="contact-info">
                    <i class="fas fa-info-circle"></i>
                    Ürün hakkında bilgi almak için:
                  </p>
                  <div class="contact-details">
                    <p><i class="fas fa-phone"></i>{{ contact?.phone || '+90 (538) 515 3191' }}</p>
                    <p><i class="fas fa-envelope"></i>{{ contact?.email || 'kozaniklimlendirme@gmail.com' }}</p>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- About Section -->
    <section class="page-section" id="about">
      <div class="container">
        <div class="text-center">
          <h2 class="section-heading text-uppercase">Hakkımızda</h2>
          <h3 class="section-subheading text-muted">Güvenilir hizmet, kaliteli ürünler.</h3>
        </div>
        <div class="about-content">
          <div class="about-image">
            <img 
              :src="about?.image || '/placeholder-about.jpg'"
              :alt="about?.title || 'Hakkımızda'"
              class="img-fluid"
              loading="lazy"
              decoding="async"
              fetchpriority="low"
              width="600"
              height="400"
            />
          </div>
          <div class="about-text">
            <h3>{{ about?.title }}</h3>
            <div class="content" v-html="about?.content"></div>
          </div>
        </div>
      </div>
    </section>

    <!-- Contact Section -->
    <section id="contact" class="section">
      <div class="container">
        <div class="text-center">
          <h2 class="section-heading text-uppercase">İletişim</h2>
          <h3 class="section-subheading text-muted">Bize ulaşın.</h3>
        </div>
        
        <div class="contact-content">
          <div class="info-card">
            <div class="info-items">
              <div class="info-item clickable" @click="openGoogleMaps">
                <i class="pi pi-map-marker"></i>
                <div>
                  <strong>Adres</strong>
                  <p>
                    {{ contact?.address }}
                    <span class="hover-text">
                      <i class="fas fa-location-arrow"></i>
                      Google Haritalarda göster
                    </span>
                  </p>
                </div>
              </div>
              
              <div class="info-item clickable" @click="handlePhoneClick">
                <i class="pi pi-phone"></i>
                <div>
                  <strong>Telefon:</strong>
                  <p>{{ contact?.phone || defaultPhone }}</p>
                </div>
              </div>
              
              <div class="info-item clickable" @click="sendEmail">
                <i class="pi pi-envelope"></i>
                <div>
                  <strong>E-posta:</strong>
                  <p>{{ contact?.email || 'kozaniklimlendirme@gmail.com' }} <i class="fas fa-external-link-alt"></i></p>
                </div>
              </div>
              
              <div class="info-item">
                <i class="pi pi-clock"></i>
                <div>
                  <strong>Çalışma Saatleri</strong>
                  <div class="working-hours">
                    <p>
                      <span>Hafta İçi:</span>
                      <span>{{ contact?.weekdayHours || 'Belirtilmemiş' }}</span>
                    </p>
                    <p>
                      <span>Cumartesi:</span>
                      <span>{{ contact?.saturdayHours || 'Belirtilmemiş' }}</span>
                    </p>
                    <p>
                      <span>Pazar:</span>
                      <span>{{ contact?.sundayHours || 'Belirtilmemiş' }}</span>
                    </p>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<style scoped>
.home {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

/* Hero Section */
.masthead {
  position: relative;
  width: 100%;
  height: 100vh;
  min-height: 600px;
  display: flex;
  align-items: center;
  overflow: hidden;
  background: none;
}

.masthead-content {
  position: relative;
  z-index: 2;
  max-width: 800px;
  margin: 0 auto;
  text-align: center;
  padding: 2rem;
}

.masthead-subheading {
  font-size: 1.5rem;
  color: #fff;
  margin-bottom: 1rem;
  text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.5);
  font-family: 'Montserrat', sans-serif;
  font-weight: 500;
}

.masthead-heading {
  font-size: 3rem;
  font-weight: 700;
  color: #fff;
  margin-bottom: 2rem;
  text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.5);
}

.masthead-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.4);
  z-index: 1;
}

.masthead-background {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 0;
}

.fallback-background {
  background: linear-gradient(135deg, rgba(30, 58, 138, 0.95), rgba(30, 144, 255, 0.9));
}

.masthead-video {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  min-width: 100%;
  min-height: 100%;
  width: auto;
  height: auto;
  object-fit: cover;
}

.btn-discover {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  padding: 1rem 2rem;
  background: #1e90ff;
  color: #fff;
  text-decoration: none;
  border-radius: 50px;
  font-weight: 600;
  transition: all 0.3s ease;
  text-transform: uppercase;
  letter-spacing: 1px;
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.btn-discover:hover {
  transform: translateY(-3px);
  box-shadow: 0 4px 15px rgba(30, 144, 255, 0.3);
  background: #0077e6;
}

.btn-discover i {
  font-size: 1.2rem;
}

@media (max-width: 768px) {
  .masthead-heading {
    font-size: 2rem;
  }
  
  .masthead-subheading {
    font-size: 1.2rem;
  }
}

/* Services Section */
.page-section {
  padding: 6rem 0;
}

.section-heading {
  font-size: 2.5rem;
  margin-top: 0;
  margin-bottom: 1.5rem;
  text-align: center !important;
}

.section-subheading {
  font-size: 1.1rem;
  font-weight: 400;
  font-style: italic;
  margin-bottom: 4rem;
  text-align: center !important;
  color: #666;
}

.service-item {
  margin-bottom: 3rem;
  transition: transform 0.3s ease;
  text-align: left;
}

.service-item:hover {
  transform: translateY(-5px);
}

.service-image {
  margin-bottom: 1.5rem;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
  aspect-ratio: 4/3;
  background-color: #f8f9fa;
}

.service-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s ease;
  display: block;
}

.service-item:hover .service-image img {
  transform: scale(1.05);
}

.service-heading {
  font-size: 1.5rem;
  margin-bottom: 1rem;
  color: var(--primary-color);
  text-align: left;
}

.service-description {
  color: #666;
  line-height: 1.8;
  text-align: left;
}

/* Contact Section Styles */
.contact-content {
  max-width: 1200px;
  margin: 0 auto;
}

.info-card {
  background: #fff;
  padding: 2rem;
  border-radius: 12px;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.1);
}

.info-items {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 1.5rem;
}

.info-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  padding: 1.5rem;
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.8), rgba(248, 249, 250, 0.9));
  border-radius: 10px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  border: 1px solid rgba(30, 144, 255, 0.1);
  height: 100%;
}

.info-item:hover {
  transform: translateY(-3px);
  box-shadow: 0 6px 20px rgba(30, 144, 255, 0.15);
  border-color: rgba(30, 144, 255, 0.2);
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.95), rgba(248, 249, 250, 1));
}

.info-item i {
  font-size: 2rem;
  color: #1e90ff;
  margin-bottom: 1rem;
}

.info-item strong {
  display: block;
  margin-bottom: 0.75rem;
  color: #1e3a8a;
  font-size: 1.2rem;
  font-weight: 600;
}

.info-item.clickable {
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.info-item.clickable:hover {
  transform: translateY(-3px);
  box-shadow: 0 6px 20px rgba(30, 144, 255, 0.15);
  border-color: rgba(30, 144, 255, 0.2);
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.95), rgba(248, 249, 250, 1));
}

.info-item p i {
  margin-left: 0.5rem;
  font-size: 0.9rem;
  color: #1e90ff;
  opacity: 0.8;
  transition: all 0.3s ease;
}

.info-item:hover p i {
  opacity: 1;
  transform: translateX(2px);
}

.working-hours {
  margin-top: 0.75rem;
}

.working-hours p {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-bottom: 0.75rem;
  color: #666;
}

.working-hours p span:first-child {
  font-weight: 600;
  color: #1e3a8a;
  margin-bottom: 0.25rem;
}

.working-hours p span:last-child {
  color: #1e90ff;
  font-weight: 500;
}

@media (max-width: 991px) {
  .info-items {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 576px) {
  .info-items {
    grid-template-columns: 1fr;
  }
  
  .info-item {
    padding: 1.25rem;
  }
}

.map {
  width: 100%;
  height: 400px;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.product-card {
  background: white;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
  transition: transform 0.3s ease, box-shadow 0.3s ease;
  cursor: pointer;
  height: 100%;
}

.product-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 4px 8px rgba(0,0,0,0.2);
}

.product-image {
  position: relative;
  overflow: hidden;
  aspect-ratio: 4/3;
}

.product-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.product-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(30, 144, 255, 0.7);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transition: opacity 0.3s ease;
}

.product-card:hover .product-overlay {
  opacity: 1;
}

.overlay-content {
  color: white;
  text-align: center;
}

.overlay-content i {
  font-size: 2rem;
  margin-bottom: 0.5rem;
}

.overlay-content span {
  display: block;
  font-size: 1rem;
}

.product-details {
  padding: 1.5rem;
  background: white;
}

.product-title {
  font-size: 1.25rem;
  margin-bottom: 0.75rem;
  color: #1e3a8a;
  font-weight: 700;
  display: block;
  overflow: hidden;
  text-overflow: ellipsis;
  min-height: 3rem;
  line-height: 1.4;
  word-break: break-word;
  max-width: 100%;
  background: white;
  padding: 0.5rem;
  border-radius: 4px;
}

.product-description {
  color: #666;
  font-size: 1rem;
  margin-bottom: 1rem;
  line-height: 1.8;
  text-align: left;
}

.product-price {
  font-size: 1.5rem;
  font-weight: bold;
  color: var(--primary-color);
}

/* Modern Modal Styles */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.75);
  backdrop-filter: blur(5px);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
  padding: 0;
}

.modal-content {
  background-color: #fff;
  border-radius: 12px;
  max-width: 800px;
  width: 100%;
  max-height: 85vh;
  overflow-y: auto;
  position: relative;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25);
  margin: 1rem;
}

.modal-close {
  position: absolute;
  top: 0.5rem;
  right: 0.5rem;
  background: white;
  border: none;
  width: 28px;
  height: 28px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  z-index: 10;
  transition: all 0.3s ease;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}

.modal-close:hover {
  transform: scale(1.1);
}

.modal-close i {
  font-size: 1rem;
  color: #333;
}

.modal-body {
  padding: 0;
}

.product-modal-grid {
  display: grid;
  grid-template-columns: 1fr;
  min-height: auto;
}

.product-image-container {
  position: relative;
  background: #f8f9fa;
  padding: 1rem;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 12px 12px 0 0;
  min-height: 180px;
}

.modal-product-image {
  max-width: 100%;
  height: auto;
  max-height: 200px;
  object-fit: contain;
  border-radius: 8px;
}

.product-info-container {
  padding: 1.5rem;
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.modal-title {
  font-size: 1.25rem;
  font-weight: 600;
  color: #1e3a8a;
  margin-bottom: 0.75rem;
  line-height: 1.3;
}

.price-tag {
  margin-bottom: 1rem;
  display: flex;
  align-items: baseline;
  gap: 0.5rem;
}

.price {
  font-size: 1.5rem;
  font-weight: 600;
  color: #1e90ff;
}

.currency {
  font-size: 1rem;
  color: #666;
  font-weight: 500;
}

.product-description {
  color: #4a5568;
  font-size: 0.9rem;
  line-height: 1.6;
  margin-bottom: 1rem;
}

.product-contact {
  margin-top: 0.5rem;
  padding-top: 1rem;
  border-top: 1px solid #eee;
}

.contact-info {
  font-size: 0.9rem;
  color: #2d3748;
  margin-bottom: 0.75rem;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.contact-details {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.contact-details p {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  color: #4a5568;
  font-size: 0.9rem;
}

@media (min-width: 768px) {
  .product-modal-grid {
    grid-template-columns: 1fr 1fr;
  }

  .product-image-container {
    border-radius: 12px 0 0 12px;
    min-height: 300px;
    padding: 1.5rem;
  }

  .modal-product-image {
    max-height: 250px;
  }

  .product-info-container {
    padding: 2rem;
  }

  .modal-title {
    font-size: 1.5rem;
  }
}

@media (max-width: 576px) {
  .modal-content {
    margin: 0.5rem;
    max-height: 90vh;
  }

  .product-image-container {
    min-height: 150px;
  }

  .modal-product-image {
    max-height: 180px;
  }

  .product-info-container {
    padding: 1rem;
  }

  .modal-title {
    font-size: 1.1rem;
  }

  .price {
    font-size: 1.25rem;
  }

  .product-description {
    font-size: 0.85rem;
  }
}

/* About Section */
.about-content {
  display: flex;
  gap: 4rem;
  align-items: center;
  margin-top: 3rem;
}

.about-image {
  flex: 1;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
  position: relative;
  background-color: #f8f9fa;
  min-height: 400px;
}

.about-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  position: absolute;
  top: 0;
  left: 0;
  transition: transform 0.6s ease;
  display: block;
}

.about-image::before {
  content: '';
  display: block;
  padding-top: 66.67%; /* 3:2 aspect ratio */
}

.about-text {
  flex: 1;
  padding: 2rem;
  background: rgba(255, 255, 255, 0.8);
  border-radius: 12px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.05);
}

.about-text h3 {
  color: #1e3a8a;
  margin-bottom: 1.5rem;
  font-size: 1.8rem;
}

.about-text .content {
  font-size: 1.1rem;
  line-height: 1.8;
  color: #444;
}

.about-text .content p {
  margin-bottom: 1rem;
}

.about-text .content p:last-child {
  margin-bottom: 0;
}

@media (max-width: 991px) {
  .about-content {
    flex-direction: column;
    gap: 2rem;
  }

  .about-image {
    width: 100%;
    min-height: 300px;
  }
}

/* Contact Section */
#contact .section-heading {
  text-align: center !important;
}

#contact .section-subheading {
  text-align: center !important;
}

/* Products Section */
#products .section-heading {
  text-align: center !important;
}

#products .section-subheading {
  text-align: center !important;
}

/* Footer Styles */
.footer {
  background-color: rgba(33, 37, 41, 0.98);
  color: #fff;
  padding: 4rem 0 2rem;
  margin-top: auto;
  width: 100%;
  backdrop-filter: blur(10px);
  box-shadow: 0 -2px 10px rgba(0, 0, 0, 0.1);
}

.footer-content {
  max-width: 1200px;
  margin: 0 auto;
}

.footer-sections {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 2rem;
  margin-bottom: 3rem;
}

.footer-section h3 {
  font-size: 1.5rem;
  margin-bottom: 1.5rem;
  color: #1e90ff;
  font-weight: 600;
}

.footer-section ul {
  list-style: none;
  padding: 0;
  margin: 0;
}

.footer-section ul li {
  margin-bottom: 1rem;
}

.footer-section ul li a {
  color: rgba(255, 255, 255, 0.8);
  text-decoration: none;
  transition: color 0.3s ease;
}

.footer-section ul li a:hover {
  color: #1e90ff;
}

.footer-section ul li i {
  margin-right: 0.5rem;
  color: #1e90ff;
}

.footer-section p {
  color: rgba(255, 255, 255, 0.8);
  line-height: 1.6;
}

.footer-bottom {
  display: flex;
  justify-content: space-between;
  align-items: center;
  color: rgba(255, 255, 255, 0.6);
  font-size: 0.9rem;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
  padding-top: 2rem;
}

.sitemap a {
  color: rgba(255, 255, 255, 0.6);
  text-decoration: none;
  transition: color 0.3s ease;
}

.sitemap a:hover {
  color: #1e90ff;
}

@media (max-width: 768px) {
  .footer-sections {
    grid-template-columns: 1fr;
  }

  .footer-bottom {
    flex-direction: column;
    text-align: center;
    gap: 1rem;
  }
}

@media (max-width: 768px) {
  .masthead-heading {
    font-size: 2.5rem;
    letter-spacing: 0.03em;
    line-height: 1.3;
  }
}

@media (max-width: 576px) {
  .masthead-heading {
    font-size: 2rem;
    letter-spacing: 0.02em;
  }
}

.info-item p .hover-text {
  position: absolute;
  bottom: -30px;
  left: 50%;
  transform: translateX(-50%);
  background: rgba(30, 144, 255, 0.9);
  color: white;
  padding: 0.5rem 1rem;
  border-radius: 20px;
  font-size: 0.9rem;
  opacity: 0;
  transition: all 0.3s ease;
  white-space: nowrap;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}

.info-item.clickable:hover p .hover-text {
  opacity: 1;
  bottom: -40px;
}

.info-item {
  position: relative;
}

.info-item p i {
  display: none;
}

.service-image img,
.product-image img,
.modal-product-image,
.about-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transform: translateZ(0);
  will-change: transform;
  backface-visibility: hidden;
}

@supports (content-visibility: auto) {
  .service-image,
  .product-image,
  .about-image {
    content-visibility: auto;
    contain-intrinsic-size: 300px;
  }
}

/* Performans optimizasyonları için ek stiller */
.service-image,
.product-image,
.about-image {
  background-color: #f8f9fa;
  position: relative;
  overflow: hidden;
}

.service-image::before,
.product-image::before,
.about-image::before {
  content: "";
  display: block;
  padding-top: 66.67%; /* 3:2 aspect ratio */
}

.service-image img,
.product-image img,
.about-image img {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  object-fit: cover;
  transform: translateZ(0);
  will-change: transform;
  backface-visibility: hidden;
  transition: opacity 0.3s ease;
}

img[data-src] {
  opacity: 0;
}

img:not([data-src]) {
  opacity: 1;
}

/* Modern loading animasyonu */
@keyframes shimmer {
  0% {
    background-position: -200% 0;
  }
  100% {
    background-position: 200% 0;
  }
}

.service-image::after,
.product-image::after,
.about-image::after {
  content: "";
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: linear-gradient(
    90deg,
    rgba(255, 255, 255, 0) 0,
    rgba(255, 255, 255, 0.2) 20%,
    rgba(255, 255, 255, 0.5) 60%,
    rgba(255, 255, 255, 0)
  );
  background-size: 200% 100%;
  animation: shimmer 2s infinite;
  opacity: 0;
  pointer-events: none;
}

.service-image:not(:has(img[src])):after,
.product-image:not(:has(img[src])):after,
.about-image:not(:has(img[src])):after {
  opacity: 1;
}

/* Önbelleğe alma için stil */
.preload-image {
  position: absolute;
  width: 1px;
  height: 1px;
  opacity: 0;
  pointer-events: none;
}
</style> 