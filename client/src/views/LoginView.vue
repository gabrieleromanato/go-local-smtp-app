<script setup>
import { ref } from 'vue';
import { login, setToken, setUserId  } from '@/api/auth';
import { useRouter, RouterLink } from 'vue-router';

const userData = ref({
    email: '',
    password: ''
});

const message = ref('');

const router = useRouter();

const handleSubmit = async () => {
    message.value = '';
    try {
        const response = await login(userData.value);
        if(response.token) {
            setToken(response.token);
            setUserId(response.userId);
            router.push('/');
        } else {
            message.value = 'Error during login';
        }
    } catch (error) {
        console.error(error);
        message.value = 'Error during login';
    }
};

document.title = 'Axio SMTP Server | Login';
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
        <div v-if="message" class="form-error">{{ message }}</div>
        <button type="submit" class="btn">Login</button>
        <div class="register-link">
            <RouterLink to="/register">Register</RouterLink>
        </div>
    </form>
    </div>
</template>