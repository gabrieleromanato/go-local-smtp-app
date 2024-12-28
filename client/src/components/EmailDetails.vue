<script setup>
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';
import { faTimes, faFile } from '@fortawesome/free-solid-svg-icons';
import { ref, watch } from 'vue';

const props = defineProps({
    email: Object,
    active: Boolean
});

defineEmits(['close']);

const attachments = ref([]);
const BASE_URL = import.meta.env.VITE_ASSETS_URL;
const imageExtensions = ['jpg', 'jpeg', 'png', 'gif'];
const videoExtensions = ['mp4', 'webm', 'ogg'];
const audioExtensions = ['mp3', 'wav', 'ogg'];
const documentExtensions = ['pdf', 'doc', 'docx', 'xls', 'xlsx', 'ppt', 'pptx'];

const getFileType = (filename) => {
    const extension = filename.split('.').pop();
    if (imageExtensions.includes(extension)) {
        return 'image';
    } else if (videoExtensions.includes(extension)) {
        return 'video';
    } else if (audioExtensions.includes(extension)) {
        return 'audio';
    } else if (documentExtensions.includes(extension)) {
        return 'document';
    } else {
        return '';
    }
};

watch(() => props.active, (newVal) => {
    if (newVal) {
        attachments.value = props.email.attachments.map(attachment => {
            const url = `${BASE_URL}/attachments/${attachment.filename}`;
            return {
                name: attachment.filename,
                url: url,
                type: getFileType(attachment.filename)
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
                    <a v-if="attachment.type && attachment.type !== 'image'" :href="attachment.url" target="_blank" class="document-attachment">
                        <FontAwesomeIcon :icon="faFile" />
                        <div>{{ attachment.name }}</div>
                    </a>
                    <figure v-if="attachment.type === 'image'" class="image-attachment">
                      <a :href="attachment.url" target="_blank">  
                        <img :src="attachment.url" :alt="attachment.name" class="img-fluid" />
                        <figcaption>{{ attachment.name }}</figcaption>
                      </a>
                    </figure>
                </li>
            </ul>
        </div>
    </aside>
    </div>
</template>