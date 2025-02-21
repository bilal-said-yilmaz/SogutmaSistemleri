<template>
  <div class="home">
    <!-- Hero Section -->
    <section class="hero">
      <div class="hero-content">
        <h1>Kozan'ın Güvenilir Klima ve Beyaz Eşya Servisi</h1>
        <p>15 yılı aşkın tecrübemizle Kozan'da klima satış, montaj, bakım ve beyaz eşya tamiri hizmetleri sunuyoruz.</p>
        <div class="hero-buttons">
          <a href="#services" class="btn btn-primary">Hizmetlerimiz</a>
          <a href="#contact" class="btn btn-secondary">İletişime Geçin</a>
        </div>
      </div>
    </section>

    <!-- Services Section -->
    <section id="services" class="services">
      <div class="container">
        <h2>Hizmetlerimiz</h2>
        <div class="services-grid">
          <div class="service-card">
            <i class="fas fa-snowflake"></i>
            <h3>Klima Servisi Kozan</h3>
            <p>Tüm marka klimalarınız için bakım, onarım ve montaj hizmetleri. Kozan'da aynı gün servis garantisi.</p>
          </div>
          <div class="service-card">
            <i class="fas fa-tools"></i>
            <h3>Beyaz Eşya Tamiri Kozan</h3>
            <p>Buzdolabı, çamaşır makinesi, bulaşık makinesi tamiri. Kozan'da ev ve işyerlerine hızlı servis.</p>
          </div>
          <div class="service-card">
            <i class="fas fa-temperature-low"></i>
            <h3>Soğutma Sistemleri Kozan</h3>
            <p>Endüstriyel soğutma sistemleri kurulum ve bakımı. Kozan'da işletmelere özel çözümler.</p>
          </div>
        </div>
      </div>
    </section>

    <!-- Why Choose Us Section -->
    <section class="why-us">
      <div class="container">
        <h2>Neden Bizi Seçmelisiniz?</h2>
        <div class="features-grid">
          <div class="feature">
            <i class="fas fa-clock"></i>
            <h3>Aynı Gün Servis</h3>
            <p>Kozan'da aynı gün içinde servis hizmeti sağlıyoruz.</p>
          </div>
          <div class="feature">
            <i class="fas fa-certificate"></i>
            <h3>Profesyonel Ekip</h3>
            <p>Uzman ve deneyimli teknik kadromuzla kaliteli hizmet.</p>
          </div>
          <div class="feature">
            <i class="fas fa-shield-alt"></i>
            <h3>Garantili Hizmet</h3>
            <p>Tüm tamir ve bakım işlemlerimiz garantilidir.</p>
          </div>
          <div class="feature">
            <i class="fas fa-hand-holding-usd"></i>
            <h3>Uygun Fiyat</h3>
            <p>Kozan'da rekabetçi fiyatlarla kaliteli hizmet.</p>
          </div>
        </div>
      </div>
    </section>

    <!-- Service Areas Section -->
    <section class="service-areas">
      <div class="container">
        <h2>Hizmet Bölgelerimiz</h2>
        <div class="areas-content">
          <p>Kozan merkez ve tüm mahallelerinde hizmet veriyoruz:</p>
          <ul class="areas-list">
            <li>Kozan Merkez</li>
            <li>Mahmutlu Mahallesi</li>
            <li>Hacıuşağı Mahallesi</li>
            <li>Yarımoğlu Mahallesi</li>
            <li>Çanaklı Mahallesi</li>
            <li>Ve diğer tüm mahalleler</li>
          </ul>
        </div>
      </div>
    </section>

    <section id="products" class="section">
      <h2>Ürünlerimiz</h2>
      <div class="products-grid">
        <div v-for="product in products" :key="product.id" class="product-card">
          <img :src="product.image" :alt="product.name">
          <h3>{{ product.name }}</h3>
          <p>{{ product.description }}</p>
          <span class="price">{{ product.price }} TL</span>
        </div>
      </div>
    </section>

    <section id="about" class="section">
      <h2>Hakkımızda</h2>
      <div class="about-content" v-if="about">
        <img :src="about.image" :alt="about.title">
        <div class="about-text">
          <h3>{{ about.title }}</h3>
          <p>{{ about.content }}</p>
        </div>
      </div>
    </section>

    <section id="contact" class="section">
      <h2>İletişim</h2>
      <div class="contact-content">
        <div class="contact-info">
          <h3>İletişim Bilgileri</h3>
          <div class="info-item" v-if="contact">
            <i class="pi pi-map-marker"></i>
            <div>
              <strong>Adres:</strong>
              <p>{{ contact.address }}</p>
            </div>
          </div>
          <div class="info-item" v-if="contact">
            <i class="pi pi-phone"></i>
            <div>
              <strong>Telefon:</strong>
              <p>{{ contact.phone }}</p>
            </div>
          </div>
          <div class="info-item" v-if="contact">
            <i class="pi pi-envelope"></i>
            <div>
              <strong>E-posta:</strong>
              <p>{{ contact.email }}</p>
            </div>
          </div>
          <div class="info-item">
            <i class="pi pi-clock"></i>
            <div>
              <strong>Çalışma Saatleri:</strong>
              <p>Hafta içi: {{ contact?.weekdayHours || '09:00 - 18:00' }}</p>
              <p>Cumartesi: {{ contact?.saturdayHours || '09:00 - 14:00' }}</p>
              <p>Pazar: {{ contact?.sundayHours || 'Kapalı' }}</p>
            </div>
          </div>
        </div>

        <div class="map-container">
          <iframe
            width="100%"
            height="100%"
            frameborder="0"
            style="border:0"
            :src="`https://maps.google.com/maps?q=${contact?.latitude || '37.4538'},${contact?.longitude || '35.8158'}&z=15&output=embed`"
            allowfullscreen
          ></iframe>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useStore } from 'vuex'
