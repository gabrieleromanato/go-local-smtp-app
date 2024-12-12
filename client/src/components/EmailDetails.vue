<script setup>
import { ref, watch } from 'vue';

const props = defineProps({
    email: Object,
    active: Boolean
});

defineEmits(['close']);

const getTypeFromFileExtension = (extension) => {
    const images = ['png', 'jpg', 'jpeg', 'gif', 'svg'];
    const documents = ['pdf', 'doc', 'docx', 'xls', 'xlsx', 'ppt', 'pptx'];
    const audio = ['mp3', 'wav', 'ogg', 'flac'];
    const video = ['mp4', 'avi', 'mkv', 'mov', 'wmv'];

    let mime = 'application/octet-stream';

    if (images.includes(extension)) {
        if (extension === 'svg') {
            mime = 'image/svg+xml';
        } else if (extension === 'jpg' || extension === 'jpeg') {
            mime = 'image/jpeg';
        } else {
            mime = `image/${extension}`;
        }
    } else if (documents.includes(extension)) {
        if (extension === 'pdf') {
            mime = 'application/pdf';
        } else if (extension === 'doc') {
            mime = 'application/msword';
        } else if (extension === 'docx') {
            mime = 'application/vnd.openxmlformats-officedocument.wordprocessingml.document';
        } else if (extension === 'xls') {
            mime = 'application/vnd.ms-excel';
        } else if (extension === 'xlsx') {
            mime = 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet';
        } else if (extension === 'ppt') {
            mime = 'application/vnd.ms-powerpoint';
        } else if (extension === 'pptx') {
            mime = 'application/vnd.openxmlformats-officedocument.presentationml.presentation';
        }
    } else if (audio.includes(extension)) {
        mime = `audio/${extension}`;
    } else if (video.includes(extension)) {
        mime = `video/${extension}`;
    }
    return mime;
};

const attachments = ref([]);

watch(() => props.active, (newVal) => {
    if (newVal) {
        attachments.value = props.email.attachments.map(attachment => {
            const extension = attachment.filename.split('.').pop();
            const type = getTypeFromFileExtension(extension);
            const url = attachment.content.includes('attachments') ? `https://localhost:8080/attachments/${attachment.filename}` : `data:${type};base64,${attachment.content}`;
            return {
                name: attachment.filename,
                type: type,
                url: `data:${type};base64,${attachment.content}`
            };
        });
    }
});
</script>

<template>
    <aside class="email-details" :class="{'active': active}">
        <button @click="$emit('close')" class="email-details__close">Close</button>
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
                    <a :href="attachment.url" :download="attachment.name">{{ attachment.name }}</a>
                </li>
            </ul>
        </div>
    </aside>
</template>