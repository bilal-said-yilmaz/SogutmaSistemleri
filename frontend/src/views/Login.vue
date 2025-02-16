<template>
  <div class="login-page">
    <div class="login-form">
      <h2>Admin Girişi</h2>
      <form @submit.prevent="login">
        <input v-model="form.username" type="text" placeholder="Kullanıcı Adı" required>
        <input v-model="form.password" type="password" placeholder="Şifre" required>
        <button type="submit">Giriş Yap</button>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useStore } from 'vuex'
import axios from 'axios'

const router = useRouter()
const store = useStore()

const form = ref({
  username: '',
  password: ''
})

const login = async () => {
  try {
    const response = await axios.post('/api/auth/login', form.value)
    const token = response.data.token
    localStorage.setItem('token', token)
    store.commit('setAuth', true)
    store.commit('setUser', { username: form.value.username })
    
    // Axios için default Authorization header'ını ayarla
    axios.defaults.headers.common['Authorization'] = `Bearer ${token}`
    
    router.push('/admin')
  } catch (error) {
    console.error('Login hatası:', error)
    alert('Giriş başarısız!')
  }
}
</script>

<style lang="scss" scoped>
.login-page {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: calc(100vh - 120px);
}

.login-form {
  width: 100%;
  max-width: 400px;
  padding: 2rem;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);

  h2 {
    text-align: center;
    margin-bottom: 2rem;
    color: #333;
  }

  form {
    display: flex;
    flex-direction: column;
    gap: 1rem;

    input {
      padding: 0.8rem;
      border: 1px solid #ddd;
      border-radius: 4px;
      font-size: 1rem;

      &:focus {
        outline: none;
        border-color: #007bff;
      }
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
</style> 