<template>
  <div class="admin">
    <div v-if="!isAuthenticated" class="login-form">
      <h2>Admin Girişi</h2>
      <form @submit.prevent="login">
        <input v-model="loginForm.username" type="text" placeholder="Kullanıcı Adı" required>
        <input v-model="loginForm.password" type="password" placeholder="Şifre" required>
        <button type="submit">Giriş Yap</button>
      </form>
    </div>

    <div v-else class="admin-panel">
      <nav class="admin-nav">
        <button 
          v-for="tab in tabs" 
          :key="tab.id"
          :class="{ active: currentTab === tab.id }"
          @click="currentTab = tab.id"
        >
          {{ tab.name }}
        </button>
      </nav>

      <!-- Ürünler -->
      <div v-if="currentTab === 'products'" class="tab-content">
        <h2>Ürünler</h2>
        <button @click="showAddProductModal = true" class="add-button">
          Yeni Ürün Ekle
        </button>
        
        <div class="items-grid">
          <div v-for="product in products" :key="product.id" class="item-card">
            <img :src="product.image" :alt="product.name">
            <div class="item-info">
              <h3>{{ product.name }}</h3>
              <p>{{ product.description }}</p>
              <span class="price">{{ product.price }} TL</span>
            </div>
            <div class="item-actions">
              <button @click="editProduct(product)">Düzenle</button>
              <button @click="deleteProduct(product.id)" class="delete">Sil</button>
            </div>
          </div>
        </div>
      </div>

      <!-- Hizmetler -->
      <div v-if="currentTab === 'services'" class="tab-content">
        <h2>Hizmetler</h2>
        <button @click="showAddServiceModal = true" class="add-button">
          Yeni Hizmet Ekle
        </button>
        
        <div class="items-grid">
          <div v-for="service in services" :key="service.id" class="item-card">
            <img :src="service.image" :alt="service.title">
            <div class="item-info">
              <h3>{{ service.title }}</h3>
              <p>{{ service.description }}</p>
            </div>
            <div class="item-actions">
              <button @click="editService(service)">Düzenle</button>
              <button @click="deleteService(service.id)" class="delete">Sil</button>
            </div>
          </div>
        </div>
      </div>

      <!-- Hakkımızda -->
      <div v-if="currentTab === 'about'" class="tab-content">
        <h2>Hakkımızda</h2>
        <form @submit.prevent="updateAbout" class="edit-form">
          <input v-model="about.title" type="text" placeholder="Başlık">
          <textarea v-model="about.content" placeholder="İçerik"></textarea>
          <div class="image-upload">
            <input type="file" @change="handleAboutImage" accept="image/*">
            <img v-if="about.image" :src="about.image" class="preview-image" alt="Hakkımızda resmi">
          </div>
          <button type="submit">Güncelle</button>
        </form>
      </div>

      <!-- İletişim -->
      <div v-if="currentTab === 'contact'" class="tab-content">
        <h2>İletişim Bilgileri</h2>
        <form @submit.prevent="updateContact" class="edit-form">
          <input v-model="contact.title" type="text" placeholder="Başlık">
          <input v-model="contact.phone" type="text" placeholder="Telefon">
          <input v-model="contact.email" type="email" placeholder="E-posta">
          <textarea v-model="contact.address" placeholder="Adres"></textarea>
          <button type="submit">Güncelle</button>
        </form>
      </div>
    </div>

    <!-- Ürün Ekleme/Düzenleme Modal -->
    <Dialog v-model:visible="showAddProductModal" :header="editingProduct ? 'Ürün Düzenle' : 'Yeni Ürün Ekle'">
      <form @submit.prevent="saveProduct" class="modal-form">
        <input v-model="productForm.name" type="text" placeholder="Ürün Adı" required>
        <textarea v-model="productForm.description" placeholder="Açıklama" required></textarea>
        <input v-model="productForm.price" type="number" placeholder="Fiyat" required>
        <div class="image-upload">
          <input type="file" @change="handleProductImage" accept="image/*" :required="!productForm.image">
          <img v-if="productForm.image" :src="productForm.image" class="preview-image">
        </div>
        <button type="submit">{{ editingProduct ? 'Güncelle' : 'Ekle' }}</button>
      </form>
    </Dialog>

    <!-- Hizmet Ekleme/Düzenleme Modal -->
    <Dialog v-model:visible="showAddServiceModal" :header="editingService ? 'Hizmet Düzenle' : 'Yeni Hizmet Ekle'">
      <form @submit.prevent="saveService" class="modal-form">
        <input v-model="serviceForm.title" type="text" placeholder="Hizmet Başlığı" required>
        <textarea v-model="serviceForm.description" placeholder="Açıklama" required></textarea>
        <div class="image-upload">
          <input type="file" @change="handleServiceImage" accept="image/*" :required="!serviceForm.image">
          <img v-if="serviceForm.image" :src="serviceForm.image" class="preview-image">
        </div>
        <button type="submit">{{ editingService ? 'Güncelle' : 'Ekle' }}</button>
      </form>
    </Dialog>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useStore } from 'vuex'
