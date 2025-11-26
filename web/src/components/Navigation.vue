<script setup lang="ts">
import { useGlobalState } from '../store';
import { Icon } from '@iconify/vue';
import ThemeChanger from './ThemeChanger.vue';

const { start, end, question, submitted, previousIndex, nextIndex, submit, reset } = useGlobalState();
</script>

<template>
  <template v-if="submitted">
    <div class="flex-1"></div>
    <button @click="reset" class="btn btn-warning btn-circle btn-lg">
      <Icon class="size-6" icon="heroicons-outline:refresh" />
    </button>
  </template>
  <template v-else>
    <button @click="previousIndex" :disabled="start" class="btn btn-primary btn-circle btn-lg">
      <Icon class="size-6" icon="heroicons-outline:arrow-left" />
    </button>
    <div class="flex items-center gap-5">
      <ThemeChanger />
    </div>
    <button v-if="end" :disabled="!question?.answer" @click="submit" class="btn btn-success btn-circle btn-lg">
      <Icon class="size-6" icon="heroicons-outline:check" />
    </button>
    <button v-else :disabled="!question?.answer" @click="nextIndex" class="btn btn-primary btn-circle btn-lg">
      <Icon class="size-6" icon="heroicons-outline:arrow-right" />
    </button>
  </template>
</template>
