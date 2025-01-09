<script setup>
import { ref } from 'vue';
import { register } from '@/api/auth';
import { useRouter } from 'vue-router';

const userData = ref({
    email: '',
    password: ''
});

const message = ref('');

const router = useRouter();

const handleSubmit = async () => {
    message.value = '';
    try {
        const response = await register(userData.value);
        if(response.message) {
            router.push('/login');
        } else {
            message.value = 'Error during registration';
        }
    } catch (error) {
        console.error(error);
        message.value = 'Error during registration';
    }
};

document.title = 'Axio SMTP Server | Register';
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
            <input v-model="userData.password" type="text" class="form-control" />
        </div>
        <div v-if="message" class="form-error">{{ message }}</div>
        <button type="submit" class="btn">Register</button>
    </form>
    </div>
</template>