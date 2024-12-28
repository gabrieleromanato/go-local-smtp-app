<script setup>
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';
import { faTimes } from '@fortawesome/free-solid-svg-icons';
import { ref, watch } from 'vue';

const props = defineProps({
    email: Object,
    active: Boolean
});

defineEmits(['close']);

const attachments = ref([]);
const BASE_URL = import.meta.env.VITE_ASSETS_URL;

watch(() => props.active, (newVal) => {
    if (newVal) {
        attachments.value = props.email.attachments.map(attachment => {
            const url = `${BASE_URL}/attachments/${attachment.filename}`;
            return {
                name: attachment.filename,
                url: url
            };
        });
    }
});
</script>

<template>
  <div class="email-details-wrapper" :class="{'active': active}">
    <button @click="$emit('close')" class="email-details__close">
        <FontAwesomeIcon :icon="faTimes" />
    </button>  
    <aside class="email-details">
        
        <div class="email-details__header">
            <h2 class="email-details__title">{{ email.subject }}</h2>
            <p  class="email-details__date">{{ email.sent_at }}</p>
        </div>
        <div class="email-details__body">
            <p>{{ email.body }}</p>
        </div>
        <div class="email-details__attachments" v-if="attachments.length > 0">
            <ul>
                <li v-for="attachment in attachments" :key="attachment.name">
                    <a :href="attachment.url" target="_blank">{{ attachment.name }}</a>
                </li>
            </ul>
        </div>
    </aside>
    </div>
</template>