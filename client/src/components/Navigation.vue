<script setup>
import { ref } from 'vue';
import { logout, isAuthenticated } from '@/api/auth';
import { RouterLink } from 'vue-router';

const isLogged = ref(isAuthenticated());

const handleLogout = () => {
    logout();
    window.location.href = '/login';
};
</script>

<template>
    <nav class="navigation">
        <div class="wrapper">
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