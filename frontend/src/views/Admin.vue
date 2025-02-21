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
      <div class="admin-header">
        <h1>Admin Panel</h1>
        <div class="admin-actions">
          <router-link to="/" class="btn btn-secondary">Ana Sayfaya Dön</router-link>
          <button @click="handleLogout" class="btn btn-danger">Çıkış Yap</button>
        </div>
      </div>

      <div class="admin-content">
        <div class="admin-sidebar">
          <nav class="admin-nav">
            <button 
              v-for="tab in tabs" 
              :key="tab.id"
              :class="['nav-button', { active: currentTab === tab.id }]"
              @click="currentTab = tab.id"
            >
              <i :class="getTabIcon(tab.id)"></i>
              {{ tab.name }}
            </button>
          </nav>
        </div>

        <div class="admin-main">
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
              
              <div class="working-hours">
                <h4>Çalışma Saatleri</h4>
                <div class="hours-input">
                  <label>Hafta İçi:</label>
                  <input v-model="contact.weekdayHours" type="text" placeholder="örn: 09:00 - 18:00">
                </div>
                <div class="hours-input">
                  <label>Cumartesi:</label>
                  <input v-model="contact.saturdayHours" type="text" placeholder="örn: 09:00 - 14:00">
                </div>
                <div class="hours-input">
                  <label>Pazar:</label>
                  <input v-model="contact.sundayHours" type="text" placeholder="örn: Kapalı">
                </div>
              </div>
              <button type="submit">Güncelle</button>
            </form>
          </div>

          <!-- Ana Sayfa -->
          <div v-if="currentTab === 'hero'" class="tab-content">
            <h2>Ana Sayfa</h2>
            <form @submit.prevent="updateHero" class="edit-form">
              <input v-model="hero.subheading" type="text" placeholder="Alt Başlık">
              <input v-model="hero.heading" type="text" placeholder="Başlık">
              <input v-model="hero.buttonText" type="text" placeholder="Buton Metni">
              <div class="image-upload">
                <input type="file" @change="handleHeroImage" accept="image/*" :required="!hero.backgroundImage">
                <img v-if="hero.backgroundImage" :src="hero.backgroundImage" class="preview-image">
              </div>
              <button type="submit">Güncelle</button>
            </form>
          </div>

          <!-- Footer -->
          <div v-if="currentTab === 'footer'" class="tab-content">
            <h2>Footer</h2>
            <form @submit.prevent="updateFooter" class="edit-form">
              <input v-model="footer.copyright" type="text" placeholder="Telif Hakkı">
              <input v-model="footer.socialLinks.twitter" type="text" placeholder="Twitter">
              <input v-model="footer.socialLinks.facebook" type="text" placeholder="Facebook">
              <input v-model="footer.socialLinks.instagram" type="text" placeholder="Instagram">
              <input v-model="footer.links.privacy" type="text" placeholder="Gizlilik Politikası">
              <input v-model="footer.links.terms" type="text" placeholder="Kullanım Şartları">
              <button type="submit">Güncelle</button>
            </form>
          </div>

          <!-- Users Section -->
          <div v-if="currentTab === 'users'" class="tab-content">
            <h2>Kullanıcı Yönetimi</h2>
            <button @click="openAddUserModal" class="add-button">
              Yeni Kullanıcı Ekle
            </button>
            
            <div class="users-table">
              <table>
                <thead>
                  <tr>
                    <th>ID</th>
                    <th>Kullanıcı Adı</th>
                    <th>Telefon</th>
                    <th>Rol</th>
                    <th>İşlemler</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="user in users" :key="user.id">
                    <td>{{ user.id }}</td>
                    <td>{{ user.username }}</td>
                    <td>{{ user.phone }}</td>
                    <td>{{ getRoleName(user.role_id) }}</td>
                    <td class="actions">
                      <button @click="editUser(user)" class="edit-btn">
                        <i class="fas fa-edit"></i>
                      </button>
                      <button @click="deleteUser(user.id)" class="delete-btn">
                        <i class="fas fa-trash"></i>
                      </button>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>
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

    <!-- User Add/Edit Modal -->
    <Dialog v-model:visible="showAddUserModal" :header="editingUser ? 'Kullanıcı Düzenle' : 'Yeni Kullanıcı Ekle'" class="user-modal">
      <form @submit.prevent="saveUser" class="modal-form">
        <div class="form-group">
          <label>Kullanıcı Adı</label>
          <input v-model="userForm.username" type="text" placeholder="Kullanıcı Adı" required>
        </div>

        <div class="form-group">
          <label>Telefon Numarası</label>
          <input v-model="userForm.phone" type="tel" placeholder="5XX XXX XX XX" 
                 required pattern="[0-9]{10}" 
                 @input="formatPhoneNumber">
        </div>

        <div class="form-group">
          <label>Kullanıcı Rolü</label>
          <select v-model="userForm.role_id" required class="role-select">
            <option value="">Rol Seçin</option>
            <option value="1">Admin</option>
            <option value="2">Editör</option>
          </select>
        </div>

        <template v-if="!editingUser">
          <div class="form-group">
            <label>Şifre (en az 6 karakter)</label>
            <input v-model="userForm.password" type="password" placeholder="Şifre" required minlength="6">
          </div>

          <div class="form-group">
            <label>Şifre (Tekrar)</label>
            <input v-model="userForm.confirmPassword" type="password" placeholder="Şifreyi tekrar girin" required minlength="6">
          </div>

          <div class="password-requirements">
            <p>Şifre gereksinimleri:</p>
            <ul>
              <li>En az 6 karakter uzunluğunda olmalı</li>
              <li>Şifre ve şifre tekrarı aynı olmalı</li>
            </ul>
          </div>
        </template>

        <div v-if="formError" class="error-message">{{ formError }}</div>

        <div class="modal-actions">
          <button type="button" @click="showAddUserModal = false" class="btn btn-secondary">İptal</button>
          <button type="submit" class="btn btn-primary">{{ editingUser ? 'Güncelle' : 'Ekle' }}</button>
        </div>
      </form>
    </Dialog>
  </div>
