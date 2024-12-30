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
  <table class="emails table-fixed border-collapse border border-slate-500">
    <thead>
      <tr>
        <th class="border border-slate-200">Subject</th>
        <th class="border border-slate-200">Date</th>
        <th class="border border-slate-200">From</th>
        <th colspan="2" class="border border-slate-200">To</th>
      </tr>
    </thead>
    <tbody>
      <tr v-for="email in emails" :key="email.id">
        <td class="border border-slate-200">{{ email.subject }}</td>
        <td class="border border-slate-200">{{ email.sent_at }}</td>
        <td class="border border-slate-200">{{ email.from }}</td>
        <td class="border border-slate-200">{{ email.to[0] }}</td>
        <td class="border border-slate-200">
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
