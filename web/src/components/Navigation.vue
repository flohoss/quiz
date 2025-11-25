<script setup lang="ts">
import { computed } from 'vue';
import { useGlobalState } from '../store';
import { Icon } from '@iconify/vue';

const state = useGlobalState();
const end = computed(() => state.index.value === state.questions.value.length - 1);
</script>

<template>
  <div class="flex gap-5">
    <button @click="state.previousIndex" :disabled="state.index.value === 0" class="btn btn-secondary btn-circle btn-lg">
      <Icon class="size-6" icon="heroicons-outline:arrow-left" />
    </button>
    <button v-if="end && state.selected.value" @click="state.submit" class="btn btn-primary btn-circle btn-lg">
      <Icon class="size-6" icon="heroicons-outline:check" />
    </button>
    <button v-else :disabled="end || !state.selected.value" @click="state.nextIndex" class="btn btn-primary btn-circle btn-lg">
      <Icon class="size-6" icon="heroicons-outline:arrow-right" />
    </button>
  </div>
</template>