</template>

<script setup>
import { ref, computed, watch, nextTick } from 'vue'
import { useStore } from 'vuex'
import axios from 'axios'
import Dialog from 'primevue/dialog'
import { useRouter } from 'vue-router'

const store = useStore()
const isAuthenticated = computed(() => store.state.isAuthenticated)
const router = useRouter()

const currentTab = ref('products')
const tabs = [
  { id: 'products', name: 'Ürünler' },
  { id: 'services', name: 'Hizmetler' },
  { id: 'about', name: 'Hakkımızda' },
  { id: 'contact', name: 'İletişim' },
  { id: 'users', name: 'Kullanıcılar' }
]

const loginForm = ref({
  username: '',
  password: ''
})

const products = ref([])
const services = ref([])
const about = ref({})
const contact = ref({})
const hero = ref({
  subheading: 'Klima Kozan\'a Hoş Geldiniz',
  heading: 'Profesyonel Klima ve Beyaz Eşya Çözümleri',
  buttonText: 'Hizmetlerimizi Keşfedin',
  backgroundImage: ''
})
const footer = ref({
  copyright: 'Copyright © Klima Kozan 2024',
  socialLinks: {
    twitter: '#',
    facebook: '#',
    instagram: '#'
  },
  links: {
    privacy: 'Gizlilik Politikası',
    terms: 'Kullanım Şartları'
  }
})
const users = ref([])

const showAddProductModal = ref(false)
const showAddServiceModal = ref(false)
const editingProduct = ref(null)
const editingService = ref(null)
const showAddUserModal = ref(false)
const editingUser = ref(null)

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

const userForm = ref({
  username: '',
  phone: '',
  role_id: '',
  password: '',
  confirmPassword: ''
})

const formError = ref('')

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
    const [productsRes, servicesRes, aboutRes, contactRes, heroRes, footerRes, usersRes] = await Promise.all([
      axios.get('/api/admin/products'),
      axios.get('/api/admin/services'),
      axios.get('/api/admin/about'),
      axios.get('/api/admin/contact'),
      axios.get('/api/admin/hero'),
      axios.get('/api/admin/footer'),
      axios.get('/api/admin/users')
    ])

    products.value = productsRes.data
    services.value = servicesRes.data
    about.value = aboutRes.data
    contact.value = contactRes.data
    hero.value = heroRes.data
    footer.value = footerRes.data
    users.value = usersRes.data
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

