<script setup>
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";
import { faTimes, faPaperPlane } from "@fortawesome/free-solid-svg-icons";
import { ref } from "vue";
import { sendEmail } from "@/api/emails";

defineProps({
  visible: Boolean,
});

const emits = defineEmits(["close", "send"]);

const message = ref({
  email: "",
  message: "",
  message_html: "",
  subject: "",
  recipient: "",
  attachments: [],
});

const handleFileUpload = (e) => {
  const files = e.target.files;
  message.value.attachments = files;
};

const handleSubmit = () => {
  sendEmail(message.value).then(() => {
    message.value = {
      email: "",
      message: "",
      message_html: "",
      subject: "",
      recipient: "",
      attachments: [],
    };
    emits("send");
  });
};
</script>

<template>
  <div class="email-form-container" :class="{ visible: visible }">
    <button @click="$emit('close')" class="close-button">
      <FontAwesomeIcon :icon="faTimes" />
    </button>
    <form @submit.prevent="handleSubmit" class="email-form">
      <div class="form-group">
        <label class="form-label">Email</label>
        <input type="email" class="form-control" v-model="message.email" />
      </div>
      <div class="form-group">
        <label class="form-label">Subject</label>
        <input type="text" class="form-control" v-model="message.subject" />
      </div>
      <div class="form-group">
        <label class="form-label">Message</label>
        <textarea class="form-control" v-model="message.message"></textarea>
      </div>
      <div class="form-group">
        <label class="form-label">HTML Message</label>
        <textarea class="form-control" v-model="message.message_html"></textarea>
      </div>
      <div class="form-group">
        <label class="form-label">Recipient</label>
        <input type="email" class="form-control" v-model="message.recipient" />
      </div>
      <div class="form-group">
        <label class="form-label">Attachments</label>
        <input
          type="file"
          class="form-control-file"
          multiple
          @change="handleFileUpload"
        />
      </div>
      <button type="submit" class="btn btn-block w-100">
        <FontAwesomeIcon :icon="faPaperPlane" />
      </button>
    </form>
  </div>
</template>