import axios from 'axios'
import Dialog from 'primevue/dialog'

const store = useStore()
const isAuthenticated = computed(() => store.state.isAuthenticated)

const currentTab = ref('products')
const tabs = [
  { id: 'products', name: 'Ürünler' },
  { id: 'services', name: 'Hizmetler' },
  { id: 'about', name: 'Hakkımızda' },
  { id: 'contact', name: 'İletişim' }
]

const loginForm = ref({
  username: '',
  password: ''
})

const products = ref([])
const services = ref([])
const about = ref({})
const contact = ref({})

const showAddProductModal = ref(false)
const showAddServiceModal = ref(false)
const editingProduct = ref(null)
const editingService = ref(null)

const productForm = ref({
  name: '',
  description: '',
  price: '',
  image: ''
})

const serviceForm = ref({
  title: '',
  description: '',
  image: ''
})

// Login işlemi
const login = async () => {
  try {
    const response = await axios.post('/api/auth/login', loginForm.value)
    store.commit('setAuth', true)
    store.commit('setUser', response.data)
    loginForm.value = { username: '', password: '' }
  } catch (error) {
    alert('Giriş başarısız!')
  }
}

// Veri yükleme işlemleri
const fetchData = async () => {
  try {
    const [productsRes, servicesRes, aboutRes, contactRes] = await Promise.all([
      axios.get('/api/admin/products'),
      axios.get('/api/admin/services'),
      axios.get('/api/admin/about'),
      axios.get('/api/admin/contact')
    ])

    products.value = productsRes.data
    services.value = servicesRes.data
    about.value = aboutRes.data
    contact.value = contactRes.data
  } catch (error) {
    console.error('Veri yüklenirken hata oluştu:', error)
  }
}

// Ürün işlemleri
const saveProduct = async () => {
  try {
    if (editingProduct.value) {
      await axios.put(`/api/admin/products/${editingProduct.value.id}`, productForm.value)
    } else {
      await axios.post('/api/admin/products', productForm.value)
    }
    await fetchData()
    showAddProductModal.value = false
    productForm.value = { name: '', description: '', price: '', image: '' }
    editingProduct.value = null
  } catch (error) {
    alert('İşlem başarısız!')
  }
}

const editProduct = (product) => {
  editingProduct.value = product
  productForm.value = { ...product }
  showAddProductModal.value = true
}

const deleteProduct = async (id) => {
  if (confirm('Bu ürünü silmek istediğinizden emin misiniz?')) {
    try {
      await axios.delete(`/api/admin/products/${id}`)
      await fetchData()
    } catch (error) {
      alert('Silme işlemi başarısız!')
    }
  }
}

// Hizmet işlemleri
const saveService = async () => {
  try {
    if (editingService.value) {
      await axios.put(`/api/admin/services/${editingService.value.id}`, serviceForm.value)
    } else {
      await axios.post('/api/admin/services', serviceForm.value)
    }
    await fetchData()
    showAddServiceModal.value = false
    serviceForm.value = { title: '', description: '', image: '' }
    editingService.value = null
  } catch (error) {
    alert('İşlem başarısız!')
  }
}

const editService = (service) => {
  editingService.value = service
  serviceForm.value = { ...service }
  showAddServiceModal.value = true
}

const deleteService = async (id) => {
  if (confirm('Bu hizmeti silmek istediğinizden emin misiniz?')) {
    try {
      await axios.delete(`/api/admin/services/${id}`)
      await fetchData()
    } catch (error) {
      alert('Silme işlemi başarısız!')
    }
  }
}

// Hakkımızda güncelleme
const updateAbout = async () => {
  try {
    await axios.put('/api/admin/about', about.value)
    alert('Hakkımızda bilgileri güncellendi!')
  } catch (error) {
    alert('Güncelleme başarısız!')
  }
}

