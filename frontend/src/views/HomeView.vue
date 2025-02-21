<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import axios from 'axios'
import Dialog from 'primevue/dialog'

interface Product {
  id: number
  title: string
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

// Services Data
const services = ref<Product[]>([])
const products = ref<Product[]>([])
const about = ref<{ title: string; content: string; image: string } | null>(null)
const contact = ref<Contact | null>(null)

// Contact Form
const contactForm = ref({
  name: '',
  email: '',
  phone: '',
  message: ''
})
const isSubmitting = ref(false)

// Product modal state
const showProductModal = ref(false)
const selectedProduct = ref(null)

// Add video state management
const videoError = ref(false)
const handleVideoError = () => {
  videoError.value = true
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

    services.value = servicesRes.data
    products.value = productsRes.data
    about.value = aboutRes.data
    contact.value = contactRes.data
    
    // Debug contact data
    console.log('Contact Data:', contactRes.data)
  } catch (error) {
    console.error('Veri yüklenirken hata oluştu:', error)
  }
}

// Handle Contact Form Submit
const handleContactSubmit = async () => {
  isSubmitting.value = true
  try {
    await axios.post('/api/contact', contactForm.value)
    alert('Mesajınız başarıyla gönderildi!')
    contactForm.value = {
      name: '',
      email: '',
      phone: '',
      message: ''
    }
  } catch (error) {
    console.error('Error submitting contact form:', error)
    alert('Mesaj gönderilirken bir hata oluştu. Lütfen tekrar deneyin.')
  } finally {
    isSubmitting.value = false
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
  showProductModal.value = true
}

onMounted(() => {
  fetchData()
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
  window.location.href = `mailto:${contact.value?.email}`
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
              <img :src="service.image" :alt="service.title + ' - Kozan İklimlendirme'" loading="lazy" />
            </div>
            <h3 class="service-heading">{{ service.title }}</h3>
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
                <img class="img-fluid" :src="product.image" :alt="product.name + ' - Kozan İklimlendirme'" loading="lazy" />
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

      <!-- Product Detail Modal -->
      <Dialog 
        v-model:visible="showProductModal" 
        :modal="true" 
        :dismissableMask="true"
        :style="{ width: '90%', maxWidth: '800px' }"
        class="product-modal"
      >
        <template #header>
          <h3 class="modal-title">{{ selectedProduct?.name }}</h3>
        </template>
        <template #default>
          <div class="product-modal-content" v-if="selectedProduct">
            <div class="row">
              <div class="col-md-6">
                <img :src="selectedProduct.image" :alt="selectedProduct.name" class="modal-product-image">
              </div>
              <div class="col-md-6">
                <div class="product-info">
                  <h4 class="price">{{ formatPrice(selectedProduct.price) }} TL</h4>
                  <p class="description">{{ selectedProduct.description }}</p>
                  <div class="contact-info">
                    <h5>Sipariş ve Bilgi için:</h5>
                    <div v-if="contact">
                      <p>
                        <a :href="'tel:' + contact.phone" class="contact-link">
                          <i class="pi pi-phone"></i> {{ contact.phone }}
                        </a>
                      </p>
                      <p>
                        <a :href="'mailto:' + contact.email" class="contact-link">
                          <i class="pi pi-envelope"></i> {{ contact.email }}
                        </a>
                      </p>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </template>
      </Dialog>
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
            <img :src="about?.image" :alt="about?.title" class="img-fluid">
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
                    <a :href="getGoogleMapsUrl()" target="_blank" class="address-link">
                      {{ contact?.address }}
                      <i class="pi pi-external-link"></i>
                    </a>
                  </p>
                </div>
              </div>
              
              <div class="info-item clickable" @click="makePhoneCall">
                <i class="pi pi-phone"></i>
                <div>
                  <strong>Telefon</strong>
                  <p>
                    <a :href="'tel:' + contact?.phone" class="contact-link">{{ contact?.phone }}</a>
                  </p>
                </div>
              </div>
              
              <div class="info-item clickable" @click="sendEmail">
                <i class="pi pi-envelope"></i>
                <div>
                  <strong>E-posta</strong>
                  <p>
                    <a :href="'mailto:' + contact?.email" class="contact-link">{{ contact?.email }}</a>
                  </p>
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

    <!-- Footer Section -->
    <footer class="footer">
      <div class="container">
        <div class="footer-content">
          <div class="footer-sections">
            <div class="footer-section">
              <h3>Hakkımızda</h3>
              <p>Kozan İklimlendirme olarak, 15 yılı aşkın tecrübemizle Kozan'da klima ve beyaz eşya servisi hizmeti vermekteyiz.</p>
            </div>
            <div class="footer-section">
              <h3>Hızlı Erişim</h3>
              <ul>
                <li><a href="#services">Hizmetlerimiz</a></li>
                <li><a href="#products">Ürünler</a></li>
                <li><a href="#about">Hakkımızda</a></li>
                <li><a href="#contact">İletişim</a></li>
              </ul>
            </div>
            <div class="footer-section">
              <h3>İletişim</h3>
              <ul>
                <li><i class="pi pi-phone"></i> {{ contact?.phone }}</li>
                <li><i class="pi pi-envelope"></i> {{ contact?.email }}</li>
                <li><i class="pi pi-map-marker"></i> {{ contact?.address }}</li>
              </ul>
            </div>
          </div>

          <div class="footer-center">
            <p>&copy; {{ new Date().getFullYear() }} Kozan İklimlendirme . Tüm hakları saklıdır .</p>
          </div>
        </div>
      </div>
    </footer>
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
}

.service-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s ease;
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

.contact-link {
  color: #1e90ff;
  text-decoration: none;
  transition: all 0.3s ease;
  font-weight: 500;
}

.contact-link:hover {
  color: #0077e6;
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
}

.product-title {
  font-size: 1.25rem;
  margin-bottom: 0.75rem;
  color: var(--primary-color);
  text-align: left;
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

/* Modal Styles */
.product-modal {
  .modal-title {
    margin: 0;
    color: var(--primary-color);
  }
}

.modal-product-image {
  width: 100%;
  border-radius: 8px;
  margin-bottom: 1rem;
}

.product-info {
  .price {
    font-size: 2rem;
    color: var(--primary-color);
    margin-bottom: 1rem;
  }

  .description {
    font-size: 1rem;
    line-height: 1.6;
    margin-bottom: 2rem;
  }

  .contact-info {
    background: #f8f9fa;
    padding: 1rem;
    border-radius: 8px;

    h5 {
      color: var(--primary-color);
      margin-bottom: 1rem;
    }

    p {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      margin-bottom: 0.5rem;

      i {
        color: var(--primary-color);
      }
    }
  }
}

/* Remove map styles and add address link styles */
.address-link {
  color: var(--primary-color);
  text-decoration: none;
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  transition: color 0.3s ease;
}

.address-link:hover {
  color: var(--primary-color-dark);
}

.address-link i {
  font-size: 0.9rem;
}

/* Remove contact-map class */

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
}

.about-image::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(135deg, rgba(30, 144, 255, 0.1), rgba(255, 107, 43, 0.1));
  z-index: 1;
  pointer-events: none;
}

.about-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  aspect-ratio: 4/3;
  transition: transform 0.6s ease;
}

.about-image:hover img {
  transform: scale(1.05);
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

.about-text .content p:last-child {
  margin-bottom: 0;
}

@media (max-width: 991px) {
  .about-content {
    flex-direction: column;
    gap: 2rem;
  }

  .about-image, .about-text {
    width: 100%;
  }

  .about-text {
    padding: 1.5rem;
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

/* Add these styles to the style section */
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
</style> 