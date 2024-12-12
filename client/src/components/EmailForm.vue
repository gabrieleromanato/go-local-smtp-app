<script setup>
import { ref } from 'vue'
import { sendEmail } from '@/api/emails';

defineProps({
    visible: Boolean
});

const message = ref({
    email: '',
    message: '',
    subject: '',
    recipient: '',
    attachments: []   
});

const handleFileUpload = (e) => {
    const files = e.target.files;
    message.attachments = files;
};

const handleSubmit = () => {
    sendEmail(message.value).then(() => {
        message.value = {
            email: '',
            message: '',
            subject: '',
            recipient: '',
            attachments: []
        };
        window.location.reload();
    });
};
</script>

<template>
    <form @submit.prevent="handleSubmit" class="email-form" :class="{'visible': visible}">
        <div class="form-group">
            <label class="form-label">Email</label>
            <input type="email" class="form-control" v-model="message.email" />
        </div>
        <div class="form-group">
            <label class="form-label">Subject</label>
            <input type="text" class="form-control" v-model="message.subject" />
        </div>
        <div class="form-group">
            <label class="form-label">Message</label>
            <textarea class="form-control" v-model="message.message"></textarea>
        </div>
        <div class="form-group">
            <label class="form-label">Recipient</label>
            <input type="email" class="form-control" v-model="message.recipient" />
        </div>
        <div class="form-group">
            <label class="form-label">Attachments</label>
            <input type="file" class="form-control-file" multiple @change="handleFileUpload" />
        </div>
        <button type="submit" class="button">Send</button>
    </form>
</template>