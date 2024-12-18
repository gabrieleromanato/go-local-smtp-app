<script setup>
import { ref } from 'vue';
import { logout, isAuthenticated } from '@/api/auth';

const isLogged = ref(isAuthenticated());

const handleLogout = () => {
    logout();
    window.location.href = '/login';
};
</script>

<template>
    <nav class="navigation">
        <div class="wrapper">
            <ul>
                <li v-if="isLogged">
                    <button @click="handleLogout">Logout</button>
                </li>
            </ul>
        </div>
    </nav>
</template>