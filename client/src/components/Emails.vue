<script setup>
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";
import { faTrashAlt } from "@fortawesome/free-regular-svg-icons";
import { faEye } from "@fortawesome/free-regular-svg-icons";

defineProps(["emails"]);

const emits = defineEmits(["select", "delete"]);

const handleSelect = (email) => {
  emits("select", email);
};

const handleDelete = (email) => {
  emits("delete", email.id);
};
</script>

<template>
  <table class="emails table table-striped">
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
          <button @click="handleSelect(email)" class="select">
            <FontAwesomeIcon :icon="faEye" />
          </button>
          <button @click="handleDelete(email)" class="delete">
            <FontAwesomeIcon :icon="faTrashAlt" />
          </button>
        </td>
      </tr>
    </tbody>
  </table>
</template>
