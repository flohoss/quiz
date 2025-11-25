<script setup lang="ts">
import { useGlobalState } from '../store';
import { Icon } from '@iconify/vue';

const { start, end, question, previousIndex, nextIndex, submit } = useGlobalState();
</script>

<template>
  <div class="flex gap-5">
    {{ question }}
    <button @click="previousIndex" :disabled="start" class="btn btn-secondary btn-circle btn-lg">
      <Icon class="size-6" icon="heroicons-outline:arrow-left" />
    </button>
    <button v-if="end" :disabled="!question?.answer" @click="submit" class="btn btn-primary btn-circle btn-lg">
      <Icon class="size-6" icon="heroicons-outline:check" />
    </button>
    <button v-else :disabled="!question?.answer" @click="nextIndex" class="btn btn-primary btn-circle btn-lg">
      <Icon class="size-6" icon="heroicons-outline:arrow-right" />
    </button>
  </div>
</template>
