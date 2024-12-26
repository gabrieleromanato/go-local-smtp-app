<script setup>
import { ref } from 'vue';

const query = ref('');

const emits = defineEmits(['search', 'reset']);

const handleSearch = () => {
    emits('search', query.value);
};

const handleReset = () => {
    query.value = '';
    emits('reset');
};
</script>

<template>
    <form class="search-form d-flex align-items-center">
        <input type="text" v-model="query" class="form-control me-2" placeholder="Search" />
        <button @click="handleSearch" type="button" class="btn me-2">Search</button>
        <button @click="handleReset" type="button" class="btn">Clear</button>
    </form>
</template>