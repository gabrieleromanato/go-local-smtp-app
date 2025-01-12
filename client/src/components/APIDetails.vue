<script setup>
import { ref } from 'vue';

defineProps({
    visible: Boolean
});

defineEmits(['toggle']);

const data = ref({
    api: {
        name: 'Emails API',
        description: 'This is a RESTful API that allows you to send and receive emails.',
        endpoints: [
            {
                method: 'POST',
                path: '/login',
                description: 'Login to the application',
                body: {
                    email: 'string',
                    password: 'string'
                },
                response: {
                    token: 'string'
                }
            },
            {
                method: 'GET',
                path: '/api/emails',
                description: 'Get all emails',
                headers: {
                    Authorization: 'Bearer token'
                },
                response: {
                    emails: 'array',
                    pages: 'number',
                    page: 'number'
                },
                body: {
                    page: 'number'
                }
            },
            {
                method: 'POST',
                path: '/api/emails',
                description: 'Send an email',
                body: {
                    email: 'string',
                    message: 'string',
                    message_html: 'string',
                    subject: 'string',
                    recipient: 'string',
                    attachments: 'array',
                    user_id: 'number'
                },
                headers: {
                    Authorization: 'Bearer token'
                },
                response: {
                    message: 'string'
                }
            },
            {
                method: 'DELETE',
                path: '/api/emails/:id',
                description: 'Delete an email by ID',
                headers: {
                    Authorization: 'Bearer token'
                },
                response: {
                    message: 'string'
                }
            }
        ]
    }});
</script>

<template>
    <div class="api-details-toggle">
        <button @click="$emit('toggle')" class="api-details-toggle__button">
            <span v-if="visible">Hide API Details</span>
            <span v-else>Show API Details</span>
        </button>
    </div>
  <div class="api-details" :class="{'visible': visible}">
    <h1>{{ data.api.name }}</h1>
    <p>{{ data.api.description }}</p>
    <h2>Endpoints</h2>
    <ul>
      <li v-for="endpoint in data.api.endpoints" :key="endpoint.path">
        <h3><span>{{ endpoint.method }}</span> <span>{{ endpoint.path }}</span></h3>
        <p>{{ endpoint.description }}</p>
        <h4>Request</h4>
        <pre>{{ JSON.stringify(endpoint.body, null, 2) }}</pre>
        <h4>Response</h4>
        <pre>{{ JSON.stringify(endpoint.response, null, 2) }}</pre>
      </li>
    </ul>
  </div>
</template>