// Hero güncelleme
const updateHero = async () => {
  try {
    await axios.put('/api/admin/hero', hero.value)
    alert('Ana sayfa içeriği güncellendi!')
  } catch (error) {
    alert('Güncelleme başarısız!')
  }
}

// Footer güncelleme
const updateFooter = async () => {
  try {
    await axios.put('/api/admin/footer', footer.value)
    alert('Footer içeriği güncellendi!')
  } catch (error) {
    alert('Güncelleme başarısız!')
  }
}

// Resim yükleme işlemleri
const handleAboutImage = (event) => {
  const file = event.target.files[0]
  if (file) {
    const reader = new FileReader()
    reader.onload = (e) => {
      about.value.image = e.target.result
    }
    reader.readAsDataURL(file)
  }
}

// İletişim işlemleri
const handleProductImage = (event) => {
  const file = event.target.files[0]
  if (file) {
    const reader = new FileReader()
    reader.onload = (e) => {
      productForm.value.image = e.target.result
    }
    reader.readAsDataURL(file)
  }
}

const handleServiceImage = (event) => {
  const file = event.target.files[0]
  if (file) {
    const reader = new FileReader()
    reader.onload = (e) => {
      serviceForm.value.image = e.target.result
    }
    reader.readAsDataURL(file)
  }
}

// Hero resmi yükleme
const handleHeroImage = (event) => {
  const file = event.target.files[0]
  if (file) {
    const reader = new FileReader()
    reader.onload = (e) => {
      hero.value.backgroundImage = e.target.result
    }
    reader.readAsDataURL(file)
  }
}

// Add Google Maps URL generator
const getGoogleMapsUrl = () => {
  if (!contact.value?.latitude || !contact.value?.longitude) return '#'
  return `https://www.google.com/maps?q=${contact.value.latitude},${contact.value.longitude}`
}

// Tab ikonlarını belirle
const getTabIcon = (tabId) => {
  const icons = {
    products: 'fas fa-box',
    services: 'fas fa-tools',
    about: 'fas fa-info-circle',
    contact: 'fas fa-address-card',
    hero: 'fas fa-home',
    footer: 'fas fa-shoe-prints',
    users: 'fas fa-users'
  }
  return icons[tabId]
}

// Çıkış işlemi
const handleLogout = () => {
  store.dispatch('logout')
  router.push('/login')
}

// Kullanıcı işlemleri
const formatPhoneNumber = (event) => {
  // Remove non-numeric characters
  let phone = event.target.value.replace(/\D/g, '')
  
  // Limit to 10 digits
  if (phone.length > 10) {
    phone = phone.slice(0, 10)
  }
  
  userForm.value.phone = phone
}

const saveUser = async () => {
  try {
    formError.value = ''

    // Validate phone number
    if (userForm.value.phone.length !== 10) {
      formError.value = 'Geçerli bir telefon numarası giriniz (5XX XXX XX XX)'
      return
    }

    // Validate password match for new users
    if (!editingUser.value && userForm.value.password !== userForm.value.confirmPassword) {
      formError.value = 'Şifreler eşleşmiyor'
      return
    }

    const userData = {
      username: userForm.value.username,
      phone: userForm.value.phone,
      role_id: parseInt(userForm.value.role_id),
      password: userForm.value.password
    }

    console.log('Gönderilen kullanıcı verisi:', userData)

    if (editingUser.value) {
      await axios.put(`/api/admin/users/${editingUser.value.id}`, userData)
    } else {
      await axios.post('/api/admin/users', userData)
    }

    await fetchData()
    showAddUserModal.value = false
    userForm.value = { username: '', phone: '', role_id: '', password: '', confirmPassword: '' }
    editingUser.value = null
    formError.value = ''
  } catch (error) {
    console.error('Hata detayı:', error.response?.data)
    formError.value = error.response?.data?.error || 'İşlem başarısız oldu. Lütfen tekrar deneyin.'
  }
}

const editUser = (user) => {
  editingUser.value = user
  userForm.value = {
    username: user.username,
    phone: user.phone || '',
    role_id: user.role_id,
    password: '',
    confirmPassword: ''
  }
  showAddUserModal.value = true
}

