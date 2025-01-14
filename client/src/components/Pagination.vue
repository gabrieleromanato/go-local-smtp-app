<script setup>
import { ref, watch } from "vue";
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";
import { faArrowLeft, faArrowRight } from "@fortawesome/free-solid-svg-icons";

const props = defineProps(["pages", "currentPage"]);
const emits = defineEmits(["next", "previous"]);

const links = ref([]);

const updatePaginationLinks = () => {
  const pages = props.pages;
  const currentPage = props.currentPage;
  const maxVisibleLinks = 4;

  let start = Math.max(1, currentPage - Math.floor(maxVisibleLinks / 2));
  let end = Math.min(pages, start + maxVisibleLinks - 1);

  if (end - start + 1 < maxVisibleLinks) {
    start = Math.max(1, end - maxVisibleLinks + 1);
  }

  links.value = [];
  for (let i = start; i <= end; i++) {
    links.value.push(i);
  }
};

watch(
  () => props.currentPage,
  () => {
    updatePaginationLinks();
  },
  { immediate: true }
);

const nextPage = () => {
  if (props.currentPage === props.pages) return;
  emits("next", props.currentPage + 1);
};

const previousPage = () => {
  if (props.currentPage === 1) return;
  emits("previous", props.currentPage - 1);
};
</script>

<template>
  <div class="pagination">
    <button @click="previousPage" class="btn" :disabled="currentPage === 1">
      <FontAwesomeIcon :icon="faArrowLeft" />
    </button>
    <button
      v-for="link in links"
      :key="link"
      @click="emits('next', link)"
      class="btn"
      :class="{ active: currentPage === link }">
      {{ link }}
    </button>
    <button @click="nextPage" class="btn" :disabled="currentPage === pages">
      <FontAwesomeIcon :icon="faArrowRight" />
    </button>
  </div>
</template>
