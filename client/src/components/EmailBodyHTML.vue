<script setup>
import { ref, onMounted } from 'vue';
const props = defineProps(['body'])

const html = ref('');
const dataUrl = ref('');

onMounted(() => {
  html.value = props.body;
  dataUrl.value = `data:text/html;base64,${btoa(html.value)}`;
});
</script>

<template>
  <iframe class="body-html" :src="dataUrl" width="100%" height="100%"></iframe>
</template>