const deleteUser = async (id) => {
  if (confirm('Bu kullanıcıyı silmek istediğinizden emin misiniz?')) {
    try {
      await axios.delete(`/api/admin/users/${id}`)
      await fetchData()
    } catch (error) {
      alert('Silme işlemi başarısız!')
    }
  }
}

// Rol adını getiren fonksiyon
const getRoleName = (roleId) => {
  switch (roleId) {
    case 1:
      return 'Admin'
    case 2:
      return 'Editör'
    default:
      return 'Bilinmeyen'
  }
}

// Add after the getRoleName function:
const openAddUserModal = () => {
  editingUser.value = null
  userForm.value = {
    username: '',
    phone: '',
    role_id: '',
    password: '',
    confirmPassword: ''
  }
  formError.value = ''
  showAddUserModal.value = true
}

// Sayfa yüklendiğinde verileri getir
if (isAuthenticated.value) {
  fetchData()
}
</script>

<style scoped>
.admin {
  min-height: 100vh;
  background-color: #f8f9fa;
  padding-top: 76px;
}

.admin-header, .admin-actions {
  display: none;
}

.admin-content {
  display: flex;
  gap: 2rem;
  padding: 2rem;
  max-width: 1400px;
  margin: 0 auto;
}

.admin-sidebar {
  position: sticky;
  top: 92px;
  height: calc(100vh - 92px);
  overflow-y: auto;
  width: 250px;
  flex-shrink: 0;
}

.admin-main {
  flex-grow: 1;
}

