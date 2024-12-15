<script setup>
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';
import { faArrowLeft, faArrowRight } from '@fortawesome/free-solid-svg-icons';
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
        <button @click="previousPage" class="button" :disabled="currentPage === 1">
            <FontAwesomeIcon :icon="faArrowLeft" />
        </button>
        <button @click="nextPage" class="button" :disabled="currentPage === pages">
            <FontAwesomeIcon :icon="faArrowRight" />
        </button>
    </div>
</template>