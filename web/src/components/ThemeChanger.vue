<script setup lang="ts">
import { watch } from 'vue';
import { useLocalStore } from '../localStore';
import { Icon } from '@iconify/vue';

const store = useLocalStore();

watch(
  () => store.value.darkMode,
  (newDarkMode) => {
    const theme = newDarkMode ? 'business' : 'corporate';
    document.documentElement.setAttribute('data-theme', theme);
  },
  { immediate: true }
);

const toggleTheme = () => {
  store.value.darkMode = !store.value.darkMode;
};
</script>

<template>
  <label class="swap swap-rotate">
    <input type="checkbox" :checked="store.darkMode" @change="toggleTheme()" class="hidden" />

    <Icon class="swap-on size-8" icon="tabler:sun" />
    <Icon class="swap-off size-8" icon="tabler:moon" />
  </label>
</template>
