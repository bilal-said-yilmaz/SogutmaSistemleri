<template>
  <div class="app">
    <nav class="navbar">
      <div class="nav-brand">
        <router-link to="/" class="home-link" :title="isAdminRoute ? 'Anasayfaya Dön' : ''">
          <i class="pi pi-home" v-if="isAdminRoute"></i>
          <span>Kozan</span>
        </router-link>
      </div>
      <div class="nav-links" v-if="!isAdminRoute">
        <a href="#services" @click.prevent="scrollTo('services')">Hizmetlerimiz</a>
        <a href="#products" @click.prevent="scrollTo('products')">Ürünlerimiz</a>
        <a href="#about" @click.prevent="scrollTo('about')">Hakkımızda</a>
        <a href="#contact" @click.prevent="scrollTo('contact')">İletişim</a>
        <router-link to="/admin/login" class="admin-link">Admin Paneli</router-link>
      </div>
    </nav>

    <main>
      <router-view v-slot="{ Component }">
        <transition name="fade" mode="out-in">
          <component :is="Component" />
        </transition>
      </router-view>
    </main>
  </div>
</template>

<script setup>
import { computed } from 'vue';
import { useRoute } from 'vue-router';

const route = useRoute();
const isAdminRoute = computed(() => {
  return route.path.startsWith('/admin');
});

const scrollTo = (elementId) => {
  const element = document.getElementById(elementId);
  if (element) {
    element.scrollIntoView({ behavior: 'smooth' });
  }
};
</script>

<style lang="scss">
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: Arial, sans-serif;
  line-height: 1.6;
  background: #f5f5f5;
}

.app {
  min-height: 100vh;
}

.navbar {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  background: #fff;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
  padding: 1rem 2rem;
  z-index: 1000;
  display: flex;
  justify-content: space-between;
  align-items: center;

  .nav-brand {
    .home-link {
      text-decoration: none;
      color: #333;
      font-size: 1.5rem;
      font-weight: bold;
      display: flex;
      align-items: center;
      gap: 0.5rem;
      transition: color 0.3s;

      &:hover {
        color: #007bff;
      }

      i {
        font-size: 1.2rem;
      }
    }
  }

  .nav-links {
    display: flex;
    gap: 2rem;
    align-items: center;

    a {
      text-decoration: none;
      color: #333;
      font-weight: 500;
      transition: color 0.3s;

      &:hover {
        color: #007bff;
      }
    }

    .admin-link {
      color: #007bff;
    }
  }
}

main {
  margin-top: 60px;
  min-height: calc(100vh - 60px);
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

@media (max-width: 768px) {
  .navbar {
    flex-direction: column;
    padding: 1rem;

    .nav-links {
      flex-wrap: wrap;
      justify-content: center;
      gap: 1rem;
      margin-top: 1rem;
    }
  }

  main {
    margin-top: 100px;
  }
}
</style> 