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

          <!-- Users Tab -->
          <div v-if="currentTab === 'users'" class="tab-content">
            <button @click="openAddUserModal" class="add-button">
              <i class="fas fa-plus"></i> Yeni Kullanıcı Ekle
            </button>
            
            <div class="users-table">
              <table>
                <thead>
                  <tr>
                    <th>Kullanıcı Adı</th>
                    <th>E-posta</th>
                    <th>Rol</th>
                    <th>İşlemler</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="user in users" :key="user.id">
                    <td>{{ user.username }}</td>
                    <td>{{ user.email }}</td>
                    <td>{{ getRoleName(user.role_id) }}</td>
                    <td class="actions">
                      <button @click="editUser(user)" class="btn-edit">
                        <i class="fas fa-edit"></i>
                      </button>
                      <button @click="deleteUser(user.id)" class="btn-delete">
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

    <!-- Add User Modal -->
    <Dialog v-model:visible="showAddUserModal" header="Kullanıcı Ekle" :style="{ width: '400px' }">
      <div class="user-form">
        <div class="form-group">
          <label>Kullanıcı Adı</label>
          <input type="text" v-model="userForm.username" class="form-control" required>
        </div>
        <div class="form-group">
          <label>E-posta</label>
          <input type="email" v-model="userForm.email" class="form-control" required>
        </div>
        <div class="form-group">
          <label>Şifre</label>
          <input type="password" v-model="userForm.password" class="form-control" required>
        </div>
        <div class="form-group">
          <label>Rol</label>
          <select v-model="userForm.role_id" class="form-control">
            <option value="1">Admin</option>
            <option value="2">Editör</option>
          </select>
        </div>
        <div class="modal-actions">
          <button type="button" @click="closeAddUserModal" class="btn btn-secondary">İptal</button>
          <button type="button" @click="saveUser" class="btn btn-primary">
            <i class="fas fa-save"></i> Kaydet
          </button>
        </div>
      </div>
    </Dialog>

    <!-- Edit User Modal -->
    <Dialog v-model:visible="showEditUserModal" header="Kullanıcı Düzenle" :style="{ width: '400px' }">
      <div class="user-form">
        <div class="form-group">
          <label>Kullanıcı Adı</label>
          <input type="text" v-model="userForm.username" class="form-control" required>
        </div>
        <div class="form-group">
          <label>E-posta</label>
          <input type="email" v-model="userForm.email" class="form-control" required>
        </div>
        <div class="form-group">
          <label>Rol</label>
          <select v-model="userForm.role_id" class="form-control">
            <option value="1">Admin</option>
            <option value="2">Editör</option>
          </select>
        </div>
        <div class="modal-actions">
          <button type="button" @click="closeEditUserModal" class="btn btn-secondary">İptal</button>
          <button type="button" @click="updateUser" class="btn btn-primary">
            <i class="fas fa-save"></i> Güncelle
          </button>
        </div>
      </div>
    </Dialog>

    <!-- Copyright Footer -->
    <div class="copyright-footer">
      <p>&copy; {{ new Date().getFullYear() }} Kozan İklimlendirme. Tüm hakları saklıdır.</p>
    </div>
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
const users = ref([])

const showAddProductModal = ref(false)
const showAddServiceModal = ref(false)
const editingProduct = ref(null)
const editingService = ref(null)
const showAddUserModal = ref(false)
const showEditUserModal = ref(false)
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
  email: '',
  password: '',
  role_id: 2
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
    const [productsRes, servicesRes, aboutRes, contactRes, usersRes] = await Promise.all([
      axios.get('/api/admin/products'),
      axios.get('/api/admin/services'),
      axios.get('/api/admin/about'),
      axios.get('/api/admin/contact'),
      axios.get('/api/admin/users')
    ])

    products.value = productsRes.data
    services.value = servicesRes.data
    about.value = aboutRes.data
    contact.value = contactRes.data
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