.admin-nav {
  background: #fff;
  border-radius: 8px;
  padding: 1rem;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.nav-button {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.75rem 1rem;
  border: none;
  background: transparent;
  color: #666;
  font-size: 1rem;
  text-align: left;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.nav-button i {
  width: 20px;
  text-align: center;
}

.nav-button:hover {
  background: #f8f9fa;
  color: #333;
}

.nav-button.active {
  background: #007bff;
  color: #fff;
}

.btn {
  padding: 0.5rem 1rem;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.9rem;
  text-decoration: none;
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
}

.btn-secondary {
  background: #6c757d;
  color: #fff;
}

.btn-danger {
  background: #dc3545;
  color: #fff;
}

.tab-content {
  background: #fff;
  padding: 2rem;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.add-button {
  margin-bottom: 2rem;
  padding: 0.5rem 1rem;
  background: #28a745;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.items-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 2rem;
}

.item-card {
  background: white;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.item-card img {
  width: 100%;
  height: 200px;
  object-fit: cover;
}

.item-info {
  padding: 1rem;
}

.item-info h3 {
  margin: 0 0 0.5rem;
}

.item-info p {
  color: #666;
  margin: 0;
}

.item-actions {
  padding: 1rem;
  display: flex;
  gap: 1rem;
}

.item-actions button {
  flex: 1;
  padding: 0.5rem;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.item-actions button.delete {
  background: #dc3545;
  color: white;
}

.edit-form {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  max-width: 600px;
  margin: 0 auto;
}

.edit-form input,
.edit-form textarea {
  padding: 0.8rem;
  border: 1px solid #ddd;
  border-radius: 4px;
}

.edit-form textarea {
  min-height: 150px;
}

.image-upload {
  margin: 1rem 0;
}

.preview-image {
  max-width: 200px;
  margin-top: 1rem;
}

.working-hours {
  margin-top: 1.5rem;
}

.hours-input {
  display: grid;
  grid-template-columns: 120px 1fr;
  gap: 1rem;
  align-items: center;
  margin-bottom: 0.5rem;
}

.hours-input input {
  width: 100%;
  padding: 0.8rem;
  border: 1px solid #ddd;
  border-radius: 4px;
}

.modal-form {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  padding: 1rem;
}

.modal-form input,
.modal-form textarea {
  padding: 0.8rem;
  border: 1px solid #ddd;
  border-radius: 4px;
}

.modal-form textarea {
  min-height: 150px;
}

.modal-form button {
  padding: 0.8rem;
  background: #007bff;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

@media (max-width: 768px) {
  .admin {
    padding-top: 60px;
  }

  .admin-content {
    flex-direction: column;
    padding: 1rem;
  }

  .admin-sidebar {
    position: static;
    height: auto;
    width: 100%;
    margin-bottom: 1rem;
  }

  .admin-nav {
    flex-direction: row;
    flex-wrap: wrap;
    justify-content: center;
    gap: 0.5rem;
  }

  .nav-button {
    flex: 1;
    min-width: 140px;
    max-width: calc(50% - 0.5rem);
    justify-content: center;
  }

  .items-grid {
    grid-template-columns: 1fr;
  }

  .item-card {
    margin-bottom: 1rem;
  }

  .edit-form {
    padding: 1rem;
  }

  .hours-input {
    grid-template-columns: 1fr;
  }

  .hours-input label {
    margin-bottom: 0.5rem;
  }

  .coordinate-inputs {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 480px) {
  .admin-content {
    padding: 0.5rem;
  }

  .tab-content {
    padding: 1rem;
  }

  .nav-button {
    min-width: 100%;
    max-width: 100%;
  }

  .item-actions {
    flex-direction: column;
  }

  .item-actions button {
    width: 100%;
  }
}

.users-table {
  width: 100%;
  overflow-x: auto;
  margin-top: 2rem;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.users-table table {
  width: 100%;
  border-collapse: collapse;
}

.users-table th,
.users-table td {
  padding: 1rem;
  text-align: left;
  border-bottom: 1px solid #eee;
}

.users-table th {
  background: #f8f9fa;
  font-weight: 600;
  color: #333;
}

.users-table .actions {
  display: flex;
  gap: 0.5rem;
}

.users-table .edit-btn,
.users-table .delete-btn {
  padding: 0.5rem;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.users-table .edit-btn {
  background: #007bff;
  color: white;
}

.users-table .delete-btn {
  background: #dc3545;
  color: white;
}

.users-table .edit-btn:hover {
  background: #0056b3;
}

.users-table .delete-btn:hover {
  background: #c82333;
}

.error-message {
  color: #dc3545;
  font-size: 0.9rem;
  margin: 0.5rem 0;
  text-align: center;
}

/* User Modal Styles */
.user-modal {
  max-width: 500px;
}

.user-modal .form-group {
  margin-bottom: 1rem;
}

.user-modal .form-group label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 500;
  color: #333;
}

.user-modal .form-group input,
.user-modal .form-group select {
  width: 100%;
  padding: 0.8rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 1rem;
  transition: border-color 0.3s;
}

.user-modal .form-group input:focus,
.user-modal .form-group select:focus {
  border-color: #007bff;
  outline: none;
  box-shadow: 0 0 0 2px rgba(0,123,255,0.25);
}

.user-modal .modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
  margin-top: 2rem;
}

.user-modal .btn {
  padding: 0.8rem 1.5rem;
  border: none;
  border-radius: 4px;
  font-size: 1rem;
  cursor: pointer;
  transition: all 0.3s;
}

.user-modal .btn-primary {
  background: #007bff;
  color: white;
}

.user-modal .btn-primary:hover {
  background: #0056b3;
}

.user-modal .btn-secondary {
  background: #6c757d;
  color: white;
}

.user-modal .btn-secondary:hover {
  background: #5a6268;
}

.user-modal .error-message {
  color: #dc3545;
  font-size: 0.9rem;
  margin-top: 1rem;
  text-align: center;
  padding: 0.5rem;
  background: rgba(220,53,69,0.1);
  border-radius: 4px;
}

.role-select {
  background-color: #f8f9fa;
  font-size: 1rem;
  padding: 0.8rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  width: 100%;
  cursor: pointer;
  transition: all 0.3s;
}

.role-select:hover, .role-select:focus {
  border-color: #007bff;
  background-color: #fff;
}

.password-requirements {
  margin: 1rem 0;
  padding: 1rem;
  background-color: #f8f9fa;
  border-radius: 4px;
  border-left: 4px solid #007bff;
}

.password-requirements p {
  margin: 0 0 0.5rem 0;
  font-weight: 500;
  color: #333;
}

.password-requirements ul {
  margin: 0;
  padding-left: 1.5rem;
  color: #666;
}

.password-requirements li {
  margin: 0.25rem 0;
}
</style> 