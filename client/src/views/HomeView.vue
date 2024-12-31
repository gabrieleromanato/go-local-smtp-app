<script setup>
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";
import { faPaperPlane } from "@fortawesome/free-regular-svg-icons";
import { isAuthenticated } from "@/api/auth";
import { useRouter } from "vue-router";
import { ref, onMounted } from "vue";
import { getEmails, deleteEmail, searchEmails } from "@/api/emails";
import Emails from "@/components/Emails.vue";
import Pagination from "@/components/Pagination.vue";
import EmailDetails from "@/components/EmailDetails.vue";
import EmailForm from "@/components/EmailForm.vue";
import APIDetails from "@/components/APIDetails.vue";
import SearchForm from "@/components/SearchForm.vue";

const router = useRouter();
const data = ref({
  emails: [],
  pages: 0,
  currentPage: 1,
});

const email = ref({
  id: "",
  subject: "",
  sent_at: "",
  from: "",
  to: [],
  body: "",
  attachments: [],
});
const active = ref(false);
const emailFormVisible = ref(false);
const apiDetailsVisible = ref(false);
const searchQuery = ref("");

const handleNext = (page) => {
  if(!searchQuery.value) {
    getEmails(page).then((response) => {
      data.value.emails = response.emails;
      data.value.pages = response.pages;
      data.value.currentPage = response.page;
    });
  } else {
    searchEmails(searchQuery.value, page).then((response) => {
      data.value.emails = response.emails;
      data.value.pages = response.pages;
      data.value.currentPage = response.page;
    });
  }
};

const handlePrevious = (page) => {
  if(!searchQuery.value) {
    getEmails(page).then((response) => {
      data.value.emails = response.emails;
      data.value.pages = response.pages;
      data.value.currentPage = response.page;
    });
  } else {
    searchEmails(searchQuery.value, page).then((response) => {
      data.value.emails = response.emails;
      data.value.pages = response.pages;
      data.value.currentPage = response.page;
    });
  }
};

const handleSelect = (data) => {
  active.value = true;
  email.value = data;
};

const handleClose = () => {
  active.value = false;
};

const handleSend = () => {
  getEmails().then((response) => {
    data.value.emails = response.emails;
    data.value.pages = response.pages;
    data.value.currentPage = response.page;
    emailFormVisible.value = false;
  });
};

const handleDelete = (id) => {
  deleteEmail(id).then(() => {
    getEmails().then((response) => {
      data.value.emails = response.emails;
      data.value.pages = response.pages;
      data.value.currentPage = response.page;
    });
  });
};

const handleSearch = (query) => {
  searchQuery.value = query;
  searchEmails(query).then((response) => {
    data.value.emails = response.emails;
    data.value.pages = response.pages;
    data.value.currentPage = response.page;
  });
};

const handleReset = () => {
  searchQuery.value = "";
  getEmails().then((response) => {
    data.value.emails = response.emails;
    data.value.pages = response.pages;
    data.value.currentPage = response.page;
  });
};

document.title = "Axio SMTP Server | Home";

if (!isAuthenticated()) {
  router.push({ name: "login" });
}

onMounted(() => {
  getEmails().then((response) => {
    data.value.emails = response.emails;
    data.value.pages = response.pages;
    data.value.currentPage = response.page;
  });
});
</script>

<template>
  <div class="home-view">
    <section class="xl:container mx-auto bg-white shadow-md p-4 rounded-lg">
      <div class="actions flex justify-between items-center">
        <button class="btn" @click="emailFormVisible = !emailFormVisible">
          <FontAwesomeIcon :icon="faPaperPlane" />
        </button>
        <SearchForm @search="handleSearch" @reset="handleReset" />
        <EmailForm
          @send="handleSend"
          @close="emailFormVisible = !emailFormVisible"
          :visible="emailFormVisible"
        />
      </div>
      <Emails
        @delete="handleDelete"
        @select="handleSelect"
        :emails="data.emails"
        v-if="data.emails.length > 0"
      />
      <Pagination
        v-if="data.pages > 1"
        :pages="data.pages"
        :currentPage="data.currentPage"
        @next="handleNext"
        @previous="handlePrevious"
      />
      <div v-if="data.emails.length === 0" class="no-emails">
        <p>No emails found.</p>
      </div>
      <APIDetails
        :visible="apiDetailsVisible"
        @toggle="apiDetailsVisible = !apiDetailsVisible"
      />
    </section>
  </div>
  <EmailDetails @close="handleClose" :email="email" :active="active" />
</template>
