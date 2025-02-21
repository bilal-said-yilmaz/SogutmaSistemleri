<template>
  <div class="login-page">
    <!-- Normal Login Form -->
    <div v-if="!showForgotPassword" class="login-form">
      <h2>Admin Girişi</h2>
      <form @submit.prevent="login">
        <input v-model="form.username" type="text" placeholder="Kullanıcı Adı" required>
        <input v-model="form.password" type="password" placeholder="Şifre" required>
        <div class="forgot-password">
          <a href="#" @click.prevent="toggleForgotPassword">Şifremi Unuttum</a>
        </div>
        <button type="submit">Giriş Yap</button>
      </form>
    </div>

    <!-- Forgot Password Form -->
    <div v-else class="login-form">
      <h2>Şifre Yenileme</h2>
      <!-- Phone Number Input Step -->
      <form v-if="forgotPasswordStep === 1" @submit.prevent="sendVerificationCode">
        <input v-model="forgotPasswordForm.phone" type="tel" placeholder="Telefon Numarası (5XX XXX XX XX)" required pattern="[0-9]{10}">
        <button type="submit">Doğrulama Kodu Gönder</button>
        <button type="button" @click="toggleForgotPassword" class="mt-3" style="background: #6c757d;">Geri Dön</button>
      </form>

      <!-- Verification Code Input Step -->
      <form v-else-if="forgotPasswordStep === 2" @submit.prevent="verifyCode">
        <div class="verification-code">
          <input v-for="(digit, index) in verificationCode" 
                 :key="index"
                 v-model="verificationCode[index]"
                 maxlength="1"
                 @input="onCodeInput($event, index)"
                 @keydown.delete="onCodeDelete($event, index)"
                 ref="codeInputs"
                 type="text"
                 pattern="[0-9]"
                 required>
        </div>
        <div class="resend-code">
          <button type="button" @click="resendCode" :disabled="resendTimer > 0">
            Kodu Tekrar Gönder
            <span v-if="resendTimer > 0" class="timer">({{ resendTimer }}s)</span>
          </button>
        </div>
        <button type="submit">Doğrula</button>
        <button type="button" @click="forgotPasswordStep = 1" class="mt-3" style="background: #6c757d;">Geri Dön</button>
      </form>

      <!-- New Password Input Step -->
      <form v-else-if="forgotPasswordStep === 3" @submit.prevent="resetPassword">
        <input v-model="forgotPasswordForm.newPassword" type="password" placeholder="Yeni Şifre" required minlength="6">
        <input v-model="forgotPasswordForm.confirmPassword" type="password" placeholder="Yeni Şifre (Tekrar)" required minlength="6">
        <button type="submit">Şifreyi Güncelle</button>
        <button type="button" @click="forgotPasswordStep = 2" class="mt-3" style="background: #6c757d;">Geri Dön</button>
      </form>

      <div v-if="errorMessage" class="error-message">{{ errorMessage }}</div>
      <div v-if="successMessage" class="success-message">{{ successMessage }}</div>
    </div>
  </div>
</template>