import axios from 'axios'

const store = useStore()
const services = ref([])
const products = ref([])
const about = ref(null)
const contact = ref(null)

const fetchData = async () => {
  try {
    const [servicesRes, productsRes, aboutRes, contactRes] = await Promise.all([
      axios.get('/api/services'),
      axios.get('/api/products'),
      axios.get('/api/about'),
      axios.get('/api/contact')
    ])

    services.value = servicesRes.data
    products.value = productsRes.data
    about.value = aboutRes.data
    contact.value = contactRes.data
  } catch (error) {
    console.error('Veri yüklenirken hata oluştu:', error)
  }
}

onMounted(() => {
  fetchData()
})
</script>

<style lang="scss" scoped>
.section {
  padding: 4rem 2rem;
  scroll-margin-top: 60px;

  h2 {
    text-align: center;
    margin-bottom: 2rem;
    color: #333;
  }
}

.services-grid,
.products-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 2rem;
  padding: 1rem;
}

.service-card,
.product-card {
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
  overflow: hidden;
  transition: transform 0.3s;

  &:hover {
    transform: translateY(-5px);
  }

  img {
    width: 100%;
    height: 200px;
    object-fit: cover;
  }

  h3 {
    padding: 1rem;
    margin: 0;
    color: #333;
  }

  p {
    padding: 0 1rem 1rem;
    color: #666;
  }
}

.product-card {
  .price {
    display: block;
    padding: 0 1rem 1rem;
    color: #007bff;
    font-weight: bold;
  }
}

.about-content {
  display: grid;
  grid-template-columns: 1fr 2fr;
  gap: 2rem;
  align-items: center;

  img {
    width: 100%;
    border-radius: 8px;
  }

  .about-text {
    h3 {
      margin-top: 0;
      color: #333;
    }

    p {
      color: #666;
      line-height: 1.6;
    }
  }
}

@media (max-width: 768px) {
  .about-content {
    grid-template-columns: 1fr;
  }
}

.contact-content {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 2rem;

  @media (max-width: 768px) {
    grid-template-columns: 1fr;
  }
}

.contact-info {
  background: #fff;
  padding: 2rem;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);

  h3 {
    margin-bottom: 2rem;
    color: #333;
  }

  .info-item {
    display: flex;
    align-items: flex-start;
    margin-bottom: 1.5rem;

    i {
      font-size: 1.5rem;
      color: #007bff;
      margin-right: 1rem;
      margin-top: 0.25rem;
    }

    div {
      strong {
        display: block;
        margin-bottom: 0.5rem;
        color: #333;
      }

      p {
        color: #666;
        margin: 0;
        line-height: 1.4;
      }
    }
  }
}

.map-container {
  height: 400px;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}
</style> 