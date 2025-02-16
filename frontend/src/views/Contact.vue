<template>
  <div class="contact">
    <div class="contact-header">
      <h1>İLETİŞİM</h1>
      <p class="subtitle">Bize ulaşın, size yardımcı olalım.</p>
    </div>

    <div class="contact-content">
      <div class="contact-info">
        <h2>İletişim Bilgileri</h2>
        <div class="info-item">
          <i class="pi pi-map-marker"></i>
          <div>
            <strong>Adres:</strong>
            <p>{{ contact.address }}</p>
          </div>
        </div>
        <div class="info-item">
          <i class="pi pi-phone"></i>
          <div>
            <strong>Telefon:</strong>
            <p>{{ contact.phone }}</p>
          </div>
        </div>
        <div class="info-item">
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
            <p>Hafta içi: 09:00 - 18:00</p>
            <p>Cumartesi: 09:00 - 14:00</p>
            <p>Pazar: Kapalı</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import axios from 'axios';

const contact = ref({
  address: '',
  phone: '',
  email: ''
});

const fetchContact = async () => {
  try {
    const response = await axios.get('/api/contact');
    contact.value = response.data;
  } catch (error) {
    console.error('İletişim bilgileri yüklenirken hata:', error);
  }
};

onMounted(() => {
  fetchContact();
});
</script>

<style lang="scss" scoped>
.contact {
  max-width: 1200px;
  margin: 0 auto;
  padding: 2rem;

  .contact-header {
    text-align: center;
    margin-bottom: 3rem;

    h1 {
      font-size: 2.5rem;
      margin-bottom: 1rem;
    }

    .subtitle {
      color: #666;
      font-size: 1.2rem;
    }
  }

  .contact-content {
    margin-top: 2rem;
  }

  .contact-info {
    background: #fff;
    padding: 2rem;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0,0,0,0.1);

    h2 {
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
}
</style> 