<script setup>
import { ref, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { useStore } from 'vuex'
import axios from 'axios'

const router = useRouter()
const store = useStore()

const form = ref({
  username: '',
  password: ''
})

const showForgotPassword = ref(false)
const forgotPasswordStep = ref(1)
const forgotPasswordForm = ref({
  phone: '',
  newPassword: '',
  confirmPassword: ''
})
const verificationCode = ref(['', '', '', '', '', ''])
const resendTimer = ref(0)
const errorMessage = ref('')
const successMessage = ref('')
const codeInputs = ref([])

const login = async () => {
  try {
    const response = await axios.post('/api/auth/login', form.value)
    const token = response.data.token
    localStorage.setItem('token', token)
    store.commit('setAuth', true)
    store.commit('setUser', { username: form.value.username })
    
    axios.defaults.headers.common['Authorization'] = `Bearer ${token}`
    
    const redirectPath = router.currentRoute.value.query.redirect || '/admin'
    router.push(redirectPath)
  } catch (error) {
    errorMessage.value = 'Giriş başarısız! Lütfen kullanıcı adı ve şifrenizi kontrol edin.'
  }
}

const toggleForgotPassword = () => {
  showForgotPassword.value = !showForgotPassword.value
  forgotPasswordStep.value = 1
  errorMessage.value = ''
  successMessage.value = ''
  forgotPasswordForm.value = {
    phone: '',
    newPassword: '',
    confirmPassword: ''
  }
  verificationCode.value = ['', '', '', '', '', '']
}

const sendVerificationCode = async () => {
  try {
    // Format phone number
    const phoneNumber = forgotPasswordForm.value.phone.replace(/\D/g, '')
    if (phoneNumber.length !== 10) {
      throw new Error('Geçerli bir telefon numarası giriniz (5XX XXX XX XX)')
    }

    await axios.post('/api/auth/forgot-password', {
      phone: phoneNumber
    })
    
    forgotPasswordStep.value = 2
    startResendTimer()
    successMessage.value = 'Doğrulama kodu telefonunuza gönderildi'
    errorMessage.value = ''
  } catch (error) {
    errorMessage.value = error.response?.data?.message || error.message
    successMessage.value = ''
  }
}

const startResendTimer = () => {
  resendTimer.value = 120 // 2 minutes
  const timer = setInterval(() => {
    resendTimer.value--
    if (resendTimer.value === 0) {
      clearInterval(timer)
    }
  }, 1000)
}

const resendCode = async () => {
  if (resendTimer.value === 0) {
    await sendVerificationCode()
  }
}

const verifyCode = async () => {
  try {
    const code = verificationCode.value.join('')
    if (code.length !== 6) {
      throw new Error('Lütfen 6 haneli doğrulama kodunu giriniz')
    }

    await axios.post('/api/auth/verify-code', {
      phone: forgotPasswordForm.value.phone,
      code: code
    })
    
    forgotPasswordStep.value = 3
    successMessage.value = 'Doğrulama başarılı'
    errorMessage.value = ''
  } catch (error) {
    errorMessage.value = error.response?.data?.message || 'Doğrulama kodu hatalı'
    successMessage.value = ''
  }
}

const resetPassword = async () => {
  try {
    if (forgotPasswordForm.value.newPassword !== forgotPasswordForm.value.confirmPassword) {
      throw new Error('Şifreler eşleşmiyor')
    }

    await axios.post('/api/auth/reset-password', {
      phone: forgotPasswordForm.value.phone,
      password: forgotPasswordForm.value.newPassword
    })
    
    successMessage.value = 'Şifreniz başarıyla güncellendi'
    errorMessage.value = ''
    
    // 3 saniye sonra login formuna dön
    setTimeout(() => {
      toggleForgotPassword()
    }, 3000)
  } catch (error) {
    errorMessage.value = error.response?.data?.message || error.message
    successMessage.value = ''
  }
}

const onCodeInput = (event, index) => {
  const value = event.target.value
  if (value && index < 5) {
    verificationCode.value[index] = value
    verificationCode.value = [...verificationCode.value]
    nextTick(() => {
      if (codeInputs.value[index + 1]) {
        codeInputs.value[index + 1].focus()
      }
    })
  }
}

const onCodeDelete = (event, index) => {
  if (event.key === 'Backspace' && !verificationCode.value[index] && index > 0) {
    verificationCode.value[index - 1] = ''
    verificationCode.value = [...verificationCode.value]
    nextTick(() => {
      if (codeInputs.value[index - 1]) {
        codeInputs.value[index - 1].focus()
      }
    })
  }
}
</script>

<style scoped>
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
}

.login-form h2 {
  text-align: center;
  margin-bottom: 2rem;
  color: #333;
}

.login-form form {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.login-form form input {
  padding: 0.8rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 1rem;
}

.login-form form input:focus {
  outline: none;
  border-color: #007bff;
}

.login-form form button {
  padding: 0.8rem;
  background: #007bff;
  color: #fff;
  border: none;
  border-radius: 4px;
  font-size: 1rem;
  cursor: pointer;
  transition: background 0.3s;
}

.login-form form button:hover {
  background: #0056b3;
}

.forgot-password {
  text-align: right;
  margin-top: -0.5rem;
}

.forgot-password a {
  color: #007bff;
  text-decoration: none;
  font-size: 0.9rem;
}

.forgot-password a:hover {
  text-decoration: underline;
}

.verification-code {
  display: flex;
  gap: 0.5rem;
  justify-content: center;
  margin: 1rem 0;
}

.verification-code input {
  width: 40px;
  height: 40px;
  text-align: center;
  font-size: 1.2rem;
  border: 1px solid #ddd;
  border-radius: 4px;
}

.verification-code input:focus {
  outline: none;
  border-color: #007bff;
}

.resend-code {
  text-align: center;
  margin-top: 1rem;
}

.resend-code button {
  background: none;
  border: none;
  color: #007bff;
  cursor: pointer;
  font-size: 0.9rem;
}

.resend-code button:disabled {
  color: #999;
  cursor: not-allowed;
}

.timer {
  color: #666;
  font-size: 0.9rem;
  margin-left: 0.5rem;
}

.error-message {
  color: #dc3545;
  font-size: 0.9rem;
  margin-top: 0.5rem;
  text-align: center;
}

.success-message {
  color: #28a745;
  font-size: 0.9rem;
  margin-top: 0.5rem;
  text-align: center;
}

.mt-3 {
  margin-top: 1rem;
}
</style> 