// İletişim güncelleme
const updateContact = async () => {
  try {
    await axios.put('/api/admin/contact', contact.value)
    alert('İletişim bilgileri güncellendi!')
  } catch (error) {
    alert('Güncelleme başarısız!')
  }
}

// Resim yükleme işlemleri
const uploadImage = async (file) => {
  const formData = new FormData()
  formData.append('file', file)

  try {
    const token = localStorage.getItem('token')
    if (!token) {
      throw new Error('Token bulunamadı')
    }

    const response = await axios.post('/api/admin/upload', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
        'Authorization': `Bearer ${token}`
      }
    })
    return response.data.url
  } catch (error) {
    console.error('Resim yükleme hatası:', error)
    alert('Resim yüklenirken hata oluştu!')
    throw error
  }
}

const handleProductImage = async (event) => {
  const file = event.target.files[0]
  if (file) {
    try {
      const url = await uploadImage(file)
      productForm.value.image = url
    } catch (error) {
      console.error('Resim yükleme hatası:', error)
    }
  }
}

const handleServiceImage = async (event) => {
  const file = event.target.files[0]
  if (file) {
    try {
      const url = await uploadImage(file)
      serviceForm.value.image = url
    } catch (error) {
      console.error('Resim yükleme hatası:', error)
    }
  }
}

const handleAboutImage = async (event) => {
  const file = event.target.files[0]
  if (file) {
    try {
      const url = await uploadImage(file)
      about.value.image = url
    } catch (error) {
      console.error('Resim yükleme hatası:', error)
    }
  }
}

// Sayfa yüklendiğinde verileri getir
if (isAuthenticated.value) {
  fetchData()
}
</script>

<style lang="scss" scoped>
.admin {
  max-width: 1200px;
  margin: 0 auto;
  padding: 2rem;
}

.login-form {
  max-width: 400px;
  margin: 2rem auto;
  padding: 2rem;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);

  h2 {
    text-align: center;
    margin-bottom: 2rem;
  }

  form {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }
}

.admin-nav {
  display: flex;
  gap: 1rem;
  margin-bottom: 2rem;

  button {
    padding: 0.8rem 1.5rem;
    border: none;
    background: #f0f0f0;
    border-radius: 4px;
    cursor: pointer;
    transition: background 0.3s;

    &.active {
      background: #007bff;
      color: #fff;
    }
  }
}

.tab-content {
  background: #fff;
  padding: 2rem;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);

  h2 {
    margin-bottom: 2rem;
  }
}

.add-button {
  margin-bottom: 2rem;
  padding: 0.8rem 1.5rem;
  background: #28a745;
  color: #fff;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: background 0.3s;

  &:hover {
    background: #218838;
  }
}

.items-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 2rem;
}

.item-card {
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
  overflow: hidden;

  img {
    width: 100%;
    height: 200px;
    object-fit: cover;
  }

  .item-info {
    padding: 1rem;

    h3 {
      margin: 0 0 0.5rem;
    }

    p {
      color: #666;
      margin: 0 0 0.5rem;
    }

    .price {
      color: #007bff;
      font-weight: bold;
    }
  }

  .item-actions {
    padding: 1rem;
    display: flex;
    gap: 1rem;

    button {
      flex: 1;
      padding: 0.5rem;
      border: none;
      border-radius: 4px;
      cursor: pointer;
      transition: background 0.3s;

      &:first-child {
        background: #007bff;
        color: #fff;

        &:hover {
          background: #0056b3;
        }
      }

      &.delete {
        background: #dc3545;
        color: #fff;

        &:hover {
          background: #c82333;
        }
      }
    }
  }
}

.edit-form {
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
    cursor: pointer;
    transition: background 0.3s;

    &:hover {
      background: #0056b3;
    }
  }
}

.modal-form {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  padding: 1rem;

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
    cursor: pointer;
    transition: background 0.3s;

    &:hover {
      background: #0056b3;
    }
  }
}

.image-upload {
  margin: 1rem 0;

  input[type="file"] {
    width: 100%;
    padding: 0.5rem;
    border: 1px dashed #ddd;
    border-radius: 4px;
  }

  .preview-image {
    max-width: 200px;
    margin-top: 1rem;
    border-radius: 4px;
  }
}

.coordinates {
  display: none;
}
</style> 