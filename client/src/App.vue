<script setup>
import { RouterView, useRouter } from 'vue-router'
import { checkToken, isAuthenticated  } from './api/auth';
import Navigation from './components/Navigation.vue';

const checkTokenValidity = () => {
  const token = localStorage.getItem('app-token');
  const router = useRouter();
  if (token) {
    checkToken()
      .then((resp) => {
        if(resp.error) {
          localStorage.removeItem('app-token');
          router.push('/login');
        }
      })
      .catch((err) => {
        console.log(err);
      });
  }
};

checkTokenValidity();
</script>

<template>
 <main class="app">
  <Navigation v-if="isAuthenticated()"/> 
  <RouterView />
  </main>
</template>
