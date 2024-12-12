<script setup>
defineProps(['emails']);

const emits = defineEmits(['select', 'delete']);

const handleSelect = (email) => {
    emits('select', email);
};

const handleDelete = (email) => {
    emits('delete', email.id);
};
</script>

<template>
  <div class="emails">
    <header>
        <div>Subject</div>
        <div>Date</div>
        <div>From</div>
        <div>To</div>
    </header>
      <ul class="emails-list">
        <li v-for="email in emails" :key="email.id" class="email">
          <div class="email-header" @click="handleSelect(email)">
            <h4 class="email-title">{{ email.subject }}</h4>
            <p class="email-date">{{ email.sent_at }}</p>
            <p class="email-from">{{ email.from }}</p>
            <p class="email-to">{{ email.to[0] }}</p>
            
          </div>
          <button @click="handleDelete(email)" class="button delete">Delete</button>
        </li>
        </ul>
    </div>
</template>
        