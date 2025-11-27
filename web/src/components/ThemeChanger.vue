<script setup lang="ts">
import { Icon } from '@iconify/vue';
import { useToggle } from '@vueuse/core';
import { IsDark } from '../main';

const toggleDark = useToggle(IsDark);
</script>

<template>
  <label class="swap swap-rotate">
    <input type="checkbox" :checked="IsDark" @change="toggleDark()" class="hidden" />

    <Icon class="swap-on size-8" icon="tabler:moon" />
    <Icon class="swap-off size-8" icon="tabler:sun" />
  </label>
</template>
