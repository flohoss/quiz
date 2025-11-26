<script setup lang="ts">
import { useGlobalState } from '../store';
import { Icon } from '@iconify/vue';
import ThemeChanger from './ThemeChanger.vue';

const { start, end, question, submitted, previousIndex, nextIndex, submit, reset } = useGlobalState();
</script>

<template>
  <div class="flex justify-between">
    <button @click="previousIndex" :disabled="start || submitted" class="nav-btn btn-secondary">
      <Icon class="size-6" icon="fa6-solid:hand-point-left" />
    </button>
    <div class="flex items-center gap-5">
      <ThemeChanger />
    </div>
    <button v-if="end && !submitted" :disabled="!question?.answer" @click="submit" class="nav-btn btn-primary">
      <Icon class="size-6" icon="fa6-solid:gifts" />
    </button>
    <button v-else-if="!submitted" :disabled="!question?.answer" @click="nextIndex" class="nav-btn btn-primary">
      <Icon class="size-6" icon="fa6-solid:sleigh" />
    </button>
    <button v-else @click="reset" class="nav-btn btn-warning">
      <Icon class="size-6" icon="heroicons-outline:refresh" />
    </button>
  </div>
</template>
