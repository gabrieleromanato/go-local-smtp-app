<script setup>
import { ref } from 'vue';
import { login, setToken  } from '@/api/auth';

const userData = ref({
    email: '',
    password: ''
});

const handleSubmit = async () => {
    try {
        const response = await login(userData.value);
        setToken(response.token);
        window.location.href = '/';
    } catch (error) {
        console.error(error);
    }
};

document.title = 'Emails | Login';
</script>

<template>
  <div class="login-view">  
    <form class="login-form" @submit.prevent="handleSubmit">
        <div class="login-logo">
            <div class="logo">
                <span class="visually-hidden">Emails</span>
            </div>
        </div>
        <div class="form-group">
            <label class="form-label">Email</label>
            <input v-model="userData.email" type="email" class="form-control" />  
        </div>
        <div class="form-group">
            <label class="form-label">Password</label>
            <input v-model="userData.password" type="password" class="form-control" />
        </div>
        <button type="submit" class="button">Login</button>
    </form>
    </div>
</template>