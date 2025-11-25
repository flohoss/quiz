<script setup lang="ts">
import { useGlobalState } from '../store';
import { Icon } from '@iconify/vue';

const state = useGlobalState();
</script>

<template>
  <div class="flex gap-5">
    <button @click="state.previousIndex" :disabled="state.index.value === 0" class="btn btn-secondary btn-square">
      <Icon class="size-6" icon="heroicons-outline:arrow-left" />
    </button>
    <button
      :disabled="state.index.value === state.questions.value.length - 1 || !state.selected.value"
      @click="state.nextIndex"
      class="btn btn-secondary btn-square"
    >
      <Icon class="size-6" icon="heroicons-outline:arrow-right" />
    </button>
  </div>
</template>
