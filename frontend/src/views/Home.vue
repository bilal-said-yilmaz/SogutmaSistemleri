<template>
  <div class="home">
    <section id="services" class="section">
      <h2>Hizmetlerimiz</h2>
      <div class="services-grid">
        <div v-for="service in services" :key="service.id" class="service-card">
          <img :src="service.image" :alt="service.title">
          <h3>{{ service.title }}</h3>
          <p>{{ service.description }}</p>
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
      <div class="contact-content" v-if="contact">
        <div class="contact-info">
          <h3>{{ contact.title }}</h3>
          <p><i class="pi pi-phone"></i> {{ contact.phone }}</p>
          <p><i class="pi pi-envelope"></i> {{ contact.email }}</p>
          <p><i class="pi pi-map-marker"></i> {{ contact.address }}</p>
        </div>
        <div class="contact-form">
          <form @submit.prevent="submitForm">
            <input v-model="form.name" type="text" placeholder="Adınız" required>
            <input v-model="form.email" type="email" placeholder="E-posta" required>
            <textarea v-model="form.message" placeholder="Mesajınız" required></textarea>
            <button type="submit">Gönder</button>
          </form>
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
const form = ref({
  name: '',
  email: '',
  message: ''
})

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

const submitForm = async () => {
  try {
    await axios.post('/api/contact/submit', form.value)
    alert('Mesajınız başarıyla gönderildi!')
    form.value = { name: '', email: '', message: '' }
  } catch (error) {
    alert('Mesaj gönderilirken bir hata oluştu.')
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

.contact-content {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 2rem;

  .contact-info {
    h3 {
      margin-top: 0;
      color: #333;
    }

    p {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      color: #666;
      margin: 1rem 0;

      i {
        color: #007bff;
      }
    }
  }

  .contact-form {
    form {
      display: flex;
      flex-direction: column;
      gap: 1rem;

      input,
      textarea {
        padding: 0.8rem;
        border: 1px solid #ddd;
        border-radius: 4px;
        font-size: 1rem;

        &:focus {
          outline: none;
          border-color: #007bff;
        }
      }

      textarea {
        min-height: 150px;
        resize: vertical;
      }

      button {
        padding: 0.8rem;
        background: #007bff;
        color: #fff;
        border: none;
        border-radius: 4px;
        font-size: 1rem;
        cursor: pointer;
        transition: background 0.3s;

        &:hover {
          background: #0056b3;
        }
      }
    }
  }
}

@media (max-width: 768px) {
  .about-content,
  .contact-content {
    grid-template-columns: 1fr;
  }
}
</style> 