// Tab ikonlarını belirle
const getTabIcon = (tabId) => {
  const icons = {
    products: 'fas fa-box',
    services: 'fas fa-tools',
    about: 'fas fa-info-circle',
    contact: 'fas fa-address-card',
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

const resetUserForm = () => {
  userForm.value = {
    username: '',
    email: '',
    password: '',
    role_id: 2
  }
}

const saveUser = async () => {
  try {
    const userData = {
      username: userForm.value.username,
      email: userForm.value.email,
      password: userForm.value.password,
      role_id: parseInt(userForm.value.role_id)
    }
    await axios.post('/api/admin/users', userData)
    await fetchData()
    closeAddUserModal()
    alert('Kullanıcı başarıyla eklendi')
  } catch (error) {
    if (error.response?.data?.error) {
      alert(error.response.data.error)
    } else {
      alert('Kullanıcı eklenirken bir hata oluştu')
    }
  }
}

const editUser = (user) => {
  userForm.value = {
    id: user.id,
    username: user.username,
    email: user.email,
    role_id: user.role_id
  }
  showEditUserModal.value = true
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
  resetUserForm()
  showAddUserModal.value = true
}

const closeAddUserModal = () => {
  showAddUserModal.value = false
}

const closeEditUserModal = () => {
  showEditUserModal.value = false
}

const updateUser = async () => {
  try {
    const userData = {
      username: userForm.value.username,
      email: userForm.value.email,
      role_id: parseInt(userForm.value.role_id)
    }
    await axios.put(`/api/admin/users/${userForm.value.id}`, userData)
    await fetchData()
    closeEditUserModal()
    alert('Kullanıcı başarıyla güncellendi')
  } catch (error) {
    if (error.response?.data?.error) {
      alert(error.response.data.error)
    } else {
      alert('Kullanıcı güncellenirken bir hata oluştu')
    }
  }
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
  display: flex;
  flex-direction: column;
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
}

.users-table table {
  width: 100%;
  border-collapse: collapse;
  margin-top: 1rem;
}

.users-table th,
.users-table td {
  padding: 1rem;
  text-align: left;
  border-bottom: 1px solid #eee;
}

.users-table th {
  background-color: #f8f9fa;
  font-weight: 600;
}

.users-table tr:hover {
  background-color: #f8f9fa;
}

.actions {
  display: flex;
  gap: 0.5rem;
}

.btn-edit,
.btn-delete {
  padding: 0.5rem;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s;
}

.btn-edit {
  background-color: #007bff;
  color: white;
}

.btn-delete {
  background-color: #dc3545;
  color: white;
}

.btn-edit:hover {
  background-color: #0056b3;
}

.btn-delete:hover {
  background-color: #c82333;
}

.user-form {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  padding: 1rem 0;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.form-group label {
  font-weight: 500;
  color: #333;
}

.form-control {
  padding: 0.5rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 1rem;
}

.form-control:focus {
  outline: none;
  border-color: #007bff;
}

select.form-control {
  cursor: pointer;
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

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
  margin-top: 1.5rem;
  padding-top: 1rem;
  border-top: 1px solid #eee;
}

.btn {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1.5rem;
  border: none;
  border-radius: 4px;
  font-size: 1rem;
  cursor: pointer;
  transition: all 0.3s ease;
}

.btn i {
  font-size: 1rem;
}

.btn-primary {
  background-color: #007bff;
  color: white;
}

.btn-primary:hover {
  background-color: #0056b3;
}

.btn-secondary {
  background-color: #6c757d;
  color: white;
}

.btn-secondary:hover {
  background-color: #5a6268;
}

/* Copyright Footer */
.copyright-footer {
  text-align: center;
  padding: 1rem;
  background-color: #f8f9fa;
  border-top: 1px solid #eee;
  color: #666;
  font-size: 0.9rem;
  margin-top: auto;
  width: 100%;
}

@media (max-width: 768px) {
  .admin {
    padding-top: 60px;
  }
  
  .copyright-footer {
    padding: 0.75rem;
    font-size: 0.85rem;
  }
}
</style> 