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
  <table class="emails">
    <thead>
      <tr>
        <th>Subject</th>
        <th>Date</th>
        <th>From</th>
        <th colspan="2">To</th>
      </tr>
    </thead>
      <tbody>
        <tr v-for="email in emails" :key="email.id">
          <td>{{ email.subject }}</td>
          <td>{{ email.sent_at }}</td>
          <td>{{ email.from }}</td>
          <td>{{ email.to[0] }}</td>
          <td>
            <button @click="handleSelect(email)" class="button select">View</button>
            <button @click="handleDelete(email)" class="button delete">Delete</button>
          </td>
        </tr>
        </tbody>
    </table>
</template>
        