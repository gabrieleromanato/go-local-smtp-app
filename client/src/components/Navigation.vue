<script setup>
import { ref } from "vue";
import { logout, isAuthenticated } from "@/api/auth";
import { RouterLink, useRouter } from "vue-router";

const isLogged = ref(isAuthenticated());
const router = useRouter();

const handleLogout = () => {
  logout();
  isLogged.value = false;
  router.push("/login");
};
</script>

<template>
  <nav class="navigation mt-8 mb-8">
    <div class="bg-white shadow-md p-4 rounded-lg xl:container mx-auto flex justify-between items-center">
      <div class="logo">
        <RouterLink to="/">
          <span class="visually-hidden">Emails</span>
        </RouterLink>
      </div>
      <ul>
        <li v-if="isLogged">
          <button @click="handleLogout">Logout</button>
        </li>
      </ul>
    </div>
  </nav>
</template>
