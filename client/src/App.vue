<script setup>
import { RouterView } from 'vue-router'
import { checkToken  } from './api/auth';

const checkTokenValidity = () => {
  const token = localStorage.getItem('app-token');
  if (token) {
    checkToken()
      .then((resp) => {
        if(resp.error) {
          localStorage.removeItem('app-token');
          window.location.href = '/login';
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
  <RouterView />
  </main>
</template>
