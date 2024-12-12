<script setup>
import { isAuthenticated } from '@/api/auth';
import { useRouter } from 'vue-router';
import { ref, onMounted } from 'vue';
import { getEmails, deleteEmail  } from '@/api/emails';
import Emails from '@/components/Emails.vue';
import Pagination from '@/components/Pagination.vue';
import EmailDetails from '@/components/EmailDetails.vue';

const router = useRouter();
const data = ref({
    emails: [],
    pages: 0,
    currentPage: 1
});

const email = ref({
    id: '',
    subject: '',
    sent_at: '',
    from: '',
    to: [],
    body: '',
    attachments: []
});
const active = ref(false);

const handleNext = (page) => {
    getEmails(page).then(response => {
        data.value.emails = response.emails;
        data.value.pages = response.pages;
        data.value.currentPage = response.page;
    });
};

const handlePrevious = (page) => {
    getEmails(page).then(response => {
        data.value.emails = response.emails;
        data.value.pages = response.pages;
        data.value.currentPage = response.page;
    });
};

const handleSelect = (data) => {
    active.value = true;
    email.value = data;
};

const handleClose = () => {
    active.value = false;
};

const handleDelete = (id) => {
    deleteEmail(id).then(() => {
        getEmails().then(response => {
            data.value.emails = response.emails;
            data.value.pages = response.pages;
            data.value.currentPage = response.page;
        });
    });
};

document.title = 'Emails | Home';

if (!isAuthenticated()) {
    router.push({ name: 'login' });
}

onMounted(() => {
    getEmails().then(response => {
        data.value.emails = response.emails;
        data.value.pages = response.pages;
        data.value.currentPage = response.page;
        console.log(data.value.emails);
    });
});
</script>

<template>
  <div class="home-view">
    <section class="container">
        <Emails @delete="handleDelete" @select="handleSelect" :emails="data.emails" v-if="data.emails.length > 0" />
        <Pagination :pages="data.pages" :currentPage="data.currentPage" @next="handleNext" @previous="handlePrevious" />
    </section>
  </div>
  <EmailDetails @close="handleClose" :email="email" :active="active" />
</template>
