<script setup>
const props = defineProps(['pages', 'currentPage']);

const emits = defineEmits(['next', 'previous']);

const nextPage = () => {
    if (props.currentPage === props.pages) return;
    emits('next', props.currentPage + 1);
};

const previousPage = () => {
    if (props.currentPage === 1) return;
    emits('previous', props.currentPage - 1);
};
</script>

<template>
    <div class="pagination">
        <button @click="previousPage" class="button" :disabled="currentPage === 1">Previous</button>
        <button @click="nextPage" class="button" :disabled="currentPage === pages">Next</button>
    </div>
</template>